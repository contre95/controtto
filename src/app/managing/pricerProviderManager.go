package managing

import (
	"controtto/src/domain/pnl"
	"fmt"
	"log/slog"
	"sync"
)

// PriceProviderManager handles all price provider operations
type PriceProviderManager struct {
	providers map[string]*pnl.PriceProvider
	mu        sync.RWMutex // protects the providers map
}

type QueryPriceReq struct {
	AssetSymbolA string
	AssetSymbolB string
}

type QueryPriceResp struct {
	Price         float64
	ProviderName  string
	ProviderColor string
}

func NewPriceProviderManager(in map[string]pnl.PriceProvider) *PriceProviderManager {
	providers := make(map[string]*pnl.PriceProvider, len(in))
	for key, val := range in {
		v := val // crea una copia para tomar su dirección
		providers[key] = &v
	}
	ppm := &PriceProviderManager{
		providers: providers,
	}
	for key, provider := range providers {
		provider.PriceAPI = provider.Init(provider.Token)
		fmt.Printf("initializing %s: IsSet=%t, NeedsToken=%t\n", key, provider.IsSet, provider.NeedsToken)
		err := ppm.UpdateProvider(key, provider.Token, provider.IsSet)
		if err != nil {
			panic(fmt.Sprintf("could not update provider %s: %v", key, err))
		}
		if provider.IsSet && provider.PriceAPI == nil {
			panic(fmt.Sprintf("provider %s is enabled but PriceAPI is not implemented", key))
		}
	}
	return ppm
}

func (ppm *PriceProviderManager) ListProviders(all bool) map[string]*pnl.PriceProvider {
	ppm.mu.RLock()
	defer ppm.mu.RUnlock()
	filtered := make(map[string]*pnl.PriceProvider)
	for k, v := range ppm.providers {
		if all || v.IsSet {
			filtered[k] = v
		}
	}
	return filtered
}

func (ppm *PriceProviderManager) UpdateProvider(key string, token string, enable bool) error {
	ppm.mu.Lock()
	defer ppm.mu.Unlock()
	provider, ok := ppm.providers[key]
	if !ok {
		return ErrProviderNotFound
	}
	provider.Token = token
	provider.IsSet = enable
	if provider.NeedsToken && enable {
		if token == "" {
			return ErrEmptyToken
		} else {
			provider.PriceAPI = provider.Init(token)
		}
	}
	return nil
}

// QueryPrice gets the current price for an asset pair
func (ppm *PriceProviderManager) QueryPrice(req QueryPriceReq) (*QueryPriceResp, error) {
	slog.Info("Querying prices")
	if req.AssetSymbolA == "" || req.AssetSymbolB == "" {
		return nil, ErrInvalidAssetPair
	}
	resp := &QueryPriceResp{}
	for _, provider := range ppm.providers {
		slog.Info("Checking", "provider", provider.ProviderName, "set", provider.IsSet)
		if !provider.IsSet {
			continue
		}
		price, err := provider.PriceAPI.GetCurrentPrice(req.AssetSymbolA, req.AssetSymbolB)
		if err != nil {
			slog.Error("Error getting price", "provider", provider.ProviderName, "error", err)
			continue
		}
		resp.Price = price
		resp.ProviderColor = provider.Color
		resp.ProviderName = provider.ProviderName
		return resp, nil
	}
	slog.Error("Error getting price", "error", ErrProviderNotFound)
	return nil, ErrProviderNotFound
}

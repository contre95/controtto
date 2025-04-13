package querying

import (
	"controtto/src/domain/pnl"
	"errors"
	"fmt"
	"log/slog"
)

type QueryPriceReq struct {
	AssetSymbolA string
	AssetSymbolB string
}
type QueryPriceResp struct {
	Price         float64
	ProviderName  string
	ProviderColor string
}

type PriceQuerier struct {
	providers pnl.PriceProviders
}

func NewPriceQuerier(p pnl.PriceProviders) *PriceQuerier {
	return &PriceQuerier{p}
}

// GetPrice returns the current base asset value expressed in terms of the quote value
// If is fails to retrieve the value it will set it to 0 (zero)
func (aq *PriceQuerier) GetPrice(req QueryPriceReq) (*QueryPriceResp, error) {
	resp := QueryPriceResp{Price: 0}
	for _, m := range aq.providers {
		// price, err := m.GetCurrentPrice(req.AssetSymbolA, req.AssetSymbolB)
		price, err := m.GetCurrentPrice(req.AssetSymbolA, req.AssetSymbolB)
		if err != nil {
			slog.Error("Could not get current price", "base", req.AssetSymbolA, "quote", req.AssetSymbolB, "provider", m.ProviderName, "error", err)
			continue
		}
		slog.Info("Provider queried", "base", req.AssetSymbolA, "quote", req.AssetSymbolB, "price", price, "provider", m.ProviderName)
		resp.Price = price
		resp.ProviderName = m.ProviderName
		resp.ProviderColor = m.Color
		return &resp, nil
	}
	return &resp, errors.New(fmt.Sprintf("Could not get the price of %s in %s", req.AssetSymbolA, req.AssetSymbolB))
}

package managing

import (
	"controtto/src/domain/pnl"
	"errors"
	"sync"
	"time"
)

var (
	ErrMarketNotFound   = errors.New("market traders not found")
	ErrEmptyToken       = errors.New("market traders not configured")
	ErrInvalidTrade     = errors.New("invalid trade parameters")
	ErrMarketNotHealthy = errors.New("market API not healthy")
)

// MarketManager handles all market operations
type MarketManager struct {
	markets map[string]*pnl.Market
	mu      sync.RWMutex // protects the markets map
}

func NewMarketManager(in map[string]pnl.Market) *MarketManager {
	markets := make(map[string]*pnl.Market, len(in))
	for key, val := range in {
		v := val // copy value to avoid referencing the same instance
		markets[key] = &v
	}
	mm := &MarketManager{
		markets: markets,
	}
	for key, market := range markets {
		err := mm.UpdateMarket(key, market.Token)
		if err != nil {
			panic("could not initialize markets")
		}
	}
	return mm
}

func (c *MarketManager) GetMarkets(all bool) map[string]*pnl.Market {
	c.mu.RLock()
	defer c.mu.RUnlock()
	filtered := make(map[string]*pnl.Market)
	for k, v := range c.markets {
		if all || v.IsSet {
			filtered[k] = v
		}
	}
	return filtered
}

func (m *MarketManager) UpdateMarket(key string, token string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	market, ok := m.markets[key]
	if !ok {
		return ErrMarketNotFound
	}
	market.IsSet = token != ""
	market.Token = token // Clear token
	if market.IsSet {
		market.API = market.Init(token)
	}
	m.markets[key] = market
	return nil
}

// getMarket gets a market's instance
func (m *MarketManager) getMarket(key string) (*pnl.Market, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	market, ok := m.markets[key]
	if !ok {
		return nil, ErrMarketNotFound
	}
	if !market.IsSet {
		return nil, ErrEmptyToken
	}
	return market, nil
}

// ListTraders returns all configured markets
func (m *MarketManager) ListTraders(all bool) map[string]*pnl.Market {
	return m.GetMarkets(all) // true = only return configured markets
}

// ExecuteBuy executes a buy order
func (m *MarketManager) ExecuteBuy(marketKey string, options pnl.TradeOptions) (*pnl.Trade, error) {
	if options.Amount <= 0 {
		return nil, ErrInvalidTrade
	}
	market, err := m.getMarket(marketKey)
	if err != nil {
		return nil, err
	}
	if !market.API.HealthCheck() {
		return nil, ErrMarketNotHealthy
	}
	return market.API.Buy(options)
}

// ExecuteSell executes a sell order
func (m *MarketManager) ExecuteSell(marketKey string, options pnl.TradeOptions) (*pnl.Trade, error) {
	if options.Amount <= 0 {
		return nil, ErrInvalidTrade
	}
	market, err := m.getMarket(marketKey)
	if err != nil {
		return nil, err
	}
	if !market.API.HealthCheck() {
		return nil, ErrMarketNotHealthy
	}
	return market.API.Sell(options)
}

// FetchBalance gets asset balance from a market
func (m *MarketManager) FetchBalance(marketKey, symbol string) (float64, error) {
	market, err := m.getMarket(marketKey)
	if err != nil {
		return 0, err
	}
	return market.API.FetchAssetAmount(symbol)
}

// ImportTrades imports historical trades
func (m *MarketManager) ImportTrades(marketKey string, pair pnl.TradingPair, since time.Time) ([]pnl.Trade, error) {
	market, err := m.getMarket(marketKey)
	if err != nil {
		return nil, err
	}
	return market.API.ImportTrades(pair, since)
}

// CheckHealth verifies market connection
func (m *MarketManager) CheckHealth(marketKey string) (bool, error) {
	market, err := m.getMarket(marketKey)
	if err != nil {
		return false, err
	}
	return market.API.HealthCheck(), nil
}

// GetAccountDetails fetches account information
func (m *MarketManager) GetAccountDetails(marketKey string) (string, error) {
	market, err := m.getMarket(marketKey)
	if err != nil {
		return "", err
	}
	return market.API.AccountDetails()
}

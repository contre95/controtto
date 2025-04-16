package managing

import (
	"controtto/src/domain/pnl"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrTraderNotFound   = errors.New("market traders not found")
	ErrEmptyToken       = errors.New("market traders not configured")
	ErrInvalidTrade     = errors.New("invalid trade parameters")
	ErrMarketNotHealthy = errors.New("market API not healthy")
)

// MarketManager handles all market operations
type MarketManager struct {
	traders map[string]*pnl.MarketTrader
	mu      sync.RWMutex // protects the traders map
}

func NewMarketManager(in map[string]pnl.MarketTrader) *MarketManager {
	traders := make(map[string]*pnl.MarketTrader, len(in))
	for key, val := range in {
		v := val // copy value to avoid referencing the same instance
		traders[key] = &v
	}
	mm := &MarketManager{
		traders: traders,
	}
	for key, trader := range traders {
		err := mm.UpdateTrader(key, trader.Token)
		if err != nil {
			panic("could not initialize traders")
		}
	}
	return mm
}

func (c *MarketManager) UpdateMarketTraderToken(key, token string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	trader, ok := c.traders[key]
	if !ok {
		return fmt.Errorf("market trader %q not found", key)
	}
	trader.Token = token
	trader.IsSet = token != ""
	trader.Details = "Updated " + time.Now().Format(time.RFC3339)
	c.traders[key] = trader
	return nil
}

func (c *MarketManager) GetMarketTraders(all bool) map[string]*pnl.MarketTrader {
	c.mu.RLock()
	defer c.mu.RUnlock()
	filtered := make(map[string]*pnl.MarketTrader)
	for k, v := range c.traders {
		if all || v.IsSet {
			filtered[k] = v
		}
	}
	return filtered
}

func (m *MarketManager) UpdateTrader(key string, token string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	trader, ok := m.traders[key]
	if !ok {
		return ErrTraderNotFound
	}
	trader.IsSet = token != ""
	trader.Token = token // Clear token
	if trader.IsSet {
		trader.API = trader.Init(token)
	}
	m.traders[key] = trader
	return nil
}

// getTrader gets a trader's instance
func (m *MarketManager) getTrader(key string) (*pnl.MarketTrader, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	trader, ok := m.traders[key]
	if !ok {
		return nil, ErrTraderNotFound
	}
	if !trader.IsSet {
		return nil, ErrEmptyToken
	}
	return trader, nil
}

// ListTraders returns all configured traders
func (m *MarketManager) ListTraders(all bool) map[string]*pnl.MarketTrader {
	return m.GetMarketTraders(all) // true = only return configured traders
}

// ExecuteBuy executes a buy order
func (m *MarketManager) ExecuteBuy(marketKey string, options pnl.TradeOptions) (*pnl.Trade, error) {
	if options.Amount <= 0 {
		return nil, ErrInvalidTrade
	}
	trader, err := m.getTrader(marketKey)
	if err != nil {
		return nil, err
	}
	if !trader.API.HealthCheck() {
		return nil, ErrMarketNotHealthy
	}
	return trader.API.Buy(options)
}

// ExecuteSell executes a sell order
func (m *MarketManager) ExecuteSell(marketKey string, options pnl.TradeOptions) (*pnl.Trade, error) {
	if options.Amount <= 0 {
		return nil, ErrInvalidTrade
	}
	trader, err := m.getTrader(marketKey)
	if err != nil {
		return nil, err
	}
	if !trader.API.HealthCheck() {
		return nil, ErrMarketNotHealthy
	}
	return trader.API.Sell(options)
}

// FetchBalance gets asset balance from a market
func (m *MarketManager) FetchBalance(marketKey, symbol string) (float64, error) {
	trader, err := m.getTrader(marketKey)
	if err != nil {
		return 0, err
	}
	return trader.API.FetchAssetAmount(symbol)
}

// ImportTrades imports historical trades
func (m *MarketManager) ImportTrades(marketKey string, pair pnl.TradingPair, since time.Time) ([]pnl.Trade, error) {
	trader, err := m.getTrader(marketKey)
	if err != nil {
		return nil, err
	}
	return trader.API.ImportTrades(pair, since)
}

// CheckHealth verifies market connection
func (m *MarketManager) CheckHealth(marketKey string) (bool, error) {
	trader, err := m.getTrader(marketKey)
	if err != nil {
		return false, err
	}
	return trader.API.HealthCheck(), nil
}

// GetAccountDetails fetches account information
func (m *MarketManager) GetAccountDetails(marketKey string) (string, error) {
	trader, err := m.getTrader(marketKey)
	if err != nil {
		return "", err
	}
	return trader.API.AccountDetails()
}

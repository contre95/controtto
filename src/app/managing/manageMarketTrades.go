// application/market_manager.go
package managing

import (
	"controtto/src/app/config"
	"controtto/src/domain/pnl"
	"errors"
	"time"
)

var (
	ErrTraderNotFound      = errors.New("market trader not found")
	ErrTraderNotConfigured = errors.New("market trader not configured")
	ErrInvalidTrade        = errors.New("invalid trade parameters")
	ErrMarketNotHealthy    = errors.New("market API not healthy")
)

// MarketManager handles all market operations
type MarketManager struct {
	cfg *config.ConfigManager // Source of truth for trader configurations
}

func NewMarketManager(cfg *config.ConfigManager) *MarketManager {
	return &MarketManager{cfg: cfg}
}

// getTrader gets a fresh trader instance directly from config
func (m *MarketManager) getTrader(key string) (*pnl.MarketTrader, error) {
	traders := m.cfg.GetMarketTraders(true) // true = only return configured traders
	for k, trader := range traders {
		if k == key {
			if !trader.IsSet {
				return nil, ErrTraderNotConfigured
			}
			return &trader, nil
		}
	}
	return nil, ErrTraderNotFound
}

// ListTraders returns all configured traders
func (m *MarketManager) ListTraders(all bool) pnl.MarketTraders {
	return m.cfg.GetMarketTraders(all) // true = only return configured traders
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

	if !trader.MarketAPI.HealthCheck() {
		return nil, ErrMarketNotHealthy
	}

	return trader.MarketAPI.Buy(options)
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

	if !trader.MarketAPI.HealthCheck() {
		return nil, ErrMarketNotHealthy
	}

	return trader.MarketAPI.Sell(options)
}

// FetchBalance gets asset balance from a market
func (m *MarketManager) FetchBalance(marketKey, symbol string) (float64, error) {
	trader, err := m.getTrader(marketKey)
	if err != nil {
		return 0, err
	}
	return trader.MarketAPI.FetchAssetAmount(symbol)
}

// ImportTrades imports historical trades
func (m *MarketManager) ImportTrades(marketKey string, pair pnl.TradingPair, since time.Time) ([]pnl.Trade, error) {
	trader, err := m.getTrader(marketKey)
	if err != nil {
		return nil, err
	}
	return trader.MarketAPI.ImportTrades(pair, since)
}

// CheckHealth verifies market connection
func (m *MarketManager) CheckHealth(marketKey string) (bool, error) {
	trader, err := m.getTrader(marketKey)
	if err != nil {
		return false, err
	}
	return trader.MarketAPI.HealthCheck(), nil
}

// GetAccountDetails fetches account information
func (m *MarketManager) GetAccountDetails(marketKey string) (string, error) {
	trader, err := m.getTrader(marketKey)
	if err != nil {
		return "", err
	}
	return trader.MarketAPI.AccountDetails()
}

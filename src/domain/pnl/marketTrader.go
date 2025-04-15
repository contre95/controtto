package pnl

import "time"

type MarketType string

const (
	Exchange MarketType = "Exchange"
	Broker   MarketType = "Broker"
	DEX      MarketType = "DEX"
)

type MarketTrader struct {
	IsSet       bool
	Color       string
	MarketName  string
	MarketKey   string
	Type        MarketType
	Details     string
	MarketLogo  string
	Token       string
	ProviderURL string
	Env         string
	Init        func(string) MarketAPI
	API         MarketAPI
}

type MarketTraders map[string]MarketTrader

type MarketTraderNotFound error

// MarketTrader repository interface
type TradeOptions struct {
	TradingPair   TradingPair
	Amount        float64
	Price         *float64 // Optional price for limit orders
	IsMarketOrder bool
	// Other order parameters like stop loss, take profit, etc.
}

type MarketAPI interface {
	// Fetches how much of this asset you have in the Market
	FetchAssetAmount(symbol string) (float64, error)
	AccountDetails() (string, error)
	HealthCheck() bool
	ImportTrades(tradingPair TradingPair, since time.Time) ([]Trade, error)
	Buy(options TradeOptions) (*Trade, error)
	Sell(options TradeOptions) (*Trade, error)
}

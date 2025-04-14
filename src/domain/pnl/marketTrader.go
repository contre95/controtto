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
	MarketLogo  string
	Token       string
	ProviderURL string
	Env         string
	MarketAPI
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
	AccountDetails() (string, error)
	Buy(options TradeOptions) (*Trade, error)
	Sell(options TradeOptions) (*Trade, error)
	// Fetches last transaction between Base and Quote on the Exchange and returns that in a Trade format
	ImportTrades(tradingPair TradingPair, since time.Time) ([]Trade, error)
	// Fetches amount of assets there are Spot (Exchanges)
	FetchAsset(symbol string) (float64, error)
}

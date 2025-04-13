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
	TradingPairID TradingPairID
	Amount        float64
	Price         *float64
	IsMarketOrder bool
	// Other order parameters like stop loss, take profit, etc.
}

type MarketAPI interface {
	Buy(options TradeOptions) (*Trade, error)
	Sell(options TradeOptions) (*Trade, error)
	ImportTrades(tradingPairID TradingPairID, since time.Time) ([]Trade, error)
	FetchAsset(symbol string) (float64, error)
}

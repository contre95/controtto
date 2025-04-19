package pnl

import "time"

type MarketType string

const (
	Exchange MarketType = "Exchange"
	Wallet   MarketType = "Wallet"
	Broker   MarketType = "Broker"
	DEX      MarketType = "DEX"
)

type Market struct {
	IsSet                bool
	Color                string
	MarketName           string
	MarketKey            string
	MarketTradingSymbols []string
	Type                 MarketType
	Details              string
	MarketLogo           string
	Token                string
	ProviderURL          string
	Env                  string
	Init                 func(string) MarketAPI
	API                  MarketAPI
}

type MarketTraderNotFound error

// MarketTrader repository interface
type TradeOptions struct {
	Pair          Pair
	Amount        float64
	Price         *float64 // Optional price for limit orders
	IsMarketOrder bool
	// Other order parameters like stop loss, take profit, etc.
}

type MarketAPI interface {
	// Fetches how much of this asset you have in the Market
	FetchAssetAmount(ticket string) (float64, error)
	// FetchPrice(ticket string) (float64, error)
	AccountDetails() (string, error)
	HealthCheck() bool
	ImportTrades(tradingPair Pair, since time.Time) ([]Trade, error)
	Buy(options TradeOptions) (*Trade, error)
	Sell(options TradeOptions) (*Trade, error)
}

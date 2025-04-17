package pnl

import "time"

type PairID string

// Pair represents the primary aggregate root. It contains the main context for profit and loss calculations between two assets
type Pair struct {
	ID           PairID
	BaseAsset    Asset
	QuoteAsset   Asset
	Trades       []Trade
	Calculations Calculations
}

// Calculation is a value object for a Pair and it is populated with the function Calculate. It hold the data inferred from the Pair trades
type Calculations struct {
	AvgBuyPrice              float64
	BasePrice                float64
	CurrentBaseAmountInQuote float64
	TotalBase                float64
	TotalBaseInQuote         float64
	TotalQuoteSpent          float64
	PNLAmount                float64
	PNLPercent               float64
	TotalFeeInQuote          float64
	TotalFeeInBase           float64
}

const (
	Buy TradeType = "Buy"
	// Withdraw TradeType = "Withdraw"
	Sell TradeType = "Sell"
)

type TradeType string

func GetValidTradeTypes() []TradeType {
	return []TradeType{Buy, Sell}
}

// Trade represents individual exchange trades between the asset pair.
// The first listed currency of a currency pair is called the base currency, and the second currency is called the quote currency.
type Trade struct {
	ID          string
	Timestamp   time.Time
	BaseAmount  float64
	QuoteAmount float64
	FeeInBase   float64
	FeeInQuote  float64
	TradeType   TradeType
	Price       float64
}

// TrasingPairs repository interface
type Pairs interface {
	AddPair(tp Pair) error
	ListPairs() ([]Pair, error)
	GetPair(tpid string) (*Pair, error)
	DeletePair(tpid string) error
	DeleteTrade(tid string) error
	RecordTrade(t Trade, tpid PairID) error
	ListTrades(tpid string) ([]Trade, error)
}

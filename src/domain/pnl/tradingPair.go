package pnl

import "time"

type TradingPairID string

// TradingPair represents the primary aggregate root. It contains the main context for profit and loss calculations between two assets
type TradingPair struct {
	ID           TradingPairID
	BaseAsset    Asset
	QuoteAsset   Asset
	Transactions []Transaction
	Calculations Calculations
}

// Calculation is a value object for a TradingPair and it is populated with the function Calculate. It hold the data inferred from the TradingPair transactions
type Calculations struct {
	AvgBuyPrice              float64
	BaseMarketPrice          float64
	MarketName               string
	MarketColor              string
	CurrentBaseAmountInQuote float64
	TotalBase                float64
	TotalQuoteSpent          float64
	PNLAmount                float64
	PNLPercent               float64
	TotalFeeInQuote          float64
	TotalFeeInBase           float64
}

const (
	Buy TransactionType = "Buy"
	// Withdraw TransactionType = "Withdraw"
	Sell TransactionType = "Sell"
)

type TransactionType string

func GetValidTransactionTypes() []TransactionType {
	return []TransactionType{Buy, Sell}
}

// Transaction represents individual exchange transactions between the asset pair.
// The first listed currency of a currency pair is called the base currency, and the second currency is called the quote currency.
type Transaction struct {
	ID              string
	Timestamp       time.Time
	BaseAmount      float64
	QuoteAmount     float64
	FeeInBase       float64
	FeeInQuote      float64
	TransactionType TransactionType
	Price           float64
}

// TrasingPairs repository interface
type TradingPairs interface {
	AddTradingPair(tp TradingPair) error
	ListTradingPairs() ([]TradingPair, error)
	GetTradingPair(tpid string) (*TradingPair, error)
	DeleteTradingPair(tpid string) error
	DeleteTransaction(tid string) error
	RecordTransaction(t Transaction, tpid TradingPairID) error
	ListTransactions(tpid string) ([]Transaction, error)
}

type MarketNotFound error

// Markets repository interface
type Markets interface {
	// GetCurrentPrice returns the given price of assetA expressed in terms of assetB, if the value is market is not found it returns a MarketNotFound error
	GetCurrentPrice(assetA, assetB string) (float64, error)
	Color() string
	Name() string
}

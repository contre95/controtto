package pnl

import (
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/google/uuid"
)

type TradingPairID string

// TradingPair represents the primary aggregate root. It contains the main context for profit and loss calculations between two assets
type TradingPair struct {
	ID              TradingPairID
	BaseAsset       Asset
	QuoteAsset      Asset
	Transactions    []Transaction
	AvgBuyPrice     float64
	TotalBase       float64
	TotalQuoteSpent float64
}

func NewTradingPair(base Asset, quote Asset) (*TradingPair, error) {
	// I'll leave uniqueness of this to an exception https://stackoverflow.com/questions/2660817/ddd-validation-of-unique-constraint
	tp := TradingPair{
		ID:           TradingPairID(fmt.Sprintf("%s%s", base.Symbol, quote.Symbol)),
		BaseAsset:    base,
		QuoteAsset:   quote,
		Transactions: []Transaction{},
	}
	return &tp, nil
}

func (tp *TradingPair) Validate() error {
	// TODO: Implement
	return nil
}

func (tp *TradingPair) CalculateFields() {
	// Perform any necessary validation or business logic checks here.
	tp.TotalBase = 0
	tp.TotalQuoteSpent = 0
	for _, t := range tp.Transactions {
		if t.TransactionType == Buy {
			tp.TotalBase += t.BaseAmount
			tp.TotalQuoteSpent += t.QuoteAmount
		}
		if t.TransactionType == Sell {
			tp.TotalBase -= t.BaseAmount
			tp.TotalQuoteSpent -= t.QuoteAmount
		}
	}
	tp.AvgBuyPrice = float64(tp.TotalQuoteSpent / tp.TotalBase)
	slog.Info("Fields calculated", "base", tp.TotalBase, "quote", tp.TotalQuoteSpent, "avg-buy-price", tp.AvgBuyPrice)
}

func (tp *TradingPair) NewTransaction(baseAmount float64, quoteAmount float64, timestamp time.Time, tType TransactionType) (*Transaction, error) {
	// Perform any necessary validation or business logic checks here.

	if !slices.Contains[[]TransactionType]([]TransactionType{Buy, Sell}, tType) {
		return nil, fmt.Errorf("Invalid Transaction type '%s'.", tType)
	}
	// Append the transaction to the Transactions slice.
	transaction := &Transaction{
		ID:              uuid.New().String(),
		Timestamp:       timestamp,
		BaseAmount:      baseAmount,
		QuoteAmount:     quoteAmount,
		TransactionType: tType,
	}
	// transaction.CalculateFields()
	tp.Transactions = append(tp.Transactions, *transaction)

	// Update any other relevant state of the account, like profit/loss calculations.
	return transaction, nil
}

type TransactionType string

const (
	Buy  TransactionType = "Buy"
	Sell TransactionType = "Sell"
)

// Transaction represents individual exchange transactions between the asset pair.
// The first listed currency of a currency pair is called the base currency, and the second currency is called the quote currency.
type Transaction struct {
	ID              string
	Timestamp       time.Time
	BaseAmount      float64
	QuoteAmount     float64
	TradingFee      float64
	WithdrawalFee   float64
	TransactionType TransactionType
	Price           float64
}

func (t *Transaction) CalculateFields() error {
	t.Price = t.QuoteAmount / t.BaseAmount
	return nil
}

type TradingPairs interface {
	AddTradingPair(tp TradingPair) error
	ListTradingPairs() ([]TradingPair, error)
	GetTradingPair(tpid string) (*TradingPair, error)
	DeleteTradingPair(tpid string) error
	DeleteTransaction(tid string) error
	RecordTransaction(t Transaction, tpid TradingPairID) error
	ListTransactions(tpid string) ([]Transaction, error)
}

// Markets
type Markets interface {
	GetCurrentPrice(assetA, assetB string) (float64, error)
}

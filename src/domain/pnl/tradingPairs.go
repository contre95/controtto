package pnl

import (
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/google/uuid"
)

// NewTradingPair creates a new TradingPair validating its invariants.
func NewTradingPair(base Asset, quote Asset) (*TradingPair, error) {
	// I'll leave uniqueness of this to an exception https://stackoverflow.com/questions/2660817/ddd-validation-of-unique-constraint
	tp := TradingPair{
		ID:           TradingPairID(fmt.Sprintf("%s%s", base.Symbol, quote.Symbol)),
		BaseAsset:    base,
		QuoteAsset:   quote,
		Transactions: []Transaction{},
	}
	return tp.Validate()
}

// NewTransaction creates new transaction for the given TradingPair
func (tp *TradingPair) NewTransaction(baseAmount, quoteAmount, tFee, wFee float64, timestamp time.Time, tType TransactionType) (*Transaction, error) {
	// Append the transaction to the Transactions slice.
	transaction := &Transaction{
		ID:              uuid.New().String(),
		Timestamp:       timestamp,
		BaseAmount:      baseAmount,
		QuoteAmount:     quoteAmount,
		FeeInBase:       tFee,
		FeeInQuote:      wFee,
		TransactionType: tType,
	}
	tp.Transactions = append(tp.Transactions, *transaction)
	return transaction.Validate()
}

type InvalidTransaction error

// Validate validates a TradingPair, if all fields are valid it returns itself, otherwise it returns an InvalidTradingPair error.
func (t *Transaction) Validate() (*Transaction, error) {
	// Perform any necessary validation or business logic checks here.
	if t.FeeInBase > 0 && t.FeeInQuote > 0 {
		slog.Error("Transaction Validation error", "error", "Invalid fee, can't have both on a single transaction.", "FeeInQuote", t.FeeInQuote, "FeeInBase", t.FeeInBase)
		return nil, InvalidTransaction(errors.New("Invalid base/quote amounts"))
	}
	if t.BaseAmount <= 0 || t.QuoteAmount <= 0 {
		slog.Error("Transaction Validation error", "error", "Invalid base/quote amount")
		return nil, InvalidTransaction(errors.New("Invalid base/quote amounts"))
	}
	if !slices.Contains(t.TransactionType.GetValidTypes(), t.TransactionType) {
		slog.Error("Transaction Validation error", "error", "Invalid transaction type")
		return nil, InvalidTransaction(errors.New("Invalid transaction"))
	}
	return t, nil
}

// Validate validates a TradingPair, if all fields are valid it returns itself, otherwise it returns an InvalidTradingPair error.
func (tp *TradingPair) Validate() (*TradingPair, error) {
	// Perform any necessary validation or business logic checks here.
	return tp, nil
}

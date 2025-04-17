package pnl

import (
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/google/uuid"
)

// NewTradingPair creates a new Pair validating its invariants.
func NewTradingPair(base Asset, quote Asset) (*Pair, error) {
	// I'll leave uniqueness of this to an exception https://stackoverflow.com/questions/2660817/ddd-validation-of-unique-constraint
	tp := Pair{
		ID:         TradingPairID(fmt.Sprintf("%s%s", base.Symbol, quote.Symbol)),
		BaseAsset:  base,
		QuoteAsset: quote,
		Trades:     []Trade{},
	}
	return tp.Validate()
}

// NewTrade creates new trade for the given Pair
func (tp *Pair) NewTrade(baseAmount, quoteAmount, tFee, wFee float64, timestamp time.Time, tType TradeType) (*Trade, error) {
	// Append the trade to the Trades slice.
	trade := &Trade{
		ID:          uuid.New().String(),
		Timestamp:   timestamp,
		BaseAmount:  baseAmount,
		QuoteAmount: quoteAmount,
		FeeInBase:   tFee,
		FeeInQuote:  wFee,
		TradeType:   tType,
	}
	tp.Trades = append(tp.Trades, *trade)
	return trade.Validate()
}

type InvalidTrade error

// Validate validates a Pair, if all fields are valid it returns itself, otherwise it returns an InvalidTradingPair error.
func (t *Trade) Validate() (*Trade, error) {
	// Perform any necessary validation or business logic checks here.
	if t.FeeInBase > 0 && t.FeeInQuote > 0 {
		slog.Error("Trade Validation error", "error", "Invalid fee, can't have both on a single trade.", "FeeInQuote", t.FeeInQuote, "FeeInBase", t.FeeInBase)
		return nil, InvalidTrade(errors.New("Invalid base/quote amounts"))
	}
	if t.BaseAmount <= 0 || t.QuoteAmount <= 0 {
		slog.Error("Trade Validation error", "error", "Invalid base/quote amount")
		return nil, InvalidTrade(errors.New("Invalid base/quote amounts"))
	}
	if !slices.Contains(GetValidTradeTypes(), t.TradeType) {
		slog.Error("Trade Validation error", "error", "Invalid trade type")
		return nil, InvalidTrade(errors.New("Invalid trade"))
	}
	return t, nil
}

// Validate validates a Pair, if all fields are valid it returns itself, otherwise it returns an InvalidTradingPair error.
func (tp *Pair) Validate() (*Pair, error) {
	// Perform any necessary validation or business logic checks here.
	return tp, nil
}

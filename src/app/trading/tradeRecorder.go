package trading

import (
	"controtto/src/domain/pnl"
	"fmt"
	"log/slog"
	"time"
)

// TradeRecorder handles trade recording operations
type TradeRecorder struct {
	tradingPairs pnl.Pairs
}

// NewTradeRecorder creates a new TradeRecorder instance
func NewTradeRecorder(tradingPairsRepo pnl.Pairs) *TradeRecorder {
	return &TradeRecorder{
		tradingPairs: tradingPairsRepo,
	}
}

type DeleteTradeReq struct {
	ID string `json:"id"`
}

type DeleteTradeResp struct {
	ID  string `json:"id"`
	Msg string `json:"message"`
}

func (tpm *TradeRecorder) DeleteTrade(req DeleteTradeReq) (*DeleteTradeResp, error) {
	err := tpm.tradingPairs.DeleteTrade(req.ID)
	if err != nil {
		slog.Error("Error deleting trade",
			"error", err,
			"trade_id", req.ID)
		return nil, fmt.Errorf("failed to delete trade: %w", err)
	}

	return &DeleteTradeResp{
		ID:  req.ID,
		Msg: fmt.Sprintf("Trade %s deleted successfully", req.ID),
	}, nil
}

type RecordTradeReq struct {
	PairID      string    `json:"trading_pair_id"`
	Timestamp   time.Time `json:"timestamp"`
	BaseAmount  float64   `json:"base_amount"`
	QuoteAmount float64   `json:"quote_amount"`
	FeeInBase   float64   `json:"fee_in_base"`
	FeeInQuote  float64   `json:"fee_in_quote"`
	Type        string    `json:"type"`
}

type RecordTradeResp struct {
	ID         string    `json:"id"`
	Msg        string    `json:"message"`
	RecordTime time.Time `json:"record_time"`
}

func (tpm *TradeRecorder) RecordTrade(req RecordTradeReq) (*RecordTradeResp, error) {
	// Validate trade type
	if !isValidTradeType(req.Type) {
		return nil, fmt.Errorf("invalid trade type: %s", req.Type)
	}
	fmt.Println("RecordTradeReq", req)
	tradingPair, err := tpm.tradingPairs.GetPair(req.PairID)
	if err != nil {
		slog.Error("Failed to get trading pair",
			"error", err,
			"trading_pair_id", req.PairID)
		return nil, fmt.Errorf("failed to get trading pair: %w", err)
	}

	// Create new trade
	trade, err := tradingPair.NewTrade(
		req.BaseAmount,
		req.QuoteAmount,
		req.FeeInBase,
		req.FeeInQuote,
		req.Timestamp,
		pnl.TradeType(req.Type),
	)
	if err != nil {
		slog.Error("Failed to create trade",
			"error", err,
			"trading_pair_id", req.PairID)
		return nil, fmt.Errorf("failed to create trade: %w", err)
	}

	// Persist trade
	if err := tpm.tradingPairs.RecordTrade(*trade, tradingPair.ID); err != nil {
		slog.Error("Failed to record trade",
			"error", err,
			"trade_id", trade.ID)
		return nil, fmt.Errorf("failed to record trade: %w", err)
	}

	slog.Info("Trade recorded successfully",
		"trade_id", trade.ID,
		"trading_pair_id", tradingPair.ID,
		"timestamp", req.Timestamp)

	return &RecordTradeResp{
		ID:         trade.ID,
		Msg:        "Trade recorded successfully",
		RecordTime: time.Now().UTC(),
	}, nil
}

func isValidTradeType(tType string) bool {
	for _, validType := range pnl.GetValidTradeTypes() {
		if string(validType) == tType {
			return true
		}
	}
	return false
}

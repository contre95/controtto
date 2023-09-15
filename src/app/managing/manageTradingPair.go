package managing

import (
	"controtto/src/domain/pnl"
	"fmt"
	"log/slog"
	"time"
)

type RecordTransactionReq struct {
	TradingPairID string
	Timestamp     time.Time
	BaseAmount    float64
	QuoteAmount   float64
	FeeInBase     float64
	FeeInQuote    float64
	Type          string
}

type RecordTransactionResp struct {
	ID         string
	Msg        string
	RecordTime time.Time
}

type CreateTradingPairResp struct {
	ID  string
	Msg string
}

type DeleteTradingPairReq struct {
	ID string
}

type DeleteTradingPairResp struct {
	ID  string
	Msg string
}

type DeleteTransactionReq struct {
	ID string
}

type DeleteTransactionResp struct {
	ID  string
	Msg string
}

type CreateTradingPairReq struct {
	BaseAssetSymbol  string
	QuoteAssetSymbol string
}

type TradingPairsManager struct {
	assets       pnl.Assets
	tradingPairs pnl.TradingPairs
}

func NewTradingPairManager(a pnl.Assets, tp pnl.TradingPairs) *TradingPairsManager {
	return &TradingPairsManager{a, tp}
}

func (tpm *TradingPairsManager) DeleteTradingPair(req DeleteTradingPairReq) (*DeleteTradingPairResp, error) {
	err := tpm.tradingPairs.DeleteTradingPair(req.ID)
	if err != nil {
		slog.Error("Error deleting trading pair", "error", err)
		return nil, err
	}
	return &DeleteTradingPairResp{
		ID:  req.ID,
		Msg: fmt.Sprintf("Transaction %s deleted successfully!", req.ID),
	}, nil
}

func (tpm *TradingPairsManager) DeleteTransaction(req DeleteTransactionReq) (*DeleteTransactionResp, error) {
	err := tpm.tradingPairs.DeleteTransaction(req.ID)
	if err != nil {
		slog.Error("Error deleting trading pair", "error", err)
		return nil, err
	}
	return &DeleteTransactionResp{
		ID:  req.ID,
		Msg: fmt.Sprintf("%s Transaction deleted successfully!", req.ID),
	}, nil
}

func (tpm *TradingPairsManager) RecordTransaction(req RecordTransactionReq) (*RecordTransactionResp, error) {
	var err error
	tradingPair, err := tpm.tradingPairs.GetTradingPair(string(req.TradingPairID))
	if err != nil {
		slog.Error("Could not retrieve TradingPair", "error", err)
		return nil, err
	}
	transaction, err := tradingPair.NewTransaction(req.BaseAmount, req.QuoteAmount, req.FeeInBase, req.FeeInQuote, req.Timestamp, pnl.TransactionType(req.Type))
	if err != nil {
		slog.Error("Could create transaction", "error", err)
		return nil, err
	}
	err = tpm.tradingPairs.RecordTransaction(*transaction, tradingPair.ID)
	if err != nil {
		slog.Error("Could not persist transaction", "error", err)
		return nil, err
	}
	slog.Info("Transaction created successfully.", "time", req.Timestamp.Format(time.UnixDate))
	return &RecordTransactionResp{
		ID:         transaction.ID,
		Msg:        "Transaction created successfully",
		RecordTime: time.Now(),
	}, nil
}

func (tpm *TradingPairsManager) Create(req CreateTradingPairReq) (*CreateTradingPairResp, error) {
	var err error
	var base, quote *pnl.Asset
	base, err = tpm.assets.GetAsset(req.BaseAssetSymbol)
	if err != nil {
		slog.Error("Getting asset", "Asset", req.BaseAssetSymbol, "error", err)
		return nil, err
	}
	quote, err = tpm.assets.GetAsset(req.QuoteAssetSymbol)
	if err != nil {
		slog.Error("Getting asset", "Asset", req.QuoteAssetSymbol, "error", err)
		return nil, err
	}
	tradingPair, err := pnl.NewTradingPair(*base, *quote)
	if err != nil {
		slog.Error("Creating new pair", "Base", base.Name, "Quote", quote.Name, "error", err)
		return nil, err
	}
	err = tpm.tradingPairs.AddTradingPair(*tradingPair)
	// TODO: Validate already exists asset
	if err != nil {
		slog.Error("Persisting trading pair", "Trading pair", tradingPair, "error", err)
		return nil, fmt.Errorf("Coudl not create pair ( %s/%s )", tradingPair.BaseAsset.Symbol, tradingPair.QuoteAsset.Symbol)
	}
	slog.Info("New Trading pair created", "TradingPairs", tradingPair)
	return &CreateTradingPairResp{
		ID:  string(tradingPair.ID),
		Msg: fmt.Sprintf("Trading %s pair created successfully.", string(tradingPair.ID)),
	}, nil
}

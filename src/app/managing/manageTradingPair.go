package managing

import (
	"controtto/src/app/config"
	"controtto/src/domain/pnl"
	"fmt"
	"log/slog"
	"time"
)

type RecordTradeReq struct {
	TradingPairID string
	Timestamp     time.Time
	BaseAmount    float64
	QuoteAmount   float64
	FeeInBase     float64
	FeeInQuote    float64
	Type          string
}

type RecordTradeResp struct {
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

type DeleteTradeReq struct {
	ID string
}

type DeleteTradeResp struct {
	ID  string
	Msg string
}

type CreateTradingPairReq struct {
	BaseAssetSymbol  string
	QuoteAssetSymbol string
}

type TradingPairsManager struct {
	config       *config.Config
	assets       pnl.Assets
	tradingPairs pnl.TradingPairs
}

func NewTradingPairManager(cfg *config.Config, a pnl.Assets, tp pnl.TradingPairs) *TradingPairsManager {
	return &TradingPairsManager{cfg, a, tp}
}

func (tpm *TradingPairsManager) DeleteTradingPair(req DeleteTradingPairReq) (*DeleteTradingPairResp, error) {
	err := tpm.tradingPairs.DeleteTradingPair(req.ID)
	if err != nil {
		slog.Error("Error deleting trading pair", "error", err)
		return nil, err
	}
	return &DeleteTradingPairResp{
		ID:  req.ID,
		Msg: fmt.Sprintf("Trade %s deleted successfully!", req.ID),
	}, nil
}

func (tpm *TradingPairsManager) DeleteTrade(req DeleteTradeReq) (*DeleteTradeResp, error) {
	err := tpm.tradingPairs.DeleteTrade(req.ID)
	if err != nil {
		slog.Error("Error deleting trading pair", "error", err)
		return nil, err
	}
	return &DeleteTradeResp{
		ID:  req.ID,
		Msg: fmt.Sprintf("%s Trade deleted successfully!", req.ID),
	}, nil
}

func (tpm *TradingPairsManager) RecordTrade(req RecordTradeReq) (*RecordTradeResp, error) {
	var err error
	tradingPair, err := tpm.tradingPairs.GetTradingPair(string(req.TradingPairID))
	if err != nil {
		slog.Error("Could not retrieve TradingPair", "error", err)
		return nil, err
	}
	trade, err := tradingPair.NewTrade(req.BaseAmount, req.QuoteAmount, req.FeeInBase, req.FeeInQuote, req.Timestamp, pnl.TradeType(req.Type))
	if err != nil {
		slog.Error("Could create trade", "error", err)
		return nil, err
	}
	err = tpm.tradingPairs.RecordTrade(*trade, tradingPair.ID)
	if err != nil {
		slog.Error("Could not persist trade", "error", err)
		return nil, err
	}
	slog.Info("Trade created successfully.", "time", req.Timestamp.Format(time.UnixDate))
	return &RecordTradeResp{
		ID:         trade.ID,
		Msg:        "Trade created successfully",
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

	// Uncommon pairs validation
	if !config.Load().GetUncommonPairs() {
		bt := base.Type // Assuming these are strings like "Crypto", "Forex", etc.
		qt := quote.Type

		valid := (bt == "Crypto" && qt == "Crypto") ||
			(bt == "Crypto" && qt == "Forex") ||
			(bt == "Stock" && qt == "Forex") ||
			(bt == "Forex" && qt == "Forex")

		if !valid {
			err = fmt.Errorf("pair %s/%s not allowed with uncommonPairs enabled", bt, qt)
			slog.Error("Invalid trading pair type", "BaseType", bt, "QuoteType", qt, "error", err)
			return nil, err
		}
	}

	tradingPair, err := pnl.NewTradingPair(*base, *quote)
	if err != nil {
		slog.Error("Creating new pair", "Base", base.Name, "Quote", quote.Name, "error", err)
		return nil, err
	}

	err = tpm.tradingPairs.AddTradingPair(*tradingPair)
	if err != nil {
		slog.Error("Persisting trading pair", "Trading pair", tradingPair, "error", err)
		return nil, fmt.Errorf("Could not create pair ( %s/%s )", tradingPair.BaseAsset.Symbol, tradingPair.QuoteAsset.Symbol)
	}

	slog.Info("New Trading pair created", "TradingPairs", tradingPair)

	return &CreateTradingPairResp{
		ID:  string(tradingPair.ID),
		Msg: fmt.Sprintf("Trading %s pair created successfully.", string(tradingPair.ID)),
	}, nil
}

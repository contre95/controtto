package managing

import (
	"controtto/src/app/config"
	"controtto/src/domain/pnl"
	"fmt"
	"log/slog"
)

type CreatePairResp struct {
	ID  string
	Msg string
}

type DeletePairReq struct {
	ID string
}

type DeletePairResp struct {
	ID  string
	Msg string
}

type CreatePairReq struct {
	BaseAssetSymbol  string
	QuoteAssetSymbol string
}

type PairsManager struct {
	config       *config.Manager
	assets       pnl.Assets
	tradingPairs pnl.Pairs
}

func NewPairManager(cfg *config.Manager, a pnl.Assets, tp pnl.Pairs) *PairsManager {
	return &PairsManager{cfg, a, tp}
}

func (tpm *PairsManager) DeletePair(req DeletePairReq) (*DeletePairResp, error) {
	err := tpm.tradingPairs.DeletePair(req.ID)
	if err != nil {
		slog.Error("Error deleting trading pair", "error", err)
		return nil, err
	}
	return &DeletePairResp{
		ID:  req.ID,
		Msg: fmt.Sprintf("Trade %s deleted successfully!", req.ID),
	}, nil
}

func (tpm *PairsManager) Create(req CreatePairReq) (*CreatePairResp, error) {
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

	// // Uncommon pairs validation
	if !tpm.config.Get().UncommonPairs {
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

	tradingPair, err := pnl.NewPair(*base, *quote)
	if err != nil {
		slog.Error("Creating new pair", "Base", base.Name, "Quote", quote.Name, "error", err)
		return nil, err
	}

	err = tpm.tradingPairs.AddPair(*tradingPair)
	if err != nil {
		slog.Error("Persisting trading pair", "Trading pair", tradingPair, "error", err)
		return nil, fmt.Errorf("Could not create pair ( %s/%s )", tradingPair.BaseAsset.Symbol, tradingPair.QuoteAsset.Symbol)
	}

	slog.Info("New Trading pair created", "Pairs", tradingPair)

	return &CreatePairResp{
		ID:  string(tradingPair.ID),
		Msg: fmt.Sprintf("Trading %s pair created successfully.", string(tradingPair.ID)),
	}, nil
}

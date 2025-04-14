package querying

import (
	"controtto/src/app/config"
	"controtto/src/domain/pnl"
	"errors"
	"fmt"
	"log/slog"
)

type TradingPairsQuerier struct {
	tradingPairs pnl.TradingPairs
	cfg          *config.ConfigManager
}

// NewTradingPairQuerier returns a new intereactor with all the Trading Pair related use cases.
func NewTradingPairQuerier(cfg *config.ConfigManager, a pnl.TradingPairs) *TradingPairsQuerier {
	return &TradingPairsQuerier{a, cfg}
}

// List all trading pairs without any level of detail

type ListTradingPairsReq struct{}
type ListTradingPairsResp struct {
	Pairs []pnl.TradingPair
}

func (tpq *TradingPairsQuerier) ListTradingPairs(req ListTradingPairsReq) (*ListTradingPairsResp, error) {
	var err error
	pairs, err := tpq.tradingPairs.ListTradingPairs()
	if err != nil {
		slog.Error("Error getting tading pairs list from DB", "error", err)
		return nil, err
	}
	resp := ListTradingPairsResp{
		Pairs: pairs,
	}
	return &resp, nil
}

// List Trades

type TradesReq struct {
	TradingPairID string
}
type TradesResp struct {
	Trades []pnl.Trade
}

func (tpq *TradingPairsQuerier) ListTrades(req TradesReq) (*TradesResp, error) {
	var err error
	trades, err := tpq.getTrades(req.TradingPairID)
	if err != nil {
		slog.Error("Error getting list from DB", "TradingPair", req.TradingPairID, "error", err)
		return nil, err
	}
	slog.Error("Trades retrieved succesfully", "TradingPair", req.TradingPairID, "TradeCount", len(trades))
	resp := TradesResp{
		Trades: trades,
	}
	return &resp, nil
}

// Get single trading pair

// GetTradingPairReq indicate the level of datail you want to retrieve the Trading pair
type GetTradingPairReq struct {
	TPID                 string
	WithCurrentBasePrice bool
	WithTrades           bool
	WithCalculations     bool
}

type GetTradingPairResp struct {
	Pair pnl.TradingPair
}

func (tpq *TradingPairsQuerier) GetTradingPair(req GetTradingPairReq) (*GetTradingPairResp, error) {
	var err error
	pair, err := tpq.tradingPairs.GetTradingPair(req.TPID)
	if err != nil {
		return nil, err
	}

	if req.WithCalculations {
		req.WithTrades = true
		req.WithCurrentBasePrice = true
	}

	if req.WithTrades {
		trades, err := tpq.getTrades(req.TPID)
		if err != nil {
			return nil, err
		}
		for _, t := range trades {
			t.CalculateFields()
			pair.Trades = append(pair.Trades, t)
		}
	}

	if req.WithCurrentBasePrice {
		pair.Calculations.BasePrice, pair.Calculations.ProviderName, pair.Calculations.Color, _ = tpq.getCurrentBasePrice(pair.BaseAsset.Symbol, pair.QuoteAsset.Symbol)
	}

	if req.WithCalculations {
		err := pair.Calculate()
		if err != nil {
			slog.Error("Error making calculations", "TradinPair", string(pair.ID))
			return nil, errors.New("Error making calculations")
		}
	}

	resp := GetTradingPairResp{
		Pair: *pair,
	}
	return &resp, nil
}

func (tpq *TradingPairsQuerier) getTrades(tpid string) ([]pnl.Trade, error) {
	trades, err := tpq.tradingPairs.ListTrades(tpid)
	if err != nil {
		slog.Error("Error getting trade", "Trading Pair", tpid, "error", err)
		return nil, err
	}
	return trades, nil
}

func (tpq *TradingPairsQuerier) getCurrentBasePrice(asset1, asset2 string) (float64, string, string, error) {
	var err error
	var baseAssetPrice float64 = 0
	providerName := ""
	providerColor := ""
	failedproviders := 0
	providers := tpq.cfg.GetPriceProviders()
	for _, m := range providers {
		slog.Info("Querying providers", "provider", m.ProviderName)
		if m.IsSet {
			providerName = m.ProviderName
			providerColor = m.Color
			baseAssetPrice, err = m.GetCurrentPrice(asset1, asset2)
			if err != nil {
				slog.Error("Error getting base asset price.", "provider", m.ProviderName, "error", err)
				failedproviders++
			} else {
				break
			}
		} else {
			fmt.Println("not set")
		}
	}
	if failedproviders == len(providers) {
		slog.Error("All providers failed to find the price.", "asset1", asset1, "asset2", asset2)
		return 0, "", "", err
	}
	return baseAssetPrice, providerName, providerColor, nil
}

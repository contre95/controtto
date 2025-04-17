package querying

import (
	"controtto/src/domain/pnl"
	"errors"
	"log/slog"
)

type TradingPairsQuerier struct {
	tradingPairs pnl.Pairs
}

// NewTradingPairQuerier returns a new intereactor with all the Trading Pair related use cases.
func NewTradingPairQuerier(a pnl.Pairs) *TradingPairsQuerier {
	return &TradingPairsQuerier{a}
}

// List all trading pairs without any level of detail

type ListTradingPairsReq struct{}
type ListTradingPairsResp struct {
	Pairs []pnl.Pair
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
		slog.Error("Error getting list from DB", "Pair", req.TradingPairID, "error", err)
		return nil, err
	}
	slog.Error("Trades retrieved succesfully", "Pair", req.TradingPairID, "TradeCount", len(trades))
	resp := TradesResp{
		Trades: trades,
	}
	return &resp, nil
}

// Get single trading pair

// GetTradingPairReq indicate the level of datail you want to retrieve the Trading pair
type GetTradingPairReq struct {
	TPID             string
	WithTrades       bool
	WithCalculations bool
	BasePrice        float64
}

type GetTradingPairResp struct {
	Pair pnl.Pair
}

func (tpq *TradingPairsQuerier) GetTradingPair(req GetTradingPairReq) (*GetTradingPairResp, error) {
	var err error
	pair, err := tpq.tradingPairs.GetTradingPair(req.TPID)
	if err != nil {
		return nil, err
	}

	if req.WithCalculations {
		req.WithTrades = true
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

	pair.Calculations.BasePrice = req.BasePrice

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

package querying

import (
	"controtto/src/domain/pnl"
	"errors"
	"log/slog"
)

type PairsQuerier struct {
	tradingPairs pnl.Pairs
}

// NewPairQuerier returns a new intereactor with all the Trading Pair related use cases.
func NewPairQuerier(a pnl.Pairs) *PairsQuerier {
	return &PairsQuerier{a}
}

// List all trading pairs without any level of detail

type ListPairsReq struct{}
type ListPairsResp struct {
	Pairs []pnl.Pair
}

func (tpq *PairsQuerier) ListPairs(req ListPairsReq) (*ListPairsResp, error) {
	var err error
	pairs, err := tpq.tradingPairs.ListPairs()
	if err != nil {
		slog.Error("Error getting tading pairs list from DB", "error", err)
		return nil, err
	}
	resp := ListPairsResp{
		Pairs: pairs,
	}
	return &resp, nil
}

// List Trades

type TradesReq struct {
	PairID string
}
type TradesResp struct {
	Trades []pnl.Trade
}

func (tpq *PairsQuerier) ListTrades(req TradesReq) (*TradesResp, error) {
	var err error
	trades, err := tpq.getTrades(req.PairID)
	if err != nil {
		slog.Error("Error getting list from DB", "Pair", req.PairID, "error", err)
		return nil, err
	}
	slog.Error("Trades retrieved succesfully", "Pair", req.PairID, "TradeCount", len(trades))
	resp := TradesResp{
		Trades: trades,
	}
	return &resp, nil
}

// Get single trading pair

// GetPairReq indicate the level of datail you want to retrieve the Trading pair
type GetPairReq struct {
	TPID             string
	WithTrades       bool
	WithCalculations bool
	BasePrice        float64
}

type GetPairResp struct {
	Pair pnl.Pair
}

func (tpq *PairsQuerier) GetPair(req GetPairReq) (*GetPairResp, error) {
	var err error
	pair, err := tpq.tradingPairs.GetPair(req.TPID)
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

	resp := GetPairResp{
		Pair: *pair,
	}
	return &resp, nil
}

func (tpq *PairsQuerier) getTrades(tpid string) ([]pnl.Trade, error) {
	trades, err := tpq.tradingPairs.ListTrades(tpid)
	if err != nil {
		slog.Error("Error getting trade", "Trading Pair", tpid, "error", err)
		return nil, err
	}
	return trades, nil
}

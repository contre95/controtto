package querying

import (
	"controtto/src/domain/pnl"
	"log/slog"
)

type GetTradingPairReq struct {
	TPID string
}
type GetTradingPairResp struct {
	Pair pnl.TradingPair
}
type QueryTradingPairReq struct{}
type QueryTradingPairResp struct {
	Pairs []pnl.TradingPair
}

type TransactionsReq struct {
	TradingPairID string
}
type TransactionsResp struct {
	Transactions []pnl.Transaction
}

type TradingPairsQuerier struct {
	tradingPairs pnl.TradingPairs
}

func NewTradingPairQuerier(a pnl.TradingPairs) *TradingPairsQuerier {
	return &TradingPairsQuerier{a}
}

func (tpq *TradingPairsQuerier) ListTransactions(req TransactionsReq) (*TransactionsResp, error) {
	var err error
	transactions, err := tpq.tradingPairs.ListTransactions(req.TradingPairID)
	if err != nil {
		slog.Error("Error getting list from DB", "Trading Pair", req.TradingPairID, "error", err)
		return nil, err
	}
	resp := TransactionsResp{
		Transactions: transactions,
	}
	return &resp, nil
}

func (tpq *TradingPairsQuerier) ListTradingPairs(req QueryTradingPairReq) (*QueryTradingPairResp, error) {
	var err error
	pairs, err := tpq.tradingPairs.ListTradingPairs()
	if err != nil {
		slog.Error("Error getting tading pairs list from DB", "error", err)
		return nil, err
	}
	resp := QueryTradingPairResp{
		Pairs: pairs,
	}
	return &resp, nil
}

// GetTradingPair retrieves trading pair information with associated transactions.
func (tpq *TradingPairsQuerier) GetTradingPair(req GetTradingPairReq) (*GetTradingPairResp, error) {
	var err error
	pair, err := tpq.tradingPairs.GetTradingPair(req.TPID)
	if err != nil {
		slog.Error("Error getting tading pairs list from DB", "error", err)
		return nil, err
	}
	transactions, err := tpq.tradingPairs.ListTransactions(req.TPID)
	if err != nil {
		slog.Error("Error getting transaction", "Trading Pair", req.TPID, "error", err)
		return nil, err
	}
	for _, t := range transactions {
		t.CalculateFields()
		pair.Transactions = append(pair.Transactions, t)
	}
	resp := GetTradingPairResp{
		Pair: *pair,
	}
	return &resp, nil
}

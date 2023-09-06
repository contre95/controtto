package querying

import (
	"controtto/src/domain/pnl"
	"log/slog"
)

type GetTradingPairReq struct {
	TPID          string
	WithBasePrice bool
}
type GetTradingPairResp struct {
	Pair           pnl.TradingPair
	BaseAssetPrice float64
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
	markets      pnl.Markets
}

func NewTradingPairQuerier(a pnl.TradingPairs, m pnl.Markets) *TradingPairsQuerier {
	return &TradingPairsQuerier{a, m}
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

// GetTradingPair retrieves trading pair information with associated transactions. It also returns the current base asset value expressed in terms of the quote value
// If is fails to retrieve the value it will set it to 0 (zero)
//
// TODO: Answer: Why does the presenter won't simply use two different use cases to retrieve the market value ?. We'll see.
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

	baseAssetPrice := float64(0)
	if req.WithBasePrice {
		// baseAssetPrice is expres in the quoteAsset value, but ofcourse you know that cause you are a domain expert.
		// Is this they way of handling errors? What if I wanna notify the user that the actual price is no 0.
		baseAssetPrice, err = tpq.markets.GetCurrentPrice(pair.BaseAsset.Symbol, pair.QuoteAsset.Symbol)
		if err != nil {
			slog.Error("Error getting base asset price, setting it to 0.", "error", err)
		}
	}
	resp := GetTradingPairResp{
		Pair:           *pair,
		BaseAssetPrice: baseAssetPrice,
	}
	return &resp, nil
}

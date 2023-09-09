package querying

import (
	"controtto/src/domain/pnl"
	"log/slog"
)

type TradingPairsQuerier struct {
	tradingPairs pnl.TradingPairs
	markets      pnl.Markets
}

// NewTradingPairQuerier returns a new intereactor with all the Trading Pair related use cases.
func NewTradingPairQuerier(a pnl.TradingPairs, m pnl.Markets) *TradingPairsQuerier {
	return &TradingPairsQuerier{a, m}
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

// List Transactions

type TransactionsReq struct {
	TradingPairID string
}
type TransactionsResp struct {
	Transactions []pnl.Transaction
}

func (tpq *TradingPairsQuerier) ListTransactions(req TransactionsReq) (*TransactionsResp, error) {
	var err error
	transactions, err := tpq.getTransactions(req.TradingPairID)
	if err != nil {
		slog.Error("Error getting list from DB", "TradingPair", req.TradingPairID, "error", err)
		return nil, err
	}
	slog.Error("Transactions retrieved succesfully", "TradingPair", req.TradingPairID, "TransactionCount", len(transactions))
	resp := TransactionsResp{
		Transactions: transactions,
	}
	return &resp, nil
}

// Get single trading pair

// GetTradingPairReq indicate the level of datail you want to retrieve the Trading pair
type GetTradingPairReq struct {
	TPID                 string
	WithCurrentBasePrice bool
	WithTransactions     bool
	WithCalculations     bool
}

type GetTradingPairResp struct {
	Pair           pnl.TradingPair
	BaseAssetPrice float64
}

func (tpq *TradingPairsQuerier) GetTradingPair(req GetTradingPairReq) (*GetTradingPairResp, error) {
	var err error
	pair, err := tpq.tradingPairs.GetTradingPair(req.TPID)
	if err != nil {
		return nil, err
	}

	if req.WithCalculations {
		req.WithTransactions = true
		req.WithCurrentBasePrice = true
	}

	if req.WithTransactions {
		transactions, err := tpq.getTransactions(req.TPID)
		if err != nil {
			return nil, err
		}
		for _, t := range transactions {
			t.CalculateFields()
			pair.Transactions = append(pair.Transactions, t)
		}
	}

	if req.WithCurrentBasePrice {
		pair.Calculations.CurrentBasePrice, err = tpq.getCurrentBasePrice(pair.BaseAsset.Symbol, pair.QuoteAsset.Symbol)
		if err != nil {
			return nil, err
		}
	}

	if req.WithCalculations {
		err := pair.Calculate()
		if err != nil {
			return nil, err
		}
	}

	resp := GetTradingPairResp{
		Pair: *pair,
	}
	return &resp, nil
}

func (tpq *TradingPairsQuerier) getTransactions(tpid string) ([]pnl.Transaction, error) {
	transactions, err := tpq.tradingPairs.ListTransactions(tpid)
	if err != nil {
		slog.Error("Error getting transaction", "Trading Pair", tpid, "error", err)
		return nil, err
	}
	return transactions, nil
}

func (tpq *TradingPairsQuerier) getCurrentBasePrice(asset1, asset2 string) (float64, error) {
	baseAssetPrice, err := tpq.markets.GetCurrentPrice(asset1, asset2)
	if err != nil {
		slog.Error("Error getting base asset price, setting it to 0.", "error", err)
		baseAssetPrice = 0 // Set it to 0 in case of an error
	}
	return baseAssetPrice, nil
}

package querying

import (
	"controtto/src/domain/pnl"
	"log/slog"
)

type QueryMarketReq struct {
	Symbol string
}
type QueryMarketResp struct {
	Symbol string
	Price  float64
}

type MarketsQuerier struct {
	markets pnl.Markets
}

func NewMarketQuerier(a pnl.Markets) *MarketsQuerier {
	return &MarketsQuerier{a}
}

func (aq *MarketsQuerier) GetMarketPrice(req QueryMarketReq) (*QueryMarketResp, error) {
	var err error
	price, err := aq.markets.GetCurrentPrice(req.Symbol)
	if err != nil {
		slog.Error("Could not get market current price", "market", req.Symbol, "error", err)
		return nil, err
	}
	slog.Info("Market queried", "market", req.Symbol, "price", price)
	resp := QueryMarketResp{
		Symbol: req.Symbol,
		Price:  price,
	}
	return &resp, nil
}

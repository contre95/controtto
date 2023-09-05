package querying

import (
	"controtto/src/domain/pnl"
	"log/slog"
)

type QueryMarketReq struct {
	AssetSymbolA string
	AssetSymbolB string
}
type QueryMarketResp struct {
	Price float64
}

type MarketsQuerier struct {
	markets pnl.Markets
}

func NewMarketQuerier(a pnl.Markets) *MarketsQuerier {
	return &MarketsQuerier{a}
}

func (aq *MarketsQuerier) GetMarketPrice(req QueryMarketReq) (*QueryMarketResp, error) {
	var err error
	price, err := aq.markets.GetCurrentPrice(req.AssetSymbolA, req.AssetSymbolB)
	if err != nil {
		slog.Error("Could not get market current price", "market", req.AssetSymbolA, "error", err)
		return nil, err
	}
	slog.Info("Market queried", "market", req.AssetSymbolA, "price", price)
	resp := QueryMarketResp{
		Price: price,
	}
	return &resp, nil
}

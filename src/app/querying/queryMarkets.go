package querying

import (
	"controtto/src/domain/pnl"
	"errors"
	"fmt"
	"log/slog"
)

type QueryMarketReq struct {
	AssetSymbolA string
	AssetSymbolB string
}
type QueryMarketResp struct {
	Price         float64
	ProviderName  string
	ProviderColor string
}

type MarketsQuerier struct {
	markets []pnl.Markets
}

func NewMarketQuerier(a []pnl.Markets) *MarketsQuerier {
	return &MarketsQuerier{a}
}

// GetMarketPrice returns the current base asset value expressed in terms of the quote value
// If is fails to retrieve the value it will set it to 0 (zero)
func (aq *MarketsQuerier) GetMarketPrice(req QueryMarketReq) (*QueryMarketResp, error) {
	resp := QueryMarketResp{Price: 0}
	for _, m := range aq.markets {
		price, err := m.GetCurrentPrice(req.AssetSymbolA, req.AssetSymbolB)
		if err != nil {
			slog.Error("Could not get market current price", "base", req.AssetSymbolA, "quote", req.AssetSymbolB, "provider", m.Name(), "error", err)
			continue
		}
		slog.Info("Market queried", "base", req.AssetSymbolA, "quote", req.AssetSymbolB, "price", price, "provider", m.Name())
		resp.Price = price
		resp.ProviderName = m.Name()
		resp.ProviderColor = m.Color()
		return &resp, nil
	}
	return &resp, errors.New(fmt.Sprintf("Could not get the price of %s in %s", req.AssetSymbolA, req.AssetSymbolB))
}

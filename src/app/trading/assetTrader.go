package trading

import (
	"controtto/src/app/managing"
	"controtto/src/domain/pnl"
	"errors"
	"time"
)

var (
	ErrMarketNotHealthy = errors.New("market API not healthy")
	ErrInvalidTrade     = errors.New("invalid trade parameters")
)

// Request and Response DTOs
type TradeRequest struct {
	MarketKey     string
	PairID        string
	Amount        float64
	Price         *float64 // Optional price for limit orders
	IsMarketOrder bool
}

type TradeResponse struct {
	ID          string    `json:"id"`
	Timestamp   time.Time `json:"timestamp"`
	BaseAmount  float64   `json:"base_amount"`
	QuoteAmount float64   `json:"quote_amount"`
	FeeInBase   float64   `json:"fee_in_base"`
	FeeInQuote  float64   `json:"fee_in_quote"`
	TradeType   string    `json:"trade_type"`
	Price       float64   `json:"price"`
}

type FetchTradesReq struct {
	MarketKey string    `json:"market_key"`
	PairID    string    `json:"pair_id"`
	Since     time.Time `json:"since"`
}

type AssetTrader struct {
	markets *managing.MarketManager
	assets  pnl.Assets
	pairs   pnl.Pairs
}

func NewAssetTrader(mm *managing.MarketManager, p pnl.Pairs) *AssetTrader {
	return &AssetTrader{
		markets: mm,
		pairs:   p,
	}
}

func toTradeResponse(domain pnl.Trade) TradeResponse {
	return TradeResponse{
		ID:          domain.ID,
		Timestamp:   domain.Timestamp,
		BaseAmount:  domain.BaseAmount,
		QuoteAmount: domain.QuoteAmount,
		FeeInBase:   domain.FeeInBase,
		FeeInQuote:  domain.FeeInQuote,
		TradeType:   string(domain.TradeType),
		Price:       domain.Price,
	}
}

func (m *AssetTrader) getMarket(marketKey string) (*pnl.Market, error) {
	market, exists := m.markets.GetMarkets(false)[marketKey]
	if !exists {
		return nil, errors.New("market not found")
	}
	return market, nil
}

func (m *AssetTrader) getPair(pairID string) (*pnl.Pair, error) {
	pair, err := m.pairs.GetPair(pairID)
	if err != nil {
		return nil, err
	}
	return pair, nil
}

func (m *AssetTrader) ExecuteBuy(req TradeRequest) (*TradeResponse, error) {
	if req.Amount <= 0 {
		return nil, ErrInvalidTrade
	}

	market, err := m.getMarket(req.MarketKey)
	if err != nil {
		return nil, err
	}

	if !market.API.HealthCheck() {
		return nil, ErrMarketNotHealthy
	}

	pair, err := m.getPair(req.PairID)
	if err != nil {
		return nil, err
	}

	options := pnl.TradeOptions{
		Pair:          *pair,
		Amount:        req.Amount,
		Price:         req.Price,
		IsMarketOrder: req.IsMarketOrder,
	}

	trade, err := market.API.Buy(options)
	if err != nil {
		return nil, err
	}

	resp := toTradeResponse(*trade)
	return &resp, nil
}

func (m *AssetTrader) ExecuteSell(req TradeRequest) (*TradeResponse, error) {
	if req.Amount <= 0 {
		return nil, ErrInvalidTrade
	}

	market, err := m.getMarket(req.MarketKey)
	if err != nil {
		return nil, err
	}

	if !market.API.HealthCheck() {
		return nil, ErrMarketNotHealthy
	}

	pair, err := m.getPair(req.PairID)
	if err != nil {
		return nil, err
	}

	options := pnl.TradeOptions{
		Pair:          *pair,
		Amount:        req.Amount,
		Price:         req.Price,
		IsMarketOrder: req.IsMarketOrder,
	}

	trade, err := market.API.Sell(options)
	if err != nil {
		return nil, err
	}

	resp := toTradeResponse(*trade)
	return &resp, nil
}

func (m *AssetTrader) FetchTrades(req FetchTradesReq) ([]TradeResponse, error) {
	market, err := m.getMarket(req.MarketKey)
	if err != nil {
		return nil, err
	}

	pair, err := m.getPair(req.PairID)
	if err != nil {
		return nil, err
	}

	trades, err := market.API.ImportTrades(*pair, req.Since)
	if err != nil {
		return nil, err
	}

	var response []TradeResponse
	for _, trade := range trades {
		response = append(response, toTradeResponse(trade))
	}

	return response, nil
}

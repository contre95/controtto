package trading

import (
	"controtto/src/app/managing"
	"controtto/src/app/querying"
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
	Pair          PairDTO
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

type AssetDTO struct {
	Symbol      string `json:"symbol"`
	Color       string `json:"color"`
	Name        string `json:"name"`
	CountryCode string `json:"country_code"`
	Type        string `json:"type"`
}

type PairDTO struct {
	ID         string   `json:"id"`
	BaseAsset  AssetDTO `json:"base_asset"`
	QuoteAsset AssetDTO `json:"quote_asset"`
}

type ImportTradesRequest struct {
	MarketKey string    `json:"market_key"`
	Pair      PairDTO   `json:"pair"`
	Since     time.Time `json:"since"`
}

type AssetTrader struct {
	markets *managing.MarketManager
	assets  querying.AssetsQuerier
}

func NewAssetTrader(mm *managing.MarketManager, aq querying.AssetsQuerier) *AssetTrader {
	return &AssetTrader{
		markets: mm,
		assets:  aq,
	}
}

// Helper functions for domain <-> DTO conversion
func toDomainAsset(dto AssetDTO) (*pnl.Asset, error) {
	return pnl.NewAsset(
		dto.Symbol,
		dto.Color,
		dto.Name,
		dto.CountryCode,
		dto.Type,
	)
}

func toDomainPair(dto PairDTO) (*pnl.Pair, error) {
	baseAsset, err := toDomainAsset(dto.BaseAsset)
	if err != nil {
		return nil, err
	}

	quoteAsset, err := toDomainAsset(dto.QuoteAsset)
	if err != nil {
		return nil, err
	}

	return pnl.NewPair(*baseAsset, *quoteAsset)
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

	pair, err := toDomainPair(req.Pair)
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

	pair, err := toDomainPair(req.Pair)
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

func (m *AssetTrader) ImportTrades(req ImportTradesRequest) ([]TradeResponse, error) {
	market, err := m.getMarket(req.MarketKey)
	if err != nil {
		return nil, err
	}

	pair, err := toDomainPair(req.Pair)
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

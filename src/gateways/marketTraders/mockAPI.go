package marketTraders

import (
	"controtto/src/domain/pnl"
	"errors"
	"fmt"
	"math/rand/v2"
	"time"
)

type Pair string

type Asset struct {
	Symbol string
	Name   string
}

// MockMarketAPI is a mock implementation of MarketAPI
type MockMarketAPI struct {
	Token string
}

func (b *MockMarketAPI) HealthCheck() bool {
	return b.Token == "enable"
}

func (m *MockMarketAPI) Buy(options pnl.TradeOptions) (*pnl.Trade, error) {
	if options.Amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}
	price := 100.0
	if options.Price != nil {
		price = *options.Price
	}
	return &pnl.Trade{
		ID:          fmt.Sprintf("mock-trade-buy-%s", options.Pair),
		Timestamp:   time.Time{},
		BaseAmount:  options.Amount,
		QuoteAmount: options.Amount * price,
		FeeInBase:   0.05 * options.Amount,
		FeeInQuote:  0,
		TradeType:   "Buy",
		Price:       price,
	}, nil
}

func (m *MockMarketAPI) Sell(options pnl.TradeOptions) (*pnl.Trade, error) {
	if options.Amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}
	price := 100.0
	if options.Price != nil {
		price = *options.Price
	}
	return &pnl.Trade{
		ID:          fmt.Sprintf("mock-trade-buy-%s", options.Pair),
		Timestamp:   time.Time{},
		BaseAmount:  options.Amount,
		QuoteAmount: options.Amount * price,
		FeeInBase:   0.05 * options.Amount,
		FeeInQuote:  0,
		TradeType:   "Buy",
		Price:       price,
	}, nil
}

func (m *MockMarketAPI) ImportTrades(tradingPair pnl.Pair, since time.Time) ([]pnl.Trade, error) {
	return []pnl.Trade{
		{
			ID:          "trade-001",
			Timestamp:   time.Now().Add(-2 * time.Hour),
			BaseAmount:  0.5,
			QuoteAmount: 15000,
			FeeInBase:   0.0005,
			FeeInQuote:  15,
			TradeType:   "buy",
			Price:       30000,
		},
		{
			ID:          "trade-002",
			Timestamp:   time.Now().Add(-90 * time.Minute),
			BaseAmount:  0.2,
			QuoteAmount: 6000,
			FeeInBase:   0,
			FeeInQuote:  6,
			TradeType:   "sell",
			Price:       30000,
		},
		{
			ID:          "trade-003",
			Timestamp:   time.Now().Add(-45 * time.Minute),
			BaseAmount:  1.0,
			QuoteAmount: 29500,
			FeeInBase:   0.001,
			FeeInQuote:  0,
			TradeType:   "buy",
			Price:       29500,
		},
		{
			ID:          "trade-004",
			Timestamp:   time.Now().Add(-10 * time.Minute),
			BaseAmount:  0.75,
			QuoteAmount: 22500,
			FeeInBase:   0.00075,
			FeeInQuote:  10,
			TradeType:   "sell",
			Price:       30000,
		},
	}, nil
}

func (m *MockMarketAPI) AccountDetails() (string, error) {
	return "someeamil@domain.com", nil
}

func (m *MockMarketAPI) FetchAssetAmount(symbol string) (float64, error) {
	if symbol == "ETH" {
		return 0, nil
	}
	return rand.Float64() * 100, nil
}

func NewMockMarketAPI(token string) pnl.MarketAPI {
	return &MockMarketAPI{
		Token: token,
	}
}

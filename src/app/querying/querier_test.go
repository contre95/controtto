package querying

import (
	"controtto/src/domain/pnl"
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestListAssets(t *testing.T) {
	testAssets := []pnl.Asset{
		{
			Symbol:      "AAPL",
			Color:       "#FFFFFF",
			Name:        "Apple Stocks",
			CountryCode: "US",
		},
		{
			Symbol:      "GOOGL",
			Color:       "#FF5370",
			Name:        "Google Stocks",
			CountryCode: "US",
		},
		{
			Symbol:      "TSLA",
			Color:       "#0F111A",
			Name:        "Tesla Stocks",
			CountryCode: "US",
		},
	}
	t.Run("ListAssets_Success", func(t *testing.T) {
		mockAssets := pnl.MockAssets{}
		mockAssets.ListAssetsResponse = func() ([]pnl.Asset, error) {
			return testAssets, nil
		}
		querier := NewAssetQuerier(mockAssets)
		req := QueryAssetsReq{}
		want := QueryAssetsResp{Assets: testAssets}
		got, err := querier.ListAssets(req)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if !reflect.DeepEqual(*got, want) {
			t.Errorf("ListAssets() = %v, want %v", *got, want)
		}
	})

	t.Run("ListAssets_Error", func(t *testing.T) {
		// Simulate an error in listing assets.
		mockAssets := pnl.MockAssets{}
		mockAssets.ListAssetsResponse = func() ([]pnl.Asset, error) {
			return nil, errors.New("Mock error")
		}
		querier := NewAssetQuerier(mockAssets)
		req := QueryAssetsReq{}
		_, err := querier.ListAssets(req)
		if err == nil {
			t.Error("Expected an error, got nil")
		}
	})
}

func TestGetAsset(t *testing.T) {

	t.Run("GetAsset_Success", func(t *testing.T) {
		testAsset := pnl.Asset{
			Symbol:      "AAPL",
			Color:       "#FFFFFF",
			Name:        "Apple Stocks",
			CountryCode: "US",
		}
		mockAssets := &pnl.MockAssets{}
		mockAssets.GetAssetResponse = func(symbol string) (*pnl.Asset, error) {
			if symbol != testAsset.Symbol {
				t.Errorf("Expect to get %s, got %s instead", testAsset.Color, symbol)
			}
			return &testAsset, nil
		}
		querier := NewAssetQuerier(mockAssets)
		req := QueryAssetReq{Symbol: testAsset.Symbol}
		want := QueryAssetResp{
			Asset: testAsset,
		}

		got, err := querier.GetAsset(req)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if !reflect.DeepEqual(*got, want) {
			t.Errorf("GetAsset() = %v, want %v", *got, want)
		}
	})

	t.Run("GetAsset_NotFound", func(t *testing.T) {
		mockAssets := &pnl.MockAssets{}
		mockAssets.GetAssetResponse = func(symbol string) (*pnl.Asset, error) {
			return nil, errors.New("Asset not found error")
		}
		querier := NewAssetQuerier(mockAssets)
		req := QueryAssetReq{Symbol: "MSFT"}
		_, err := querier.GetAsset(req)
		if err == nil {
			t.Error("Expected an error for asset not found, got nil")
		}
	})
}

func TestGetTradingPair(t *testing.T) {
	testTradinPair := pnl.TradingPair{
		ID: pnl.TradingPairID("BTCUSD"),
		BaseAsset: pnl.Asset{
			Symbol: "BTC",
			Name:   "Bitcoin",
		},
		QuoteAsset: pnl.Asset{
			Symbol: "USD",
			Name:   "US Dollar",
		},
		Transactions: []pnl.Transaction{
			{
				ID:              "1",
				Timestamp:       time.Now(),
				BaseAmount:      1.0,
				QuoteAmount:     50000.0,
				FeeInBase:      0.1,
				FeeInQuote:   0.0,
				TransactionType: pnl.Buy,
				Price:           50000.0,
			},
			{
				ID:              "2",
				Timestamp:       time.Now(),
				BaseAmount:      0.5,
				QuoteAmount:     25000.0,
				FeeInBase:      0.05,
				FeeInQuote:   0.0,
				TransactionType: pnl.Sell,
				Price:           50000.0,
			},
		},
		Calculations: pnl.Calculations{
			AvgBuyPrice:              50000.0,
			BaseMarketPrice:          48000.0,
			MarketName:               "Mock",
			MarketColor:              "#FFFFFF",
			CurrentBaseAmountInQuote: 48000.0,
			TotalBase:                1.5,
			TotalQuoteSpent:          75000.0,
			PNLAmount:                2500.0,
			PNLPercent:               5.0,
			TotalFeeInQuote:     0.15,
			TotalFeeInBase:  0.0,
		},
	}
	t.Run("GetTradingPair_Success", func(t *testing.T) {
		mockTradingPairs := &pnl.MockTradingPairs{}
		mockTradingPairs.GetTradingPairResponse = func(tpid string) (*pnl.TradingPair, error) {
			return &testTradinPair, nil
		}
		mockMarkets := &pnl.MockMarkets{}
		querier := NewTradingPairQuerier(mockTradingPairs, []pnl.Markets{mockMarkets})
		req := GetTradingPairReq{
			TPID: string(testTradinPair.ID),
		}
		want := GetTradingPairResp{
			Pair: testTradinPair,
		}
		got, err := querier.GetTradingPair(req)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if !reflect.DeepEqual(*got, want) {
			t.Errorf("GetTradingPair() = %v, want %v", *got, want)
		}
	})

	t.Run("GetTradingPairWithBasePrice_Success", func(t *testing.T) {
		mockTradingPairs := &pnl.MockTradingPairs{}
		mockTradingPairs.GetTradingPairResponse = func(tpid string) (*pnl.TradingPair, error) {
			return &testTradinPair, nil
		}
		mockMarkets := &pnl.MockMarkets{}
		mockMarkets.GetCurrentPriceResponse = func(assetA, assetB string) (float64, error) {
			return testTradinPair.Calculations.BaseMarketPrice, nil
		}
		querier := NewTradingPairQuerier(mockTradingPairs, []pnl.Markets{mockMarkets})
		req := GetTradingPairReq{
			TPID:                 string(testTradinPair.ID),
			WithCurrentBasePrice: true,
		}
		want := GetTradingPairResp{
			Pair: testTradinPair,
		}
		got, err := querier.GetTradingPair(req)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if !reflect.DeepEqual(*got, want) {
			t.Errorf("GetTradingPair() = %v, want %v", *got, want)
		}
	})

}

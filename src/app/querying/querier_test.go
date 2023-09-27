package querying

import (
	"controtto/src/domain/pnl"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestListAssets(t *testing.T) {
	a1, _ := pnl.NewAsset("AAPL", "#FFFFFF", "Apple Stocks", "US", "Stock")
	a2, _ := pnl.NewAsset("GOOGL", "#FF5370", "Google Stocks", "US", "Stock")
	a3, _ := pnl.NewAsset("TSLA", "#0F111A", "Tesla Stocks", "US", "Stock")
	testAssets := []pnl.Asset{*a1, *a2, *a3}
	t.Run("ListAssets_Success", func(t *testing.T) {
		mockAssets := pnl.MockAssets{}
		mockAssets.ListAssetsResponse = func() ([]pnl.Asset, error) {
			return testAssets, nil
		}
		querier := NewAssetQuerier(mockAssets)
		req := ListAssetsReq{}
		want := ListAssetsResp{Assets: testAssets}
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
		req := ListAssetsReq{}
		_, err := querier.ListAssets(req)
		if err == nil {
			t.Error("Expected an error, got nil")
		}
	})
}

func TestGetAsset(t *testing.T) {

	t.Run("GetAsset_Success", func(t *testing.T) {
		testAsset, err := pnl.NewAsset("AAPL", "#FFFFFF", "Apple", "US", "Forex")
		if err != nil {
			t.Errorf("Invalid test, asset should be valid for testing.")
		}
		mockAssets := &pnl.MockAssets{}
		mockAssets.GetAssetResponse = func(symbol string) (*pnl.Asset, error) {
			if symbol != testAsset.Symbol {
				t.Errorf("Expect to get %s, got %s instead", testAsset.Color, symbol)
			}
			return testAsset, nil
		}
		querier := NewAssetQuerier(mockAssets)
		req := QueryAssetReq{Symbol: testAsset.Symbol}
		want := QueryAssetResp{
			Asset: *testAsset,
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
	// TODO: Instanciate this assets from the domain, not like this.
	// Create a new trading pair
	mockMarkets := &pnl.MockMarkets{}
	mockMarkets.GetCurrentPriceResponse = func(a, b string) (float64, error) {
		return 65, nil
	}
	baseAsset, _ := pnl.NewAsset("BTC", "#FFFFFF", "Bitcoin", "", "Crypto")
	quoteAsset, _ := pnl.NewAsset("USD", "#FFFFFF", "US Dollar", "US", "Forex")

	tradingPair, _ := pnl.NewTradingPair(*baseAsset, *quoteAsset)
	tradingPair.Calculations.MarketName = mockMarkets.Name()
	tradingPair.Calculations.MarketColor = mockMarkets.Color()
	// Create transactions
	t1, _ := tradingPair.NewTransaction(1.0, 50000.0, 0.1, 0.0, time.Now(), "Buy")
	t2, _ := tradingPair.NewTransaction(0.5, 25000.0, 0.05, 0.0, time.Now(), "Sell")
	tp := *tradingPair // Pair without the transaction
	tradingPair.Transactions = []pnl.Transaction{*t1, *t2}
	// tradingPair.Calculate()
	t.Run("GetTradingPair_Success", func(t *testing.T) {
		mockTradingPairs := &pnl.MockTradingPairs{}
		mockTradingPairs.GetTradingPairResponse = func(tpid string) (*pnl.TradingPair, error) {
			return &tp, nil
		}
		mockTradingPairs.ListTransactionsResponse = func(tpid string) ([]pnl.Transaction, error) {
			return []pnl.Transaction{*t1, *t2}, nil
		}
		querier := NewTradingPairQuerier(mockTradingPairs, []pnl.Markets{mockMarkets})
		req := GetTradingPairReq{
			TPID:                 string(tradingPair.ID),
			WithCurrentBasePrice: false,
			WithTransactions:     true,
			WithCalculations:     false,
		}
		want := GetTradingPairResp{
			Pair: *tradingPair,
		}
		got, err := querier.GetTradingPair(req)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if !reflect.DeepEqual(*got, want) {
			t.Errorf("GetTradingPair() = %v, want %v", *got, want)
		}
	})

	t.Run("GetTradingPairWithoutTransactions_Success", func(t *testing.T) {
		mockTradingPairs := &pnl.MockTradingPairs{}
		mockTradingPairs.GetTradingPairResponse = func(tpid string) (*pnl.TradingPair, error) {
			return &tp, nil
		}
		// mockTradingPairs.ListTransactionsResponse = func(tpid string) ([]pnl.Transaction, error) {
		// 	return []pnl.Transaction{*t1, *t2}, nil
		// }
		querier := NewTradingPairQuerier(mockTradingPairs, []pnl.Markets{mockMarkets})
		req := GetTradingPairReq{
			TPID:                 string(tradingPair.ID),
			WithCurrentBasePrice: false,
			WithTransactions:     false,
			WithCalculations:     false,
		}
		fmt.Println(tradingPair.Transactions)
		tradingPair.Transactions = []pnl.Transaction{}
		fmt.Println(tradingPair.Transactions)
		want := GetTradingPairResp{
			Pair: *tradingPair,
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
			return tradingPair, nil
		}
		mockMarkets := &pnl.MockMarkets{}
		mockMarkets.GetCurrentPriceResponse = func(assetA, assetB string) (float64, error) {
			return tradingPair.Calculations.BaseMarketPrice, nil
		}
		querier := NewTradingPairQuerier(mockTradingPairs, []pnl.Markets{mockMarkets})
		req := GetTradingPairReq{
			TPID:                 string(tradingPair.ID),
			WithCurrentBasePrice: true,
		}
		want := GetTradingPairResp{
			Pair: *tradingPair,
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

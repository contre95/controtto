package pnl

type MockMarkets struct {
	GetCurrentPriceResponse func(assetA, assetB string) (float64, error)
}

func (mm *MockMarkets) Name() string  { return "Mock" }
func (mm *MockMarkets) Color() string { return "#FFFFFF" }
func (mm *MockMarkets) GetCurrentPrice(assetA, assetB string) (float64, error) {
	if mm.GetCurrentPriceResponse != nil {
		return mm.GetCurrentPriceResponse(assetA, assetB)
	}
	panic("GetCurrentPriceResponse not implemented in mock.")

}

// MockTradingPairs is a mock implementation of the TradingPairs interface for testing purposes.
type MockTradingPairs struct {
	ListTradingPairsResponse  func() ([]TradingPair, error)
	GetTradingPairResponse    func(tpid string) (*TradingPair, error)
	DeleteTradingPairResponse func(tpid string) error
	AddTradingPairResponse    func(tp TradingPair) error
	ListTradesResponse  func(tpid string) ([]Trade, error)
	RecordTradeResponse func(t Trade, tpid TradingPairID) error
	DeleteTradeResponse func(tid string) error
}

func (m *MockTradingPairs) AddTradingPair(tp TradingPair) error {
	if m.AddTradingPairResponse != nil {
		return m.AddTradingPairResponse(tp)
	}
	panic("AddTradingPairResponse not implemented in mock.")
}

// ListTradingPairs is a method of the TradingPairs interface that lists all trading pairs from the repository.
func (m *MockTradingPairs) ListTradingPairs() ([]TradingPair, error) {
	if m.ListTradingPairsResponse != nil {
		return m.ListTradingPairsResponse()
	}
	panic("ListTradingPairsResponse not implemented in mock.")
}

// GetTradingPair is a method of the TradingPairs interface that retrieves a trading pair from the repository by its ID.
func (m *MockTradingPairs) GetTradingPair(tpid string) (*TradingPair, error) {
	if m.GetTradingPairResponse != nil {
		return m.GetTradingPairResponse(tpid)
	}
	panic("GetTradingPairResponse not implemented in mock.")
}

// DeleteTradingPair is a method of the TradingPairs interface that deletes a trading pair from the repository by its ID.
func (m *MockTradingPairs) DeleteTradingPair(tpid string) error {
	if m.DeleteTradingPairResponse != nil {
		return m.DeleteTradingPairResponse(tpid)
	}
	panic("DeleteTradingPairResponse not implemented in mock.")
}

// RecordTrade is a method of the TradingPairs interface that records a trade for a trading pair.
func (m *MockTradingPairs) RecordTrade(t Trade, tpid TradingPairID) error {
	if m.RecordTradeResponse != nil {
		return m.RecordTradeResponse(t, tpid)
	}
	panic("RecordTradeResponse not implemented in mock.")

}

// ListTrades is a method of the TradingPairs interface that lists all trades for a trading pair by its ID.
func (m *MockTradingPairs) ListTrades(tpid string) ([]Trade, error) {
	if m.ListTradesResponse != nil {
		return m.ListTradesResponse(tpid)
	}
	panic("ListTradesResponse not implemented in mock.")
}

// DeleteTrade is a method of the TradingPairs interface that deletes a trade from a trading pair by its ID.
func (m *MockTradingPairs) DeleteTrade(tid string) error {
	if m.DeleteTradeResponse != nil {
		return m.DeleteTradeResponse(tid)
	}
	panic("DeleteTradeResponse not implemented in mock.")
}

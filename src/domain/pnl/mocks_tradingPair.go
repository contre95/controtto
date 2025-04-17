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

// MockPairs is a mock implementation of the Pairs interface for testing purposes.
type MockPairs struct {
	ListPairsResponse  func() ([]Pair, error)
	GetPairResponse    func(tpid string) (*Pair, error)
	DeletePairResponse func(tpid string) error
	AddPairResponse    func(tp Pair) error
	ListTradesResponse  func(tpid string) ([]Trade, error)
	RecordTradeResponse func(t Trade, tpid PairID) error
	DeleteTradeResponse func(tid string) error
}

func (m *MockPairs) AddPair(tp Pair) error {
	if m.AddPairResponse != nil {
		return m.AddPairResponse(tp)
	}
	panic("AddPairResponse not implemented in mock.")
}

// ListPairs is a method of the Pairs interface that lists all trading pairs from the repository.
func (m *MockPairs) ListPairs() ([]Pair, error) {
	if m.ListPairsResponse != nil {
		return m.ListPairsResponse()
	}
	panic("ListPairsResponse not implemented in mock.")
}

// GetPair is a method of the Pairs interface that retrieves a trading pair from the repository by its ID.
func (m *MockPairs) GetPair(tpid string) (*Pair, error) {
	if m.GetPairResponse != nil {
		return m.GetPairResponse(tpid)
	}
	panic("GetPairResponse not implemented in mock.")
}

// DeletePair is a method of the Pairs interface that deletes a trading pair from the repository by its ID.
func (m *MockPairs) DeletePair(tpid string) error {
	if m.DeletePairResponse != nil {
		return m.DeletePairResponse(tpid)
	}
	panic("DeletePairResponse not implemented in mock.")
}

// RecordTrade is a method of the Pairs interface that records a trade for a trading pair.
func (m *MockPairs) RecordTrade(t Trade, tpid PairID) error {
	if m.RecordTradeResponse != nil {
		return m.RecordTradeResponse(t, tpid)
	}
	panic("RecordTradeResponse not implemented in mock.")

}

// ListTrades is a method of the Pairs interface that lists all trades for a trading pair by its ID.
func (m *MockPairs) ListTrades(tpid string) ([]Trade, error) {
	if m.ListTradesResponse != nil {
		return m.ListTradesResponse(tpid)
	}
	panic("ListTradesResponse not implemented in mock.")
}

// DeleteTrade is a method of the Pairs interface that deletes a trade from a trading pair by its ID.
func (m *MockPairs) DeleteTrade(tid string) error {
	if m.DeleteTradeResponse != nil {
		return m.DeleteTradeResponse(tid)
	}
	panic("DeleteTradeResponse not implemented in mock.")
}

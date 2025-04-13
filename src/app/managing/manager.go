package managing

// Service just hols all the managing use cases
type Service struct {
	AssetCreator       AssetsCreator
	TradingPairManager TradingPairsManager
	MarketTradeManager MarketTradeManager
}

// NewService is the interctor for all Managing Use cases
func NewService(ac AssetsCreator, tpc TradingPairsManager, mtm MarketTradeManager) Service {
	return Service{ac, tpc, mtm}
}

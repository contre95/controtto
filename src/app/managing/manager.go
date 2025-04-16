package managing

// Service just hols all the managing use cases
type Service struct {
	AssetManager       AssetsManager
	TradingPairManager TradingPairsManager
	MarketManager      *MarketManager
}

// NewService is the interctor for all Managing Use cases
func NewService(ac AssetsManager, tpc TradingPairsManager, mtm *MarketManager) Service {
	return Service{ac, tpc, mtm}
}

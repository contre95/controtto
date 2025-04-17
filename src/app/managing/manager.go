package managing

import "errors"

var (
	ErrEmptyToken         = errors.New("market traders not configured")
	ErrInvalidTrade       = errors.New("invalid trade parameters")
	ErrMarketNotFound     = errors.New("market traders not found")
	ErrInvalidAssetPair   = errors.New("invalid asset pair")
	ErrMarketNotHealthy   = errors.New("market API not healthy")
	ErrProviderNotFound   = errors.New("price provider not found")
	ErrProviderNotHealthy = errors.New("price provider API not healthy")
)

// Service just hols all the managing use cases
type Service struct {
	AssetManager         AssetsManager
	PairManager   PairsManager
	MarketManager        *MarketManager
	PriceProviderManager *PriceProviderManager
}

// NewService is the interctor for all Managing Use cases
func NewService(ac AssetsManager, tpc PairsManager, mtm *MarketManager, ppm *PriceProviderManager) Service {
	return Service{ac, tpc, mtm, ppm}
}

package priceProviders

type PriceProviderNotFound error

// PriceProvider repository interface
type PriceProviders interface {
	// GetCurrentPrice returns the given price of assetA expressed in terms of assetB, if the value is Provider is not found it returns a MarketNotFound error
	GetCurrentPrice(assetA, assetB string) (float64, error)
	Color() string
	Name() string
}

type PriceProviderResp []struct {
	Price float64 `json:"last"`
}

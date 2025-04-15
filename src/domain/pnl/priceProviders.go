package pnl

type PriceProvider struct {
	Public       bool
	IsSet        bool
	Color        string
	ProviderLogo string
	ProviderName string
	ProviderKey  string
	NeedsToken   bool
	Token        string
	ProviderURL  string
	Env          string
	Init         func(string) PriceAPI
	PriceAPI
}

type PriceProviders map[string]PriceProvider

type PriceProviderNotFound error

// PriceProvider repository interface
type PriceAPI interface {
	// GetCurrentPrice returns the given price of assetA expressed in terms of assetB, if the value is Provider is not found it returns a MarketNotFound error
	GetCurrentPrice(assetA, assetB string) (float64, error)
}

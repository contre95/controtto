package pnl

type PriceProvider struct {
	Public            bool
	TokenSet          bool
	Color             string
	ProviderName      string
	ProviderInputName string
	Token             string
	ProviderURL       string
	Env               string
	API
}

type PriceProviders map[string]PriceProvider

type PriceProviderNotFound error

// PriceProvider repository interface
type API interface {
	// GetCurrentPrice returns the given price of assetA expressed in terms of assetB, if the value is Provider is not found it returns a MarketNotFound error
	GetCurrentPrice(assetA, assetB string) (float64, error)
}

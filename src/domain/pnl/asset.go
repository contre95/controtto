package pnl

// type Symbol string

// const (
// 	BTC  Symbol = "BTC"
// 	AMZN Symbol = "AMZN"
// 	ETH  Symbol = "ETH"
// 	EUR  Symbol = "EUR"
// 	AAPL Symbol = "AAPL"
// 	USD  Symbol = "USD"
// 	UDST Symbol = "UDST"
// )

const (
	Crypto AssetType = "Crypto"
	Forex  AssetType = "Forex"
	Stock  AssetType = "Stock"
)

type AssetType string

func GetValidTypes() []AssetType {
	return []AssetType{Crypto, Forex, Stock}
}

// Asset represents individual assets like BTC, USD, EUR, etc. the Symbol property uniquely identifies an asset.
type Asset struct {
	Symbol      string
	Color       string
	Name        string
	CountryCode string
	Type        AssetType
}

// Assets is the repository that handles the CRUD of Assets
type Assets interface {
	AddAsset(a Asset) error
	ListAssets() ([]Asset, error)
	GetAsset(symbol string) (*Asset, error)
}

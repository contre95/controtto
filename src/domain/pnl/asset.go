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

type AssetType int

// Define constants for AssetType using iota
const (
	Crypto AssetType = iota // 0
	Forex                   // 1
	Stock                   // 2
)

// String returns the string representation of the AssetType
func (a AssetType) String() string {
	return [...]string{"Crypto", "Forex", "Stock"}[a]
}

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

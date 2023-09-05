package markets

// BinanceAPI implements the Markets interface using CoinGecko APIs.
type BinanceAPI struct {
	BaseURL string
}

// NewBinanceAPI creates a new instance of CoinGeckoAPI.
func NewBinanceAPI() *BinanceAPI {
	return &BinanceAPI{
		// BaseURL: "https://api.binance.com/api/v1/ticker/price?symbol=BTCUSDT",
		BaseURL: "https://api.binance.com/api/v1/ticker/price",
	}
}

// GetCurrentPrice retrieves the current price of a cryptocurrency using CoinGecko API.
// Receives a symbol that should be in the form of BTCUSDT or USDTBTC.
func (c *BinanceAPI) GetCurrentPrice(assetA, assetB string) (float64, error) {
	// TODO: Implement me
	return 25746.3400, nil
}

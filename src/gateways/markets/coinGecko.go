package markets

import (
	"github.com/JulianToledano/goingecko"
)

// CoinGeckoAPI implements the Markets interface using CoinGecko APIs.
type CoinGeckoAPI struct {
	BaseURL string
}

// NewCoinGeckoAPI creates a new instance of CoinGeckoAPI.
func NewCoinGeckoAPI() *CoinGeckoAPI {
	return &CoinGeckoAPI{
		BaseURL: "https://api.coingecko.com/api/v3/coins/",
	}
}

// GetCurrentPrice retrieves the current price of a cryptocurrency using CoinGecko API.
func (c *CoinGeckoAPI) GetCurrentPrice(name string) (float64, error) {
	cgClient := goingecko.NewClient(nil)
	defer cgClient.Close()
	return 999, nil
}

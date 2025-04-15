package priceProviders

import (
	"controtto/src/domain/pnl"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
)

// CoinbaseAPI struct implements the Markets interface
type CoinbaseAPI struct {
	BaseURL string
}

// CoinbaseResponse is a struct to represent the response from the Coinbase API
type CoinbaseResponse struct {
	Data struct {
		Amount string `json:"amount"`
	} `json:"data"`
}

// NewCoinbaseAPI creates a new instance of CoinbaseAPI
func NewCoinbaseAPI(token string) pnl.PriceAPI {
	return &CoinbaseAPI{
		BaseURL: "https://api.coinbase.com/v2",
	}
}

func (api *CoinbaseAPI) GetCurrentPrice(assetA, assetB string) (float64, error) {
	abPrice := 1.0
	var err error
	if !slices.Contains([]string{"USDT", "USD"}, assetB) {
		abPrice, err = api.GetCurrentPrice(assetB, "USD")
		if err != nil {
			return 0, err
		}
	}
	// Build the URL for the Coinbase API request
	url := fmt.Sprintf("%s/prices/%s-%s/spot", api.BaseURL, assetA, assetB)
	// Make an HTTP GET request to the Coinbase API
	resp, err := http.Get(url)
	if err != nil {
		return 0.0, err
	}
	defer resp.Body.Close()

	// Check the HTTP status code
	if resp.StatusCode != http.StatusOK {
		return 0.0, errors.New("Failed to fetch data from Coinbase API")
	}

	// Parse the JSON response
	var coinbaseResp CoinbaseResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&coinbaseResp)
	if err != nil {
		return 0.0, err
	}

	// Convert the amount to float64
	price, err := strconv.ParseFloat(coinbaseResp.Data.Amount, 64)
	if err != nil {
		return 0.0, err
	}

	return price / abPrice, nil
}

func (api *CoinbaseAPI) Color() string { return "#0052FF" }

func (api *CoinbaseAPI) Name() string { return "Coinbase" }

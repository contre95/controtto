package markets

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// Define a struct to represent the Binance API response.
type BinanceResponse struct {
	Price string `json:"price"`
}

type BinanceAPI struct {
	BaseURL string
}

// NewBinanceAPI creates a new instance of BinanceAPI.
func NewBinanceAPI() *BinanceAPI {
	return &BinanceAPI{
		BaseURL: "https://api.binance.com/api/v3/ticker/price",
	}
}

// GetCurrentPrice retrieves the current price of a cryptocurrency pair using the Binance API.
func (c *BinanceAPI) GetCurrentPrice(assetA, assetB string) (float64, error) {
	return 555, nil
	// Construct the URL with the cryptocurrency pair symbol.
	url := fmt.Sprintf("%s?symbol=%s%s", c.BaseURL, assetA, assetB)

	// Send a GET request to the Binance API.
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Check the response status code.
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	// Parse the JSON response.
	var binanceResponse BinanceResponse
	err = json.NewDecoder(resp.Body).Decode(&binanceResponse)
	if err != nil {
		return 0, err
	}
	// Convert the price to a float64.
	price, err := stringToFloat64(binanceResponse.Price)
	if err != nil {
		return 0, err
	}
	return price, nil
}

// Helper function to convert a string to a float64.
func stringToFloat64(s string) (float64, error) {
	price, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return price, nil
}

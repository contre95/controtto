package markets

import (
	"controtto/src/domain/pnl"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
)

// Markets repository interface
type Markets interface {
	GetCurrentPrice(assetA, assetB string) (float64, error)
	Color() string
	Name() string
}

// TiingoAPI struct implements the Markets interface
type TiingoAPI struct {
	BaseURL string
	token   string
}

// TiingoResponse is a struct to represent the response from the Tiingo API
type TiingoResponse []struct {
	Price float64 `json:"last"`
}

// NewTiingoAPI creates a new instance of AVantageAPI
func NewTiingoAPI(token string) *TiingoAPI {
	return &TiingoAPI{
		BaseURL: "https://api.tiingo.com",
		token:   token,
	}
}

func (api *TiingoAPI) GetCurrentPrice(assetA, assetB string) (float64, error) {
	abPrice := 1.0
	var err error
	if !slices.Contains([]string{"USDT", "USD"}, assetB) {
		abPrice, err = api.GetCurrentPrice(assetB, "USD")
		if err != nil {
			return 0, err
		}
	}
	// Build the URL for the Tiingo API request
	url := fmt.Sprintf("%s/iex?tickers=%s&token=%s", api.BaseURL, assetA, api.token)
	// Make an HTTP GET request to the Tiingo API
	resp, err := http.Get(url)
	if err != nil {
		return 0.0, err
	}
	defer resp.Body.Close()

	// Check the HTTP status code
	if resp.StatusCode != http.StatusOK {
		return 0.0, errors.New("Failed to fetch data from Tiingo API")
	}

	// Parse the JSON response
	var tiingoResp TiingoResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&tiingoResp)
	if err != nil {
		return 0.0, err
	}
	// Check if the market data is found
	if len(tiingoResp) == 0 || tiingoResp[0].Price == 0.0 {
		return 0.0, pnl.MarketNotFound(errors.New(assetA + " market not found"))
	}

	return tiingoResp[0].Price / abPrice, nil
}

func (api *TiingoAPI) Color() string { return "#AA74EF" }

func (api *TiingoAPI) Name() string { return "Tiingo" }

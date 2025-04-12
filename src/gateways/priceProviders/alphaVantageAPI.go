package priceProviders

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
)

// Define a struct to represent the Bingx API response.
type AVantageResponse struct {
	GlobalQuote struct {
		Symbol           string `json:"01. symbol"`
		Open             string `json:"02. open"`
		High             string `json:"03. high"`
		Low              string `json:"04. low"`
		Price            string `json:"05. price"`
		Volume           string `json:"06. volume"`
		LatestTradingDay string `json:"07. latest trading day"`
		PreviousClose    string `json:"08. previous close"`
		Change           string `json:"09. change"`
		ChangePercent    string `json:"10. change percent"`
	} `json:"Global Quote"`
}

type AVantageAPI struct {
	BaseURL string
	token   string
}

// NewAVantageAPI creates a new instance of AVantageAPI
func NewAVantageAPI(token string) *AVantageAPI {
	return &AVantageAPI{
		BaseURL: "https://www.alphavantage.co/query",
		token:   token,
	}
}

// GetCurrentPrice retrieves the current price of a cryptocurrency pair using the Bingx API.
func (api *AVantageAPI) GetCurrentPrice(assetA, assetB string) (float64, error) {
	abPrice := 1.0
	var err error
	if !slices.Contains([]string{"USDT", "USD"}, assetB) {
		abPrice, err = api.GetCurrentPrice(assetB, "USD")
		if err != nil {
			return 0, err
		}
	}
	// Construct the URL with the cryptocurrency pair symbol.
	url := fmt.Sprintf("%s?function=GLOBAL_QUOTE&symbol=%s&apikey=%s", api.BaseURL, assetA, api.token)

	// Send a GET request to the Bingx API.
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
	AVantageResponse := AVantageResponse{}
	err = json.NewDecoder(resp.Body).Decode(&AVantageResponse)
	if err != nil {
		return 0, err
	}
	// Convert the price to a float64.
	price, err := stringToFloat64(AVantageResponse.GlobalQuote.Price)
	if err != nil {
		return 0, err
	}
	slog.Info("Amazon price", "price", price, "abPrice", abPrice)
	return price / abPrice, nil
}

func (api *AVantageAPI) Name() string  { return "Alpha Vantage" }
func (api *AVantageAPI) Color() string { return "#5CC6B1 " }

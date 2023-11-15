package markets

import (
	"controtto/src/domain/pnl"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Define a struct to represent the Bingx API response.
type BingxResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
		Time   int64  `json:"time"`
	} `json:"data"`
}

type BingxAPI struct {
	BaseURL string
}

// NewBingxAPI creates a new instance of BingxAPI
func NewBingxAPI() *BingxAPI {
	return &BingxAPI{
		BaseURL: "https://open-api.bingx.com/openApi/swap/v2/quote/price",
	}
}

// GetCurrentPrice retrieves the current price of a cryptocurrency pair using the Bingx API.
func (api *BingxAPI) GetCurrentPrice(assetA, assetB string) (float64, error) {
	// Construct the URL with the cryptocurrency pair symbol.
	if assetB != "USDT" {
		bPriceUSDT, err := api.GetCurrentPrice(assetB, "USDT")
		if err != nil {
			return 0, err
		}
		aPriceUSDT, err := api.GetCurrentPrice(assetA, "USDT")
		if err != nil {
			return 0, err
		}
		return aPriceUSDT / bPriceUSDT, nil
	}

	url := fmt.Sprintf("%s?symbol=%s-%s", api.BaseURL, assetA, assetB)

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
	var BingxResponse BingxResponse
	err = json.NewDecoder(resp.Body).Decode(&BingxResponse)
	if err != nil {
		return 0, err
	}
	if BingxResponse.Code != 0 {
		return 0, pnl.MarketNotFound(errors.New("Could not find: " + assetA))
	}
	// Convert the price to a float64.
	price, err := stringToFloat64(BingxResponse.Data.Price)
	if err != nil {
		return 0, err
	}
	return price, nil
}

func (api *BingxAPI) Name() string  { return "BingX" }
func (api *BingxAPI) Color() string { return "#2951F4" }

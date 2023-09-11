package markets

import (
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
func (c *BingxAPI) GetCurrentPrice(assetA, assetB string) (float64, error) {
	// Construct the URL with the cryptocurrency pair symbol.
	url := fmt.Sprintf("%s?symbol=%s-%s", c.BaseURL, assetA, assetB)

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
	if err != nil || BingxResponse.Code != 0 {
		return 0, errors.New("Failed to get price from " + c.Name())
	}
	// Convert the price to a float64.
	price, err := stringToFloat64(BingxResponse.Data.Price)
	if err != nil {
		return 0, err
	}
	return price, nil
}

func (c *BingxAPI) Name() string  { return "Bingx" }
func (c *BingxAPI) Color() string { return "#2951F4" }

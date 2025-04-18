package markets

import (
	"controtto/src/domain/pnl"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Trading212API is a client for the Trading 212 API
type Trading212API struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewTrading212API creates a new Trading212API client
func NewTrading212API(apiKey string) pnl.MarketAPI {
	baseURL := "https://live.trading212.com"
	// if isDemo {
	// 	baseURL = "https://demo.trading212.com"
	// }
	return &Trading212API{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (t *Trading212API) HealthCheck() bool {
	// return true
	_, err := t.AccountDetails()
	return err == nil
}

// Buy places a buy market order (Trading 212 API doesn't seem to have a direct buy/sell with amount in quote currency)
func (t *Trading212API) Buy(options pnl.TradeOptions) (*pnl.Trade, error) {
	// if options.Amount <= 0 {
	// 	return nil, errors.New("amount must be greater than zero")
	// }
	//
	// marketRequest := MarketRequest{
	// 	Quantity: options.Amount,
	// 	Ticker:   string(options.Pair.BaseAsset.Symbol + "_" + options.Pair.QuoteAsset.Symbol), // Assuming pnl.Pair can be directly used as Ticker
	// }
	//
	// order, err := t.placeMarketOrder(marketRequest)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// // Adapt Trading212 Order to pnl.Trade (some fields might not directly map)
	// return &pnl.Trade{
	// 	ID:          strconv.FormatInt(order.ID, 10),
	// 	Timestamp:   order.CreationTime,
	// 	BaseAmount:  order.FilledQuantity, // Assuming filled quantity is the base amount
	// 	QuoteAmount: order.FilledValue,    // Assuming filled value is the quote amount
	// 	FeeInBase:   0,                    // Fees might not be directly available in this call
	// 	FeeInQuote:  0,
	// 	TradeType:   "Buy",
	// 	Price:       0, // Price is not directly returned in the place market order, might need to fetch fills
	// }, nil
	panic("not implemented")
}

// Sell places a sell market order (Trading 212 API doesn't seem to have a direct buy/sell with amount in quote currency)
func (t *Trading212API) Sell(options pnl.TradeOptions) (*pnl.Trade, error) {
	// if options.Amount <= 0 {
	// 	return nil, errors.New("amount must be greater than zero")
	// }
	//
	// marketRequest := MarketRequest{
	// 	Quantity: options.Amount,
	// 	Ticker:   string(options.Pair.BaseAsset.Symbol + "_" + options.Pair.QuoteAsset.Symbol), // Assuming pnl.Pair can be directly used as Ticker
	// }
	//
	// order, err := t.placeMarketOrder(marketRequest)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// // Adapt Trading212 Order to pnl.Trade (some fields might not directly map)
	// return &pnl.Trade{
	// 	ID:          strconv.FormatInt(order.ID, 10),
	// 	Timestamp:   order.CreationTime,
	// 	BaseAmount:  order.FilledQuantity, // Assuming filled quantity is the base amount
	// 	QuoteAmount: order.FilledValue,    // Assuming filled value is the quote amount
	// 	FeeInBase:   0,                    // Fees might not be directly available in this call
	// 	FeeInQuote:  0,
	// 	TradeType:   "Sell",
	// 	Price:       0, // Price is not directly returned in the place market order, might need to fetch fills
	// }, nil
	panic("not implemented")
}

// ImportTrades fetches historical order data from Trading 212 API
func (t *Trading212API) ImportTrades(tradingPair pnl.Pair, since time.Time) ([]pnl.Trade, error) {
	// endpoint := fmt.Sprintf("%s/api/v0/equity/history/orders", t.baseURL)
	// req, err := http.NewRequest("GET", endpoint, nil)
	// if err != nil {
	// 	return nil, err
	// }
	// req.Header.Set("Authorization", t.apiKey)
	//
	// // Add query parameters for ticker (tradingPair) if provided
	// if tradingPair.ID != "" {
	// 	q := req.URL.Query()
	// 	q.Add("ticker", string(tradingPair.BaseAsset.Symbol+"_"+tradingPair.QuoteAsset.Symbol))
	// 	req.URL.RawQuery = q.Encode()
	// }
	//
	// resp, err := t.httpClient.Do(req)
	// if err != nil {
	// 	return nil, err
	// }
	// defer resp.Body.Close()
	//
	// if resp.StatusCode != http.StatusOK {
	// 	bodyBytes, _ := io.ReadAll(resp.Body)
	// 	return nil, fmt.Errorf("failed to fetch historical orders, status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	// }
	//
	// var paginatedResponse PaginatedResponseHistoricalOrder
	// err = json.NewDecoder(resp.Body).Decode(&paginatedResponse)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// var trades []pnl.Trade
	// for _, order := range paginatedResponse.Items {
	// 	// Basic mapping, more details might be available in other endpoints or require further processing
	// 	tradeType := pnl.Buy // Default to buy, need logic to determine from order details
	// 	if order.OrderedQuantity < 0 || (order.Type == "MARKET" && order.FilledQuantity < 0) {
	// 		tradeType = pnl.Sell
	// 	}
	//
	// 	trades = append(trades, pnl.Trade{
	// 		ID:          strconv.FormatInt(order.FillId, 10), // Using FillId as a potential unique trade ID
	// 		Timestamp:   order.DateExecuted,
	// 		BaseAmount:  abs(order.FilledQuantity),
	// 		QuoteAmount: abs(order.FillCost),
	// 		FeeInBase:   0, // Fees are in a separate Taxes field
	// 		FeeInQuote:  0,
	// 		TradeType:   tradeType,
	// 		Price:       order.FillPrice,
	// 	})
	// }
	//
	// return trades, nil
	panic("not implemented")
}

func abs(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}

// AccountDetails fetches account metadata (specifically looking for an email-like identifier)
func (t *Trading212API) AccountDetails() (string, error) {
	endpoint := fmt.Sprintf("%s/api/v0/equity/account/info", t.baseURL)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", t.apiKey)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to fetch account info, status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var account Account
	err = json.NewDecoder(resp.Body).Decode(&account)
	if err != nil {
		return "", err
	}

	fmt.Println(fmt.Sprintf("Account: %d", account.ID), nil)
	return fmt.Sprintf("Account: %d", account.ID), nil
}

// FetchAssetAmount is not directly supported by the Trading 212 API in a simple way.
// It would likely require fetching the portfolio and filtering by symbol.
func (t *Trading212API) FetchAssetAmount(symbol string) (float64, error) {
	endpoint := fmt.Sprintf("%s/api/v0/equity/portfolio", t.baseURL)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Authorization", t.apiKey)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("failed to fetch portfolio, status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var positions []Position
	err = json.NewDecoder(resp.Body).Decode(&positions)
	if err != nil {
		return 0, err
	}
	for _, pos := range positions {
		if pos.Ticker == symbol {
			return pos.Quantity, nil
		}
	}
	return 0, nil // pnl.Asset not found in portfolio
}

// placeMarketOrder makes the API call to place a market order
func (t *Trading212API) placeMarketOrder(request MarketRequest) (*Order, error) {
	panic("not implemented")
}

// --- Data Structures based on the OpenAPI specification ---

// Account schema
type Account struct {
	CurrencyCode string `json:"currencyCode"`
	ID           int64  `json:"id"`
}

// Cash schema
type Cash struct {
	Blocked  float64 `json:"blocked"`
	Free     float64 `json:"free"`
	Invested float64 `json:"invested"`
	PieCash  float64 `json:"pieCash"`
	PPL      float64 `json:"ppl"`
	Result   float64 `json:"result"`
	Total    float64 `json:"total"`
}

// HistoricalOrder schema
type HistoricalOrder struct {
	DateCreated     time.Time `json:"dateCreated"`
	DateExecuted    time.Time `json:"dateExecuted"`
	DateModified    time.Time `json:"dateModified"`
	Executor        string    `json:"executor"`
	FillCost        float64   `json:"fillCost"`
	FillId          int64     `json:"fillId"`
	FillPrice       float64   `json:"fillPrice"`
	FillResult      float64   `json:"fillResult"`
	FillType        string    `json:"fillType"`
	FilledQuantity  float64   `json:"filledQuantity"`
	FilledValue     float64   `json:"filledValue"`
	ID              int64     `json:"id"`
	LimitPrice      float64   `json:"limitPrice"`
	OrderedQuantity float64   `json:"orderedQuantity"`
	OrderedValue    float64   `json:"orderedValue"`
	ParentOrder     int64     `json:"parentOrder"`
	Status          string    `json:"status"`
	StopPrice       float64   `json:"stopPrice"`
	Taxes           []Tax     `json:"taxes"`
	Ticker          string    `json:"ticker"`
	TimeValidity    string    `json:"timeValidity"`
	Type            string    `json:"type"`
}

// PaginatedResponseHistoricalOrder schema
type PaginatedResponseHistoricalOrder struct {
	Items        []HistoricalOrder `json:"items"`
	NextPagePath string            `json:"nextPagePath"`
}

// Order schema
type Order struct {
	CreationTime   time.Time `json:"creationTime"`
	FilledQuantity float64   `json:"filledQuantity"`
	FilledValue    float64   `json:"filledValue"`
	ID             int64     `json:"id"`
	LimitPrice     float64   `json:"limitPrice"`
	Quantity       float64   `json:"quantity"`
	Status         string    `json:"status"`
	StopPrice      float64   `json:"stopPrice"`
	Strategy       string    `json:"strategy"`
	Ticker         string    `json:"ticker"`
	Type           string    `json:"type"`
	Value          float64   `json:"value"`
}

// MarketRequest schema
type MarketRequest struct {
	Quantity float64 `json:"quantity"`
	Ticker   string  `json:"ticker"`
}

// Position schema
type Position struct {
	AveragePrice    float64   `json:"averagePrice"`
	CurrentPrice    float64   `json:"currentPrice"`
	Frontend        string    `json:"frontend"`
	FxPpl           float64   `json:"fxPpl"`
	InitialFillDate time.Time `json:"initialFillDate"`
	MaxBuy          float64   `json:"maxBuy"`
	MaxSell         float64   `json:"maxSell"`
	PieQuantity     float64   `json:"pieQuantity"`
	Ppl             float64   `json:"ppl"`
	Quantity        float64   `json:"quantity"`
	Ticker          string    `json:"ticker"`
}

// Tax schema
type Tax struct {
	FillId      string    `json:"fillId"`
	Name        string    `json:"name"`
	Quantity    float64   `json:"quantity"`
	TimeCharged time.Time `json:"timeCharged"`
}

// --- Other Schemas (only including those potentially used) ---

// PlaceOrderError schema
type PlaceOrderError struct {
	Clarification string `json:"clarification"`
	Code          string `json:"code"`
}

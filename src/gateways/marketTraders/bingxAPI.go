package marketTraders

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	baseURL    = "https://api.bingx.com"
	apiVersion = "v1"
)

type BingXAPI struct {
	apiKey     string
	secretKey  string
	httpClient *http.Client
}

type TradeOptions struct {
	TradingPairID string
	Price         float64
	Amount        float64
	Type          string // "LIMIT", "MARKET", etc.
	ClientOrderID string // optional
}

type Trade struct {
	ID            string
	TradingPairID string
	Price         float64
	Amount        float64
	Fee           float64
	FeeCurrency   string
	Side          string // "BUY" or "SELL"
	Timestamp     time.Time
}

// New creates a new BingX API client
func New(apiKey, secretKey string) *BingXAPI {
	return &BingXAPI{
		apiKey:     apiKey,
		secretKey:  secretKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// GenerateToken generates an authentication token similar to TradingView's style
func (b *BingXAPI) GenerateToken(expiry time.Duration) string {
	expiryTime := time.Now().Add(expiry).Unix()
	data := fmt.Sprintf("%s:%d", b.apiKey, expiryTime)

	h := hmac.New(sha256.New, []byte(b.secretKey))
	h.Write([]byte(data))
	signature := hex.EncodeToString(h.Sum(nil))

	token := fmt.Sprintf("%s:%d:%s", b.apiKey, expiryTime, signature)
	return token
}

// Buy executes a buy order
func (b *BingXAPI) Buy(options TradeOptions) (*Trade, error) {
	return b.placeOrder(options, "BUY")
}

// Sell executes a sell order
func (b *BingXAPI) Sell(options TradeOptions) (*Trade, error) {
	return b.placeOrder(options, "SELL")
}

// placeOrder handles the common logic for placing orders
func (b *BingXAPI) placeOrder(options TradeOptions, side string) (*Trade, error) {
	endpoint := "/api/v1/trade/order"
	params := url.Values{}
	params.Add("symbol", options.TradingPairID)
	params.Add("side", side)
	params.Add("type", options.Type)
	params.Add("quantity", strconv.FormatFloat(options.Amount, 'f', -1, 64))

	if options.Type == "LIMIT" {
		params.Add("price", strconv.FormatFloat(options.Price, 'f', -1, 64))
	}

	if options.ClientOrderID != "" {
		params.Add("clientOrderID", options.ClientOrderID)
	}

	resp, err := b.sendRequest("POST", endpoint, params, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Code    int    `json:"code"`
		Message string `json:"msg"`
		Data    struct {
			OrderID string `json:"orderId"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Code != 0 {
		return nil, errors.New(result.Message)
	}

	// In a real implementation, you might want to fetch the trade details here
	return &Trade{
		ID:            result.Data.OrderID,
		TradingPairID: options.TradingPairID,
		Price:         options.Price,
		Amount:        options.Amount,
		Side:          side,
		Timestamp:     time.Now(),
	}, nil
}

// ImportTrades fetches historical trades for a trading pair
func (b *BingXAPI) ImportTrades(tradingPairID TradingPairID, since time.Time) ([]Trade, error) {
	endpoint := "/api/v1/trade/history"
	params := url.Values{}
	params.Add("symbol", string(tradingPairID))
	params.Add("startTime", strconv.FormatInt(since.UnixMilli(), 10))

	resp, err := b.sendRequest("GET", endpoint, params, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Code    int    `json:"code"`
		Message string `json:"msg"`
		Data    []struct {
			ID          string  `json:"id"`
			Price       float64 `json:"price,string"`
			Quantity    float64 `json:"quantity,string"`
			Fee         float64 `json:"fee,string"`
			FeeCurrency string  `json:"feeCurrency"`
			Side        string  `json:"side"`
			Time        int64   `json:"time"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Code != 0 {
		return nil, errors.New(result.Message)
	}

	var trades []Trade
	for _, t := range result.Data {
		trades = append(trades, Trade{
			ID:            t.ID,
			TradingPairID: string(tradingPairID),
			Price:         t.Price,
			Amount:        t.Quantity,
			Fee:           t.Fee,
			FeeCurrency:   t.FeeCurrency,
			Side:          t.Side,
			Timestamp:     time.UnixMilli(t.Time),
		})
	}

	return trades, nil
}

// FetchAsset retrieves the balance of a specific asset
func (b *BingXAPI) FetchAsset(symbol string) (float64, error) {
	endpoint := "/api/v1/account/balance"
	params := url.Values{}

	resp, err := b.sendRequest("GET", endpoint, params, true)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Code    int    `json:"code"`
		Message string `json:"msg"`
		Data    []struct {
			Asset  string `json:"asset"`
			Free   string `json:"free"`
			Locked string `json:"locked"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	if result.Code != 0 {
		return 0, errors.New(result.Message)
	}

	for _, asset := range result.Data {
		if asset.Asset == symbol {
			free, err := strconv.ParseFloat(asset.Free, 64)
			if err != nil {
				return 0, err
			}
			locked, err := strconv.ParseFloat(asset.Locked, 64)
			if err != nil {
				return 0, err
			}
			return free + locked, nil
		}
	}

	return 0, fmt.Errorf("asset %s not found", symbol)
}

// sendRequest handles the HTTP request/response cycle
func (b *BingXAPI) sendRequest(method, endpoint string, params url.Values, signed bool) (*http.Response, error) {
	var req *http.Request
	var err error

	urlStr := baseURL + endpoint

	if method == "GET" {
		if len(params) > 0 {
			urlStr += "?" + params.Encode()
		}
		req, err = http.NewRequest(method, urlStr, nil)
	} else {
		req, err = http.NewRequest(method, urlStr, bytes.NewBufferString(params.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	if err != nil {
		return nil, err
	}

	if signed {
		timestamp := time.Now().UnixMilli()
		req.Header.Set("X-BX-APIKEY", b.apiKey)
		req.Header.Set("X-BX-TIMESTAMP", strconv.FormatInt(timestamp, 10))

		// In a real implementation, you would need to properly sign the request
		// This is a simplified version
		signature := b.generateSignature(params, timestamp)
		req.Header.Set("X-BX-SIGNATURE", signature)
	}

	return b.httpClient.Do(req)
}

// generateSignature creates the HMAC-SHA256 signature for requests
func (b *BingXAPI) generateSignature(params url.Values, timestamp int64) string {
	queryString := params.Encode()
	data := fmt.Sprintf("%d%s%s", timestamp, b.apiKey, queryString)

	h := hmac.New(sha256.New, []byte(b.secretKey))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

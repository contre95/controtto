package markets

import (
	"controtto/src/domain/pnl"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type BingxMarketAPI struct {
	ApiKey    string
	ApiSecret string
}

const HOST = "https://open-api.bingx.com"

func NewBingXAPI(token string) pnl.MarketAPI {
	if len(strings.Split(token, ":")) >= 2 {
		return &BingxMarketAPI{
			ApiKey:    strings.Split(token, ":")[0],
			ApiSecret: strings.Split(token, ":")[1],
		}
	}
	return &BingxMarketAPI{}
}

func (b *BingxMarketAPI) HeatlhCheck() bool {
	return true
}

func (b *BingxMarketAPI) getParameters(payload map[string]string, urlEncode bool, timestamp int64) string {
	params := ""
	for k, v := range payload {
		encoded := v
		if urlEncode {
			encoded = url.QueryEscape(v)
			encoded = strings.ReplaceAll(encoded, "+", "%20")
		}
		params += fmt.Sprintf("%s=%s&", k, encoded)
	}
	params += fmt.Sprintf("timestamp=%d", timestamp)
	return params
}

func computeHmac256(strMessage string, strSecret string) string {
	key := []byte(strSecret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(strMessage))
	return hex.EncodeToString(h.Sum(nil))
}

func (b *BingxMarketAPI) FetchAssetAmount(symbol string) (float64, error) {
	uri := "/openApi/spot/v1/account/balance"
	method := "GET"
	timestamp := time.Now().UnixNano() / 1e6
	payload := map[string]string{
		"recvWindow": "60000",
	}
	paramStr := b.getParameters(payload, false, timestamp)
	sign := computeHmac256(paramStr, b.ApiSecret)
	urlParams := b.getParameters(payload, true, timestamp) + "&signature=" + sign
	fullURL := fmt.Sprintf("%s%s?%s", HOST, uri, urlParams)

	req, err := http.NewRequest(method, fullURL, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("X-BX-APIKEY", b.ApiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	var res struct {
		Code     int    `json:"code"`
		Msg      string `json:"msg"`
		DebugMsg string `json:"debugMsg"`
		Data     struct {
			Balances []struct {
				Asset  string `json:"asset"`
				Free   string `json:"free"`
				Locked string `json:"locked"`
			} `json:"balances"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &res); err != nil {
		slog.Error("BingxAPI: Error unmarshalling response", "error", err)
		return 0, err
	}
	if res.Code != 0 {
		slog.Error("BingxAPI: Error in response", "code", res.Code, "msg", res.Msg, "debugMsg", res.DebugMsg)
		return 0, errors.New("bingx error: " + res.Msg)
	}

	for _, asset := range res.Data.Balances {
		if strings.EqualFold(asset.Asset, symbol) {
			var amt float64
			fmt.Sscanf(asset.Free, "%f", &amt)
			return amt, nil
		}
	}
	return 0, errors.New("asset not found")
}

// Unimplemented methods for now
func (b *BingxMarketAPI) Buy(options pnl.TradeOptions) (*pnl.Trade, error) {
	panic("Buy not implemented")
}

func (b *BingxMarketAPI) Sell(options pnl.TradeOptions) (*pnl.Trade, error) {
	panic("Sell not implemented")
}

func (b *BingxMarketAPI) ImportTrades(pair pnl.Pair, since time.Time) ([]pnl.Trade, error) {
	// Validate input parameters
	if pair.BaseAsset.Symbol == "" || pair.QuoteAsset.Symbol == "" {
		return nil, fmt.Errorf("invalid pair: base and quote symbols required")
	}
	uri := "/openApi/spot/v1/trade/historyOrders"
	method := "GET"
	timestamp := time.Now().UnixMilli()
	endTime := time.Now().UnixMilli()

	payload := map[string]string{
		"symbol":    fmt.Sprintf("%s-%s", pair.BaseAsset.Symbol, pair.QuoteAsset.Symbol),
		"startTime": fmt.Sprintf("%d", since.UnixMilli()),
		"endTime":   fmt.Sprintf("%d", endTime),
		"timestamp": fmt.Sprintf("%d", timestamp),
		"limit":     "500", // Use maximum allowed limit
	}

	// Create parameter string for signing
	var paramBuilder strings.Builder
	keys := make([]string, 0, len(payload))
	for k := range payload {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i, k := range keys {
		if i > 0 {
			paramBuilder.WriteString("&")
		}
		paramBuilder.WriteString(k)
		paramBuilder.WriteString("=")
		paramBuilder.WriteString(payload[k])
	}
	paramStr := paramBuilder.String()

	// Compute signature
	sign := computeHmac256(paramStr, b.ApiSecret)

	// Create request
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s?%s&signature=%s", HOST, uri, paramStr, sign), nil)
	if err != nil {
		slog.Error("BingxAPI: Error creating request", "error", err)
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("X-BX-APIKEY", b.ApiKey)
	req.Header.Set("Accept", "application/json")

	// Execute request with timeout
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("BingxAPI: Error executing request", "error", err)
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Handle non-200 responses
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		slog.Error("BingxAPI: Unexpected status code",
			"status", resp.StatusCode,
			"url", req.URL.Redacted(),
			"body", string(body))
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	// Parse response
	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Orders []struct {
				OrderID     int64   `json:"orderId"`
				Symbol      string  `json:"symbol"`
				Price       string  `json:"price"`
				StopPrice   string  `json:"StopPrice"`
				OrigQty     string  `json:"origQty"`
				ExecutedQty string  `json:"executedQty"`
				CumQuoteQty string  `json:"cummulativeQuoteQty"`
				Status      string  `json:"status"`
				Type        string  `json:"type"`
				Side        string  `json:"side"`
				Time        int64   `json:"time"`
				UpdateTime  int64   `json:"updateTime"`
				Fee         float64 `json:"fee"`
				FeeAsset    string  `json:"feeAsset"`
				AvgPrice    float64 `json:"avgPrice"`
			} `json:"orders"`
		} `json:"data"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("BingxAPI: Error reading response body", "error", err)
		return nil, fmt.Errorf("read body failed: %w", err)
	}

	if err := json.Unmarshal(body, &response); err != nil {
		slog.Error("BingxAPI: Error unmarshalling response", "error", err, "body", string(body))
		return nil, fmt.Errorf("parse response failed: %w", err)
	}

	if response.Code != 0 {
		slog.Error("BingxAPI: Error in response", "code", response.Code, "msg", response.Msg)
		return nil, fmt.Errorf("api error: %d - %s", response.Code, response.Msg)
	}

	var trades []pnl.Trade
	for _, order := range response.Data.Orders {
		// Skip orders that aren't filled
		if order.Status != "FILLED" {
			continue
		}

		// Parse timestamps
		tradeTime := time.Unix(0, order.UpdateTime*int64(time.Millisecond))
		if tradeTime.Before(since) {
			continue
		}

		// Parse numeric values
		price, err := strconv.ParseFloat(order.Price, 64)
		if err != nil {
			// Fallback to average price if available
			if order.AvgPrice > 0 {
				price = order.AvgPrice
			} else {
				slog.Warn("BingxAPI: Skipping trade with invalid price", "price", order.Price, "error", err)
				continue
			}
		}

		quantity, err := strconv.ParseFloat(order.ExecutedQty, 64)
		if err != nil {
			slog.Warn("BingxAPI: Skipping trade with invalid quantity", "qty", order.ExecutedQty, "error", err)
			continue
		}

		// Determine trade type
		var tradeType pnl.TradeType
		if order.Side == "BUY" {
			tradeType = pnl.Buy
		} else {
			tradeType = pnl.Sell
		}

		// Handle fees (fee is already float64 in the response)
		feeInBase := 0.0
		feeInQuote := 0.0
		if strings.EqualFold(order.FeeAsset, pair.BaseAsset.Symbol) {
			feeInBase = math.Abs(order.Fee)
		} else if strings.EqualFold(order.FeeAsset, pair.QuoteAsset.Symbol) {
			feeInQuote = math.Abs(order.Fee)
		}

		// Calculate quote amount
		quoteAmount, err := strconv.ParseFloat(order.CumQuoteQty, 64)
		if err != nil {
			quoteAmount = price * quantity // Fallback calculation
		}

		trades = append(trades, pnl.Trade{
			Timestamp:   tradeTime,
			BaseAmount:  quantity,
			QuoteAmount: quoteAmount,
			FeeInBase:   feeInBase,
			FeeInQuote:  feeInQuote,
			TradeType:   tradeType,
			Price:       price,
			ID:          strconv.FormatInt(order.OrderID, 10),
		})
	}

	return trades, nil
}

func (b *BingxMarketAPI) HealthCheck() bool {
	details, err := b.AccountDetails()
	if err != nil {
		return false
	}
	fmt.Println("Account details:", details)
	return true
}

func (b *BingxMarketAPI) AccountDetails() (string, error) {
	uri := "/openApi/v1/account/apiPermissions"
	method := "GET"
	timestamp := time.Now().UnixNano() / 1e6

	payload := map[string]string{}

	paramStr := b.getParameters(payload, false, timestamp)
	sign := computeHmac256(paramStr, b.ApiSecret)
	urlParams := b.getParameters(payload, true, timestamp) + "&signature=" + sign
	fullURL := fmt.Sprintf("%s%s?%s", HOST, uri, urlParams)

	req, err := http.NewRequest(method, fullURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("X-BX-APIKEY", b.ApiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var res struct {
		Permissions []int    `json:"permissions"`
		IPAddresses []string `json:"ipAddresses"`
		Note        string   `json:"note"`
		APIKey      string   `json:"apiKey"`
	}
	fmt.Print(string(body))
	if err := json.Unmarshal(body, &res); err != nil {
		return "", err
	}

	if len(res.Note) == 0 {
		return "", errors.New("no IP addresses found")
	}
	allowedIPs := strings.Join(res.IPAddresses, ", ")
	// fmt.Printf("%q Allowed IPs: %q", res.Data.Note, allowedIPs)
	return fmt.Sprintf("%q Allowed IPs: %q", res.Note, allowedIPs), nil
}

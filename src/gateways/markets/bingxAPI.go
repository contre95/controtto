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
	uri := "/openApi/swap/v2/trade/allFillOrders"
	method := "GET"
	timestamp := time.Now().UnixNano() / 1e6
	payload := map[string]string{
		"symbol":    fmt.Sprintf("%s-%s", pair.BaseAsset.Symbol, pair.QuoteAsset.Symbol),
		"startTs":   fmt.Sprintf("%d", since.UnixMilli()),
		"endTs":     fmt.Sprintf("%d", time.Now().UnixMilli()),
		"timestamp": fmt.Sprintf("%d", timestamp), // Add timestamp to payload
	}

	// Create parameter string for signing (no URL encoding)
	var paramBuilder strings.Builder
	keys := make([]string, 0, len(payload))
	for k := range payload {
		keys = append(keys, k)
	}
	sort.Strings(keys) // Parameters must be sorted alphabetically

	for _, k := range keys {
		paramBuilder.WriteString(k)
		paramBuilder.WriteString("=")
		paramBuilder.WriteString(payload[k])
		paramBuilder.WriteString("&")
	}
	paramStr := strings.TrimSuffix(paramBuilder.String(), "&")

	// Compute signature
	sign := computeHmac256(paramStr, b.ApiSecret)

	// Create URL-encoded parameters for the request
	urlParams := url.Values{}
	for k, v := range payload {
		urlParams.Add(k, v)
	}
	urlParams.Add("signature", sign)

	fullURL := fmt.Sprintf("%s%s?%s", HOST, uri, urlParams.Encode())

	req, err := http.NewRequest(method, fullURL, nil)
	if err != nil {
		slog.Error("BingxAPI: Error creating request", "error", err)
		return nil, fmt.Errorf("create request failed: %v", err)
	}
	req.Header.Set("X-BX-APIKEY", b.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("BingxAPI: Error executing request", "error", err)
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		slog.Error("BingxAPI: Unexpected status code", "status", resp.StatusCode, "body", string(body))
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("BingxAPI: Error reading response body", "error", err)
		return nil, fmt.Errorf("read body failed: %v", err)
	}
	fmt.Println(string(body))
	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			FillOrders []struct {
				FilledTm   string `json:"filledTm"`
				Volume     string `json:"volume"`
				Price      string `json:"price"`
				Amount     string `json:"amount"`
				Commission string `json:"commission"`
				Currency   string `json:"currency"`
				OrderID    string `json:"orderId"`
				FilledTime string `json:"filledTime"`
				Symbol     string `json:"symbol"`
			} `json:"fill_orders"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		slog.Error("BingxAPI: Error unmarshalling response", "error", err, "body", string(body))
		return nil, fmt.Errorf("parse response failed: %v", err)
	}

	if response.Code != 0 {
		slog.Error("BingxAPI: Error in response", "code", response.Code, "msg", response.Msg)
		return nil, fmt.Errorf("api error: %d - %s", response.Code, response.Msg)
	}

	var trades []pnl.Trade
	for _, tradeData := range response.Data.FillOrders {
		tradeTime, err := time.Parse(time.RFC3339, tradeData.FilledTm)
		if err != nil {
			slog.Warn("BingxAPI: Skipping trade with invalid timestamp", "time", tradeData.FilledTm, "error", err)
			continue
		}

		if tradeTime.Before(since) {
			continue
		}

		// Parse string values to float64
		price, err := strconv.ParseFloat(tradeData.Price, 64)
		if err != nil {
			slog.Warn("BingxAPI: Skipping trade with invalid price", "price", tradeData.Price, "error", err)
			continue
		}

		quantity, err := strconv.ParseFloat(tradeData.Volume, 64)
		if err != nil {
			slog.Warn("BingxAPI: Skipping trade with invalid quantity", "volume", tradeData.Volume, "error", err)
			continue
		}

		commission, err := strconv.ParseFloat(tradeData.Commission, 64)
		if err != nil {
			slog.Warn("BingxAPI: Skipping trade with invalid commission", "commission", tradeData.Commission, "error", err)
			continue
		}

		quoteAmount, err := strconv.ParseFloat(tradeData.Amount, 64)
		if err != nil {
			quoteAmount = price * quantity // Fallback calculation
		}

		// Determine trade type based on commission sign
		// Negative commission means it was deducted (usually for taker trades)
		tradeType := pnl.Sell
		if strings.Contains(tradeData.Commission, "-") {
			tradeType = pnl.Buy
		}

		// Handle fees based on currency
		feeInBase := 0.0
		feeInQuote := 0.0
		if strings.EqualFold(tradeData.Currency, pair.BaseAsset.Symbol) {
			feeInBase = math.Abs(commission)
		} else if strings.EqualFold(tradeData.Currency, pair.QuoteAsset.Symbol) {
			feeInQuote = math.Abs(commission)
		}

		trades = append(trades, pnl.Trade{
			Timestamp:   tradeTime,
			BaseAmount:  quantity,
			QuoteAmount: quoteAmount,
			FeeInBase:   feeInBase,
			FeeInQuote:  feeInQuote,
			TradeType:   tradeType,
			Price:       price,
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

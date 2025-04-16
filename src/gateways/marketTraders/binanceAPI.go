package marketTraders

import (
	"controtto/src/domain/pnl"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type BinanceMarketAPI struct {
	ApiKey    string
	ApiSecret string
	BaseURL   string
}

func NewBinanceAPI(token string) pnl.MarketAPI {
	api := &BinanceMarketAPI{
		BaseURL: "https://api.binance.com",
	}
	if len(strings.Split(token, ":")) >= 2 {
		api.ApiKey = strings.Split(token, ":")[0]
		api.ApiSecret = strings.Split(token, ":")[1]
	}
	return api
}

func (b *BinanceMarketAPI) HeatlhCheck() bool {
	_, err := b.ping()
	return err == nil
}

func (b *BinanceMarketAPI) ping() (string, error) {
	url := fmt.Sprintf("%s/api/v3/ping", b.BaseURL)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("binance ping returned status: %d", resp.StatusCode)
	}
	return "pong", nil
}

func (b *BinanceMarketAPI) getSignedParams(payload map[string]string) string {
	params := url.Values{}
	for k, v := range payload {
		params.Add(k, v)
	}
	params.Add("timestamp", strconv.FormatInt(time.Now().UnixNano()/1e6, 10))
	params.Add("recvWindow", "5000") // Binance recommends 5000ms

	return params.Encode()
}

func (b *BinanceMarketAPI) signRequest(params string) string {
	mac := hmac.New(sha256.New, []byte(b.ApiSecret))
	mac.Write([]byte(params))
	return hex.EncodeToString(mac.Sum(nil))
}

func (b *BinanceMarketAPI) makeRequest(method, endpoint string, params map[string]string, signed bool) ([]byte, error) {
	var req *http.Request
	var err error

	urlStr := fmt.Sprintf("%s%s", b.BaseURL, endpoint)

	if method == http.MethodGet {
		queryParams := b.getSignedParams(params)
		if signed {
			signature := b.signRequest(queryParams)
			queryParams += "&signature=" + signature
		}
		urlStr += "?" + queryParams
		req, err = http.NewRequest(method, urlStr, nil)
	} else {
		// Handle POST requests if needed
		return nil, errors.New("POST requests not implemented")
	}

	if err != nil {
		return nil, err
	}

	if signed {
		req.Header.Set("X-MBX-APIKEY", b.ApiKey)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("binance API returned status: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func (b *BinanceMarketAPI) FetchAssetAmount(symbol string) (float64, error) {
	if b.ApiKey == "" || b.ApiSecret == "" {
		return 0, errors.New("missing API credentials")
	}

	body, err := b.makeRequest(
		http.MethodGet,
		"/api/v3/account",
		make(map[string]string),
		true,
	)
	if err != nil {
		return 0, err
	}

	var accountInfo struct {
		AccountType string `json:"accountType"`
		Balances    []struct {
			Asset  string `json:"asset"`
			Free   string `json:"free"`
			Locked string `json:"locked"`
		} `json:"balances"`
	}

	if err := json.Unmarshal(body, &accountInfo); err != nil {
		return 0, fmt.Errorf("failed to parse account info: %v", err)
	}
	fmt.Println(accountInfo.AccountType)
	for _, balance := range accountInfo.Balances {
		if strings.EqualFold(balance.Asset, symbol) {
			free, err := strconv.ParseFloat(balance.Free, 64)
			if err != nil {
				return 0, fmt.Errorf("failed to parse free amount: %v", err)
			}
			locked, err := strconv.ParseFloat(balance.Locked, 64)
			if err != nil {
				return 0, fmt.Errorf("failed to parse locked amount: %v", err)
			}
			fmt.Println("Free:", free, "Locked:", locked, "balance:", balance.Asset)
			return free + locked, nil
		}
	}

	return 0, fmt.Errorf("asset %s not found", symbol)
}

func (b *BinanceMarketAPI) HealthCheck() bool {
	// First check basic connectivity
	if _, err := b.ping(); err != nil {
		return false
	}

	// Then check account access if credentials are available
	if b.ApiKey != "" && b.ApiSecret != "" {
		_, err := b.FetchAssetAmount("BTC") // Test with BTC as it's likely to exist
		return err == nil
	}

	return true
}

func (b *BinanceMarketAPI) AccountDetails() (string, error) {
	body, err := b.makeRequest(
		http.MethodGet,
		"/api/v3/account",
		make(map[string]string),
		true,
	)
	if err != nil {
		return "", err
	}

	var accountInfo struct {
		MakerCommission  int      `json:"makerCommission"`
		TakerCommission  int      `json:"takerCommission"`
		BuyerCommission  int      `json:"buyerCommission"`
		SellerCommission int      `json:"sellerCommission"`
		CanTrade         bool     `json:"canTrade"`
		CanWithdraw      bool     `json:"canWithdraw"`
		CanDeposit       bool     `json:"canDeposit"`
		UpdateTime       int64    `json:"updateTime"`
		AccountType      string   `json:"accountType"`
		Permissions      []string `json:"permissions"`
	}

	if err := json.Unmarshal(body, &accountInfo); err != nil {
		return "", fmt.Errorf("failed to parse account info: %v", err)
	}

	return fmt.Sprintf(
		"Account Type: %s, Permissions: %v, Trading: %t, Withdrawals: %t, Deposits: %t",
		accountInfo.AccountType,
		accountInfo.Permissions,
		accountInfo.CanTrade,
		accountInfo.CanWithdraw,
		accountInfo.CanDeposit,
	), nil
}

// Unimplemented methods remain the same
func (b *BinanceMarketAPI) Buy(options pnl.TradeOptions) (*pnl.Trade, error) {
	panic("Buy not implemented")
}

func (b *BinanceMarketAPI) Sell(options pnl.TradeOptions) (*pnl.Trade, error) {
	panic("Sell not implemented")
}

func (b *BinanceMarketAPI) ImportTrades(pair pnl.TradingPair, since time.Time) ([]pnl.Trade, error) {
	panic("ImportTrades not implemented")
}

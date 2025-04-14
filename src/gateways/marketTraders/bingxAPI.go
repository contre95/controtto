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
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type BingxMarketAPI struct {
	ApiKey    string
	ApiSecret string
}

func NewBingXAPI(token string) *BingxMarketAPI {
	if len(strings.Split(token, ":")) >= 2 {
		return &BingxMarketAPI{
			ApiKey:    strings.Split(token, ":")[0],
			ApiSecret: strings.Split(token, ":")[1],
		}
	}
	return &BingxMarketAPI{}
}

func (b *BingxMarketAPI) getParameters(uri, method string, payload map[string]string, urlEncode bool, timestamp int64) string {
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

func (b *BingxMarketAPI) FetchAsset(symbol string) (float64, error) {
	const HOST = "https://open-api.bingx.com"
	uri := "/openApi/spot/v1/account/balance"
	method := "GET"
	timestamp := time.Now().UnixMilli()

	payload := map[string]string{
		"recvWindow": "60000",
	}

	paramStr := b.getParameters(uri, method, payload, false, timestamp)
	sign := computeHmac256(paramStr, b.ApiSecret)
	urlParams := b.getParameters(uri, method, payload, true, timestamp) + "&signature=" + sign
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
	fmt.Println(string(body))
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
	fmt.Println("Not found")
	return 30, nil
}

// Unimplemented methods for now
func (b *BingxMarketAPI) Buy(options pnl.TradeOptions) (*pnl.Trade, error) {
	panic("Buy not implemented")
}

func (b *BingxMarketAPI) Sell(options pnl.TradeOptions) (*pnl.Trade, error) {
	panic("Sell not implemented")
}

func (b *BingxMarketAPI) ImportTrades(pair pnl.TradingPair, since time.Time) ([]pnl.Trade, error) {
	panic("ImportTrades not implemented")
}

func (b *BingxMarketAPI) AccountDetails() (string, error) {
	panic("AccountDetails not implemented")
}

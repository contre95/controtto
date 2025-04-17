package marketTraders

import (
	"controtto/src/domain/pnl"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type BitcoinWalletAPI struct {
	Address string
	BaseURL string // e.g., https://blockstream.info/api
}

func NewBitcoinWalletAPI(token string) pnl.MarketAPI {
	parts := strings.Split(token, ":")
	address := parts[0]
	baseURL := "https://blockstream.info/api"
	if len(parts) > 1 {
		baseURL = parts[1]
	}
	return &BitcoinWalletAPI{
		Address: address,
		BaseURL: baseURL,
	}
}

func (b *BitcoinWalletAPI) HealthCheck() bool {
	_, err := b.FetchAssetAmount("BTC")
	return err == nil
}

func (b *BitcoinWalletAPI) FetchAssetAmount(symbol string) (float64, error) {
	if symbol != "BTC" {
		return 0, errors.New("only BTC supported")
	}

	url := fmt.Sprintf("%s/address/%s", b.BaseURL, b.Address)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("wallet API returned status: %d", resp.StatusCode)
	}

	var data struct {
		ChainStats struct {
			FundedTxoSum int64 `json:"funded_txo_sum"`
			SpentTxoSum  int64 `json:"spent_txo_sum"`
		} `json:"chain_stats"`
		MempoolStats struct {
			FundedTxoSum int64 `json:"funded_txo_sum"`
			SpentTxoSum  int64 `json:"spent_txo_sum"`
		} `json:"mempool_stats"`
	}
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, fmt.Errorf("failed to parse wallet data: %v", err)
	}

	totalSats := (data.ChainStats.FundedTxoSum - data.ChainStats.SpentTxoSum) +
		(data.MempoolStats.FundedTxoSum - data.MempoolStats.SpentTxoSum)

	return float64(totalSats) / 1e8, nil
}

func (b *BitcoinWalletAPI) AccountDetails() (string, error) {
	balance, err := b.FetchAssetAmount("BTC")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("BTC wallet address: %s\nBalance: %.8f BTC", b.Address, balance), nil
}

// No-op for wallet
func (b *BitcoinWalletAPI) Buy(options pnl.TradeOptions) (*pnl.Trade, error) {
	panic("Buy not supported for BTC wallet")
}

func (b *BitcoinWalletAPI) Sell(options pnl.TradeOptions) (*pnl.Trade, error) {
	panic("Sell not supported for BTC wallet")
}

func (b *BitcoinWalletAPI) ImportTrades(pair pnl.Pair, since time.Time) ([]pnl.Trade, error) {
	panic("ImportTrades not implemented for BTC wallet")
}

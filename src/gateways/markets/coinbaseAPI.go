package markets

import (
	"controtto/src/domain/pnl"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const COINBASE_HOST = "api.coinbase.com"

type CoinbaseMarketAPI struct {
	ApiKey    string
	ApiSecret string
}

type accountsResponse struct {
	Accounts []struct {
		UUID             string `json:"uuid"`
		Name             string `json:"name"`
		Currency         string `json:"currency"`
		AvailableBalance struct {
			Value    string `json:"value"`
			Currency string `json:"currency"`
		} `json:"available_balance"`
		Default   bool       `json:"default"`
		Active    bool       `json:"active"`
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt time.Time  `json:"updated_at"`
		DeletedAt *time.Time `json:"deleted_at"`
		Type      string     `json:"type"`
		Ready     bool       `json:"ready"`
		Hold      struct {
			Value    string `json:"value"`
			Currency string `json:"currency"`
		} `json:"hold"`
		RetailPortfolioID string `json:"retail_portfolio_id"`
		Platform          string `json:"platform"`
	} `json:"accounts"`
}

type transactionResponse struct {
	Pagination struct {
		EndingBefore  string `json:"ending_before"`
		StartingAfter string `json:"starting_after"`
		Limit         int    `json:"limit"`
		Order         string `json:"order"`
		PreviousURI   string `json:"previous_uri"`
		NextURI       string `json:"next_uri"`
	} `json:"pagination"`
	Transactions []struct {
		Amount struct {
			Amount   string `json:"amount"`
			Currency string `json:"currency"`
		}
		CreatedAt    time.Time `json:"created_at"`
		ID           string    `json:"id"`
		IDEMO        string    `json:"idem"`
		NativeAmount struct {
			Amount   string `json:"amount"`
			Currency string `json:"currency"`
		} `json:"native_amount"`
		Network struct {
			Hash           string `json:"hash"`
			NetworkName    string `json:"network_name"`
			Status         string `json:"status"`
			TransactionFee struct {
				Amount   string `json:"amount"`
				Currency string `json:"currency"`
			} `json:"transaction_fee"`
		} `json:"network"`
		Resource string `json:"resource"`
		Status   string `json:"status"`
		To       struct {
			Address  string `json:"address"`
			Resource string `json:"resource"`
		} `json:"to"`
		Buy struct {
			Fee struct {
				Amount   string `json:"amount"`
				Currency string `json:"currency"`
			} `json:"fee"`
			ID                string `json:"id"`
			PaymentMethodName string `json:"payment_method_name"`
			Subtotal          struct {
				Amount   string `json:"amount"`
				Currency string `json:"currency"`
			} `json:"subtotal"`
			Total struct {
				Amount   string `json:"amount"`
				Currency string `json:"currency"`
			} `json:"total"`
		} `json:"buy"`
		Sell struct {
			Fee struct {
				Amount   string `json:"amount"`
				Currency string `json:"currency"`
			} `json:"fee"`
			ID                string `json:"id"`
			PaymentMethodName string `json:"payment_method_name"`
			Subtotal          struct {
				Amount   string `json:"amount"`
				Currency string `json:"currency"`
			} `json:"subtotal"`
			Total struct {
				Amount   string `json:"amount"`
				Currency string `json:"currency"`
			} `json:"total"`
		} `json:"sell"`
		Type string `json:"type"`
	} `json:"data"`
}

func (a *accountsResponse) String() string {
	accountSummaries := []string{}
	for _, account := range a.Accounts {
		summary := fmt.Sprintf("Account: %s, Currency: %s, Available Balance: %s %s, Hold: %s %s", account.Name, account.Currency, account.AvailableBalance.Value, account.AvailableBalance.Currency, account.Hold.Value, account.Hold.Currency)
		accountSummaries = append(accountSummaries, summary)
	}
	return strings.Join(accountSummaries, "\n")
}

func NewCoinbaseAPI(token string) pnl.MarketAPI {
	if len(strings.Split(token, ":")) >= 2 {
		secret := strings.ReplaceAll(strings.TrimSpace(strings.Split(token, ":")[1]), "\\n", "\n")
		return &CoinbaseMarketAPI{
			ApiKey:    strings.TrimSpace(strings.Split(token, ":")[0]),
			ApiSecret: secret,
		}
	}
	return &CoinbaseMarketAPI{}
}

func (c *CoinbaseMarketAPI) HealthCheck() bool {
	_, err := c.getAllAccounts()
	if err != nil {
		log.Printf("Coinbase API health check failed: %v", err)
		return false
	}
	return true
}

func (c *CoinbaseMarketAPI) FetchAssetAmount(symbol string) (float64, error) {
	accounts, err := c.getAllAccounts()
	if err != nil {
		return 0, err
	}
	for _, account := range accounts.Accounts {
		if account.Currency == symbol {
			var amount float64
			_, err := fmt.Sscanf(account.AvailableBalance.Value, "%f", &amount)
			if err != nil {
				return 0, err
			}
			return amount, nil
		}
	}
	return 0, fmt.Errorf("asset %s not found", symbol)
}

func (c *CoinbaseMarketAPI) Buy(options pnl.TradeOptions) (*pnl.Trade, error) {
	return nil, errors.New("not implemented")
}
func (c *CoinbaseMarketAPI) Sell(options pnl.TradeOptions) (*pnl.Trade, error) {
	return nil, errors.New("not implemented")
}

func (c *CoinbaseMarketAPI) AccountDetails() (string, error) {
	accounts, err := c.getAllAccounts()
	if err != nil {
		return "", err
	}
	return accounts.String(), nil
}

func (c *CoinbaseMarketAPI) ImportTrades(tradingPair pnl.Pair, since time.Time) ([]pnl.Trade, error) {
	accountID, err := c.getAccountForBaseAsset(tradingPair.BaseAsset.Symbol)
	if err != nil {
		return nil, fmt.Errorf("error getting account for base asset %s: %v", tradingPair.BaseAsset.Symbol, err)
	}
	if accountID == "" {
		return nil, fmt.Errorf("no account found for base asset %s", tradingPair.BaseAsset.Symbol)
	}
	pairs, err := c.getTradesForAccount(accountID, tradingPair, since)
	if err != nil {
		return nil, fmt.Errorf("error getting trades for account %s: %v", accountID, err)
	}
	return pairs, nil
}

func (c *CoinbaseMarketAPI) getTradesForAccount(accountID string, tradingPair pnl.Pair, since time.Time) ([]pnl.Trade, error) {
	resp, err := c.makeRequest(http.MethodGet, fmt.Sprintf("/v2/accounts/%s/transactions", accountID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("coinbase transaction summary returned status: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var transactions transactionResponse
	if err := json.Unmarshal(body, &transactions); err != nil {
		return nil, err
	}
	pairs := []pnl.Trade{}
	for _, tx := range transactions.Transactions {
		if tx.CreatedAt.Before(since) {
			continue
		}
		if tx.Type != "buy" && tx.Type != "sell" {
			// We only care about buy/sell transactions
			continue
		}
		if tx.NativeAmount.Currency != tradingPair.QuoteAsset.Symbol || tx.Amount.Currency != tradingPair.BaseAsset.Symbol {
			continue
		}
		var quoteAmount float64
		var fee float64
		var tradeType pnl.TradeType
		var feeCurrency string
		switch tx.Type {
		case "sell":
			_, err := fmt.Sscanf(tx.Sell.Total.Amount, "%f", &quoteAmount)
			if err != nil {
				return nil, fmt.Errorf("error parsing base amount for tx %s: %v", tx.ID, err)
			}
			_, err = fmt.Sscanf(tx.Sell.Fee.Amount, "%f", &fee)
			if err != nil {
				return nil, fmt.Errorf("error parsing base amount for tx %s: %v", tx.ID, err)
			}
			feeCurrency = tx.Sell.Fee.Currency
			tradeType = pnl.Sell
		case "buy":
			_, err := fmt.Sscanf(tx.Buy.Total.Amount, "%f", &quoteAmount)
			if err != nil {
				return nil, fmt.Errorf("error parsing base amount for tx %s: %v", tx.ID, err)
			}
			_, err = fmt.Sscanf(tx.Buy.Fee.Amount, "%f", &fee)
			if err != nil {
				return nil, fmt.Errorf("error parsing base amount for tx %s: %v", tx.ID, err)
			}
			feeCurrency = tx.Buy.Fee.Currency
			tradeType = pnl.Buy
		}
		var baseAmount float64
		_, err = fmt.Sscanf(tx.Amount.Amount, "%f", &baseAmount)
		if err != nil {
			return nil, fmt.Errorf("error parsing base amount for tx %s: %v", tx.ID, err)
		}
		trade := pnl.Trade{
			ID:          tx.ID,
			Timestamp:   tx.CreatedAt,
			BaseAmount:  baseAmount,
			QuoteAmount: quoteAmount,
			FeeInBase:   0,
			FeeInQuote:  0,
			TradeType:   tradeType,
			Price:       quoteAmount / baseAmount,
		}
		switch feeCurrency {
		case tradingPair.BaseAsset.Symbol:
			trade.FeeInBase = fee
		case tradingPair.QuoteAsset.Symbol:
			trade.FeeInQuote = fee
		}
		pairs = append(pairs, trade)
	}
	return pairs, nil
}

func (c *CoinbaseMarketAPI) getAccountForBaseAsset(baseAsset string) (string, error) {
	accounts, err := c.getAllAccounts()
	if err != nil {
		return "", err
	}
	for _, account := range accounts.Accounts {
		if account.Currency == baseAsset {
			return account.UUID, nil
		}
	}
	return "", nil
}

func (c *CoinbaseMarketAPI) getAllAccounts() (*accountsResponse, error) {
	resp, err := c.makeRequest(http.MethodGet, "/api/v3/brokerage/accounts")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("coinbase accounts returned status: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var accounts accountsResponse
	if err := json.Unmarshal(body, &accounts); err != nil {
		return nil, err
	}
	return &accounts, nil
}

func (c *CoinbaseMarketAPI) makeRequest(method string, endpoint string) (*http.Response, error) {
	jwtToken, err := c.generateJWT(method, endpoint)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, "https://"+COINBASE_HOST+endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+jwtToken)

	client := &http.Client{Timeout: 10 * time.Second}
	return client.Do(req)
}

func (c *CoinbaseMarketAPI) generateJWT(method string, uri string) (string, error) {
	block, _ := pem.Decode([]byte(c.ApiSecret))
	if block == nil {
		return "", errors.New("failed to parse API secret as PEM")
	}

	privateKey, err := jwt.ParseECPrivateKeyFromPEM([]byte(c.ApiSecret))
	if err != nil {
		return "", errors.New("failed to parse private key")
	}

	claims := map[string]any{
		"sub": c.ApiKey,
		"iss": "cdp",
		"exp": time.Now().Add(time.Minute * 2).Unix(),
		"nbf": time.Now().Unix(),
		"uri": fmt.Sprintf("%s %s%s", method, COINBASE_HOST, uri),
	}

	nonce := fmt.Sprintf("%d", time.Now().UnixNano())
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims(claims))
	token.Header["kid"] = c.ApiKey
	token.Header["nonce"] = nonce

	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %v", err)
	}

	return signedToken, nil
}

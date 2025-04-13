package config

import (
	"controtto/src/domain/pnl"
	"controtto/src/gateways/marketTraders"
	"controtto/src/gateways/priceProviders"
	"fmt"
	"maps"
	"os"
	"strings"
	"sync"
)

const (
	PORT           = "CONTROTTO_PORT"
	UNCOMMON_PAIRS = "CONTROTTO_UNCOMMON_PAIRS"
)

type Config struct {
	mu sync.RWMutex

	Port           string
	UncommonPairs  bool
	PriceProviders map[string]pnl.PriceProvider
	MarketTraders  map[string]pnl.MarketTrader
}

// Load initializes the configuration from environment variables.
// If PORT is not set, it defaults to "8000".
func Load() *Config {
	cfg := &Config{
		Port:          os.Getenv(PORT),
		UncommonPairs: os.Getenv(UNCOMMON_PAIRS) == "true",
	}
	// Load tokens from environment or set defaults
	cfg.PriceProviders = loadProviders()
	cfg.MarketTraders = loadMarketTraders()

	if cfg.Port == "" {
		cfg.Port = "8000"
		os.Setenv(PORT, "8000")
	}
	return cfg
}
func loadMarketTraders() map[string]pnl.MarketTrader {
	trading212 := os.Getenv("CONTROTTO_TRADING212_TRADER_TOKEN")
	pancake := os.Getenv("CONTROTTO_PANCAKE_TRADER_TOKEN")
	binance := os.Getenv("CONTROTTO_BINANCE_TRADER_TOKEN")
	return map[string]pnl.MarketTrader{
		"pancake": {
			IsSet:       pancake != "",
			Env:         "CONTROTTO_PANCAKE_TRADER_TOKEN",
			MarketName:  "Pancake",
			MarketKey:   "pancake_trader",
			Color:       "#6D0000",
			Type:        pnl.DEX,
			Token:       pancake,
			ProviderURL: "https://controtto.io/docs/pancake",
			MarketLogo:  "/assets/img/marketPancake.png",
			MarketAPI:   marketTraders.NewMockMarketAPI(pancake),
		},
		"binance": {
			IsSet:       binance != "",
			Env:         "CONTROTTO_BINANCE_TRADER_TOKEN",
			MarketName:  "Binance",
			MarketKey:   "binance_trader",
			Color:       "#6D0000",
			Type:        pnl.Exchange,
			Token:       binance,
			ProviderURL: "https://controtto.io/docs/demo",
			MarketLogo:  "/assets/img/marketBinance.png",
			MarketAPI:   marketTraders.NewMockMarketAPI(binance),
		},
		"trading212": {
			IsSet:       trading212 != "",
			Env:         "CONTROTTO_TRADING212_TRADER_TOKEN",
			MarketName:  "Trading212",
			MarketKey:   "trading212_trader",
			Color:       "#6D0000",
			Type:        pnl.Broker,
			Token:       trading212,
			ProviderURL: "https://controtto.io/docs/demo",
			MarketLogo:  "/assets/img/marketTrading212.png",
			MarketAPI:   marketTraders.NewMockMarketAPI(trading212),
		},
	}
}

// loadProviders reads the tokens from the environment variables and returns a map of priceProviders.PriceProviders structs.
func loadProviders() map[string]pnl.PriceProvider {
	tingoToken := os.Getenv("CONTROTTO_TIINGO_TOKEN")
	avantageToken := os.Getenv("CONTROTTO_AVANTAGE_TOKEN")
	coinbaseEnable := os.Getenv("CONTROTTO_COINBASE_ENABLE")
	binanceEnable := os.Getenv("CONTROTTO_BINANCE_ENABLE")
	bingxEnable := os.Getenv("CONTROTTO_BINGX_ENABLE")
	return map[string]pnl.PriceProvider{
		"bingx": {
			IsSet:        bingxEnable == "",
			Env:          "CONTROTTO_BINGX_ENABLE",
			ProviderName: "BingX",
			ProviderURL:  "https://docs.cdp.coinbase.com/",
			NeedsToken:   false,
			ProviderLogo: "/assets/img/marketBingx.png",
			Color:        "#2954FE",
			ProviderKey:  "bingx_toggle",
			PriceAPI:     priceProviders.NewBingxAPI(),
		},
		"binance": {
			IsSet:        binanceEnable == "", // No token needed for Binance
			Env:          "CONTROTTO_BINANCE_ENABLE",
			ProviderName: "Binance",
			ProviderURL:  "https://docs.cdp.coinbase.com/",
			ProviderKey:  "binance_toggle",
			NeedsToken:   false,
			Color:        "#EFB72D",
			ProviderLogo: "/assets/img/marketBinance.png",
			PriceAPI:     priceProviders.NewBinanceAPI(),
		},
		"coinbase": {
			IsSet:        coinbaseEnable == "",
			Env:          "CONTROTTO_COINBASE_ENABLE",
			ProviderName: "Coinbase",
			ProviderURL:  "https://docs.cdp.coinbase.com/",
			ProviderKey:  "coinbase_toggle",
			NeedsToken:   false,
			ProviderLogo: "/assets/img/marketCoinbase.png",
			Token:        "",
			Color:        "#0052FF",
			PriceAPI:     priceProviders.NewCoinbaseAPI(),
		},
		"avantage": {
			IsSet:        avantageToken != "",
			Env:          "CONTROTTO_AVANTAGE_TOKEN",
			ProviderName: "Alpha Vantage",
			NeedsToken:   true,
			ProviderURL:  "https://www.alphavantage.co/support/#api-key",
			ProviderKey:  "vantage_token",
			ProviderLogo: "/assets/img/priceVantage.png",
			Token:        avantageToken,
			Color:        "#C2F4E1",
			PriceAPI:     priceProviders.NewAVantageAPI(avantageToken),
		},
		"tiingo": {
			IsSet:        tingoToken != "",
			Env:          "CONTROTTO_TIINGO_TOKEN",
			NeedsToken:   true,
			ProviderName: "Tiingo",
			ProviderURL:  "https://www.tiingo.com/account/api/token",
			ProviderKey:  "tiingo_token",
			ProviderLogo: "/assets/img/priceTiingo.png",
			Token:        tingoToken,
			Color:        "#AA74EF",
			PriceAPI:     priceProviders.NewTiingoAPI(tingoToken),
		},
	}
}

// GetPort returns the configured port.
func (c *Config) GetPort() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Port
}

// SetUncommonPairs sets or unsets the uncommon pairs flag and updates the env var.
func (c *Config) SetUncommonPairs(enabled bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.UncommonPairs = enabled
	if enabled {
		os.Setenv(UNCOMMON_PAIRS, "true")
	} else {
		os.Setenv(UNCOMMON_PAIRS, "false")
	}
}

// GetUncommonPairs returns whether uncommon pairs are enabled.
func (c *Config) GetUncommonPairs() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.UncommonPairs
}

func (c *Config) GetMarketTraders() map[string]pnl.MarketTrader {
	c.mu.RLock()
	defer c.mu.RUnlock()
	// Return a copy to prevent external modification
	traders := make(map[string]pnl.MarketTrader, len(c.MarketTraders))
	maps.Copy(traders, c.MarketTraders)
	return traders
}

func (c *Config) UpdateMarketTraderToken(key, token string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	trader, ok := c.MarketTraders[key]
	if !ok {
		return fmt.Errorf("market trader %q not found", key)
	}
	switch key {
	case "binance", "trading212", "pancake":
		trader.Token = token
		trader.MarketAPI = marketTraders.NewMockMarketAPI(token)
		os.Setenv("CONTROTTO_"+strings.ToUpper(key)+"_TRADER_TOKEN", token)
		trader.IsSet = token != ""
	default:
		return fmt.Errorf("unsupported market trader %q", key)
	}
	c.MarketTraders[key] = trader
	return nil
}

func (c *Config) GetPriceProviders() map[string]pnl.PriceProvider {
	c.mu.RLock()
	defer c.mu.RUnlock()
	// Return a copy to prevent external modification
	providers := make(map[string]pnl.PriceProvider, len(c.PriceProviders))
	maps.Copy(providers, c.PriceProviders)
	return providers
}

func (c *Config) UpdateProviderToken(key, token string, toggle bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	provider, ok := c.PriceProviders[key]
	if !ok {
		return fmt.Errorf("price provider %q not found", key)
	}
	// Update the provider based on the key
	switch key {
	case "avantage":
		provider.Token = token
		provider.PriceAPI = priceProviders.NewAVantageAPI(token)
		os.Setenv("CONTROTTO_AVANTAGE_TOKEN", token)
		// Set IsSet to true if either toggle is true or a non-empty token is provided
		provider.IsSet = toggle || token != ""
	case "tiingo":
		provider.Token = token
		provider.PriceAPI = priceProviders.NewTiingoAPI(token)
		os.Setenv("CONTROTTO_TIINGO_TOKEN", token)
		// Set IsSet to true if either toggle is true or a non-empty token is provided
		provider.IsSet = toggle || token != ""
	case "bingx", "binance", "coinbase":
		// For providers that don't need a token, IsSet is controlled by toggle only
		provider.IsSet = toggle
	default:
		return fmt.Errorf("unsupported price provider %q", key)
	}
	// Save the updated provider back to the map
	c.PriceProviders[key] = provider
	return nil
}

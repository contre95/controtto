package config

import (
	"controtto/src/domain/pnl"
	"controtto/src/gateways/priceProviders"
	"fmt"
	"maps"
	"os"
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

	if cfg.Port == "" {
		cfg.Port = "8000"
		os.Setenv(PORT, "8000")
	}
	return cfg
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
			IsSet:             bingxEnable == "",
			Env:               "CONTROTTO_BINGX_ENABLE",
			ProviderName:      "BingX",
			ProviderURL:       "https://docs.cdp.coinbase.com/",
			NeedsToken:        false,
			Color:             "#2954FE",
			ProviderInputName: "bingx_toggle",
			API:               priceProviders.NewCoinbaseAPI(),
		},
		"binance": {
			IsSet:             binanceEnable == "", // No token needed for Binance
			Env:               "CONTROTTO_BINANCE_ENABLE",
			ProviderName:      "Binance",
			ProviderURL:       "https://docs.cdp.coinbase.com/",
			ProviderInputName: "binance_toggle",
			NeedsToken:        false,
			Color:             "#EFB72D",
			API:               priceProviders.NewCoinbaseAPI(),
		},
		"coinbase": {
      IsSet:             coinbaseEnable == "", 
			Env:               "CONTROTTO_COINBASE_ENABLE",
			ProviderName:      "Coinbase",
			ProviderURL:       "https://docs.cdp.coinbase.com/",
			ProviderInputName: "coinbase_toggle",
			NeedsToken:        false,
			Token:             "",
			Color:             "#0052FF",
			API:               priceProviders.NewCoinbaseAPI(),
		},
		"avantage": {
      IsSet:             coinbaseEnable == "",
			Env:               "CONTROTTO_AVANTAGE_TOKEN",
			ProviderName:      "Alpha Vantage",
			NeedsToken:        true,
			ProviderURL:       "https://www.alphavantage.co/support/#api-key",
			ProviderInputName: "vantage_token",
			Token:             avantageToken,
			Color:             "#C2F4E1",
			API:               priceProviders.NewAVantageAPI(avantageToken),
		},
		"tiingo": {
			IsSet:             tingoToken != "",
			Env:               "CONTROTTO_TIINGO_TOKEN",
			NeedsToken:        true,
			ProviderName:      "Tiingo",
			ProviderURL:       "https://www.tiingo.com/account/api/token",
			ProviderInputName: "tiingo_token",
			Token:             tingoToken,
			Color:             "#AA74EF",
			API:               priceProviders.NewTiingoAPI(tingoToken),
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

// GetPriceProviders returns the map of price providers.
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
		provider.API = priceProviders.NewAVantageAPI(token)
		os.Setenv("CONTROTTO_AVANTAGE_TOKEN", token)
		// Set IsSet to true if either toggle is true or a non-empty token is provided
		provider.IsSet = toggle || token != ""
	case "tiingo":
		provider.Token = token
		provider.API = priceProviders.NewTiingoAPI(token)
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

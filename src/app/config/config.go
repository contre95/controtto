package config

import (
	"controtto/src/domain/pnl"
	"controtto/src/gateways/priceProviders"
	"fmt"
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
	avantageToken := os.Getenv("CONTROTTO_AVANTAGE_TOKEN")
	tingoToken := os.Getenv("CONTROTTO_TIINGO_TOKEN")
	return map[string]pnl.PriceProvider{
		"avantage": {
			TokenSet:          avantageToken != "",
			Env:               []string{"CONTROTTO_AVANTAGE_TOKEN"},
			ProviderName:      "Alpha Vantage",
			ProviderURL:       "https://www.alphavantage.co/support/#api-key",
			ProviderInputName: "vantage_token",
			Token:             avantageToken,
			Color:             "",
			API:               priceProviders.NewAVantageAPI(avantageToken),
		},
		"tiingo": {
			TokenSet:          tingoToken != "",
			Env:               []string{"CONTROTTO_TIINGO_TOKEN"},
			ProviderName:      "Tiingo",
			ProviderURL:       "https://www.tiingo.com/account/api/token",
			ProviderInputName: "tiingo_token",
			Token:             tingoToken,
			Color:             "",
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
	for k, v := range c.PriceProviders {
		providers[k] = v
	}
	return providers
}

// UpdateProviderToken updates the token for the specified price provider and its corresponding environment variable.
func (c *Config) UpdateProviderToken(key, token string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	provider, ok := c.PriceProviders[key]
	if !ok {
		return fmt.Errorf("price provider %q not found", key)
	}

	// Update the provider's token and related fields
	provider.Token = token
	provider.TokenSet = token != ""
	switch key {
	case "avantage":
		provider.API = priceProviders.NewAVantageAPI(token)
		os.Setenv("CONTROTTO_AVANTAGE_TOKEN", token)
	case "tiingo":
		provider.API = priceProviders.NewTiingoAPI(token)
		os.Setenv("CONTROTTO_TIINGO_TOKEN", token)
	default:
		return fmt.Errorf("unsupported price provider %q", key)
	}

	// Save the updated provider back to the map
	c.PriceProviders[key] = provider
	return nil
}

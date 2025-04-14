package config

import (
	"controtto/src/domain/pnl"
	"os"
	"sync"
)

const (
	PORT                = "CONTROTTO_PORT"
	UNCOMMON_PAIRS      = "CONTROTTO_UNCOMMON_PAIRS"
	PREFIX              = "CONTROTTO_"
	TRADER_SUFIX        = "_TRADER_TOKEN"
	PRIVATE_PRICE_SUFIX = "_PRICER_TOKEN"
	PUBLIC_PRICE_SUFIX  = "_PRICER_ENABLED"
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
	cfg.PriceProviders = loadProviders()
	cfg.MarketTraders = loadMarketTraders()

	if cfg.Port == "" {
		cfg.Port = "8000"
		os.Setenv(PORT, "8000")
	}
	return cfg
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

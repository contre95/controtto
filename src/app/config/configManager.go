package config

import (
	"os"
	"strings"
	"sync"
)

const (
	PORT           = "PORT"
	DEFAULT_PORT   = "8000"
	UNCOMMON_PAIRS = "UNCOMMON_PAIRS"
)

type ConfigManager struct {
	prefix        string
	mu            sync.RWMutex
	Port          string
	uncommonPairs bool
}

// NewConfig initializes the configuration from environment variables.
// If PORT is not set, it defaults to "8000".
func NewConfig(prefix string) *ConfigManager {
	cfg := &ConfigManager{
		Port:          os.Getenv(strings.ToUpper(prefix) + PORT),
		uncommonPairs: os.Getenv(strings.ToUpper(prefix)+UNCOMMON_PAIRS) == "true",
	}
	cfg.prefix = strings.ToUpper(prefix)
	if cfg.Port == "" {
		cfg.Port = DEFAULT_PORT
		os.Setenv(PORT, DEFAULT_PORT)
	}
	return cfg
}

// GetPort returns the configured port.
func (c *ConfigManager) GetPort() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Port
}

func (c *ConfigManager) SetUncommonPairs(enabled bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.uncommonPairs = enabled
	if enabled {
		os.Setenv(c.prefix+UNCOMMON_PAIRS, "true")
	} else {
		os.Setenv(c.prefix+UNCOMMON_PAIRS, "false")
	}
}

// GetUncommonPairs returns whether uncommon pairs are enabled.
func (c *ConfigManager) GetUncommonPairs() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.uncommonPairs
}

package config

import "sync"

// Manager holds the application configuration and provides thread-safe access to it.
type Manager struct {
	mu     sync.RWMutex
	config *Config
}

// NewManager creates a new ConfigManager.
func NewManager(config *Config) *Manager {
	return &Manager{config: config}
}

// Get returns the current configuration.
func (m *Manager) Get() *Config {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.config
}

// Update updates the configuration.
func (m *Manager) Update(config *Config) {
	// If Config is updated on runtime, we don't want to change the env vars. Config will be re-loaded to env var value on restart.
	m.mu.Lock()
	defer m.mu.Unlock()
	// We are re-setting these values to their original ones, as they are not expected to be changed at runtime.
	config.DBPath = m.config.DBPath
	config.Port = m.config.Port
	config.LoadSampleData = m.config.LoadSampleData
	m.config = config
}

package config

// Service just hols all the managing use cases
type Service struct {
	ConfigManager *ConfigManager
}

// NewService is the interctor for all Managing Use cases
func NewService(cfg *ConfigManager) Service {
	return Service{cfg}
}

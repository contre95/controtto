package config

import (
	"log/slog"
	"os"
	"strconv"
)

const (
	DefaultPort           = 8000
	DefaultDBPath         = "./data/pnl.db"
	DefaultUncommonPairs  = false
	DefaultLoadSampleData = false
)

func LoadFromEnv() (*Config, error) {
	port, err := getEnvInt("CONTROTTO_PORT", DefaultPort)
	if err != nil {
		return nil, err
	}

	dbPath := getEnvString("CONTROTTO_DB_PATH", DefaultDBPath)

	uncommonPairs, err := getEnvBool("CONTROTTO_UNCOMMON_PAIRS", DefaultUncommonPairs)
	if err != nil {
		return nil, err
	}

	loadSampleData, err := getEnvBool("CONTROTTO_LOAD_SAMPLE_DATA", DefaultLoadSampleData)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Port:           port,
		DBPath:         dbPath,
		UncommonPairs:  uncommonPairs,
		LoadSampleData: loadSampleData,
	}

	return cfg, nil
}

func getEnvString(key string, defaultVal string) string {
	val := os.Getenv(key)
	slog.Info("Loading env var", "key", key, "value", val)
	if val == "" {
		return defaultVal
	}
	return val
}

func getEnvBool(key string, defaultVal bool) (bool, error) {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal, nil
	}
	val, err := strconv.ParseBool(valStr)
	if err != nil {
		return false, err
	}
	return val, nil
}

func getEnvInt(key string, defaultVal int) (int, error) {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal, nil
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return 0, err
	}
	return val, nil
}

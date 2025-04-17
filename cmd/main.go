package main

import (
	"controtto/src/app/config"
	"controtto/src/app/managing"
	"controtto/src/app/querying"
	"controtto/src/gateways/sqlite"
	"controtto/src/rest"
	"log/slog"
	"os"
)

const (
	PREFIX              = "CONTROTTO_"
	TRADER_PREFIX       = "CONTROTTO_TRADER_"
	PRICER_PREFIX       = "CONTROTTO_PRICER_"
	TRADER_SUFIX        = "_TOKEN"
	PRIVATE_PRICE_SUFIX = "_TOKEN"
	PUBLIC_PRICE_SUFIX  = "_ENABLED"
)

func main() {
	// Database
	dbPath := "pnl.db"
	if dbPathEnv := os.Getenv("CONTROTTO_DB_PATH"); dbPathEnv != "" {
		dbPath = dbPathEnv
	}
	slog.Info("Initiating SQLite path", "path", dbPath)
	sqliteDB, err := sqlite.NewSQLite(dbPath)
	if err != nil {
		slog.Error("Error creating SQLite:", "error", err)
		panic("Bye")
	}

	// Load configuration
	cfg := config.NewConfig(PREFIX)
	ac := managing.NewAssetCreator(sqliteDB)
	tpc := managing.NewTradingPairManager(cfg, sqliteDB, sqliteDB)
	tpq := querying.NewTradingPairQuerier(sqliteDB)
	mtm := managing.NewMarketManager(traders)
	ppm := managing.NewPriceProviderManager(pricers)
	manager := managing.NewService(*ac, *tpc, mtm, ppm)
	aq := querying.NewAssetQuerier(sqliteDB)
	config := config.NewService(cfg)
	querier := querying.NewService(*aq, *tpq)

	port := cfg.Port
	slog.Info("Initiating server", "port", port)
	rest.Run(&config, &manager, &querier, port)
}

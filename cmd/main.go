package main

import (
	"controtto/src/app/config"
	"controtto/src/app/managing"
	"controtto/src/app/querying"
	"controtto/src/app/trading"
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
	config := config.NewService(cfg)
	ac := managing.NewAssetCreator(sqliteDB)
	pc := managing.NewPairManager(cfg, sqliteDB, sqliteDB)
	pq := querying.NewPairQuerier(sqliteDB)
	mm := managing.NewMarketManager(traders)
	ppm := managing.NewPriceProviderManager(pricers)
	aq := querying.NewAssetQuerier(sqliteDB)
	manager := managing.NewService(*ac, *pc, mm, ppm)
	querier := querying.NewService(*aq, *pq)
	tr := trading.NewTradeRecorder(sqliteDB)
	at := trading.NewAssetTrader(mm, sqliteDB)
	trader := trading.NewService(*at, *tr)
	// *trading.NewAssetTrader(markets map[string]pnl.Market)

	port := cfg.Port
	slog.Info("Initiating server", "port", port)
	rest.Run(&config, &manager, &querier, &trader)
}

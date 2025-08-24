package main

import (
	"controtto/src/app/config"
	"controtto/src/app/managing"
	"controtto/src/app/querying"
	"controtto/src/app/trading"
	"controtto/src/gateways/sqlite"
	"controtto/src/rest"
	"log/slog"
)

const (
	PREFIX              = "CONTROTTO_"
	MARKET_PREFIX       = "CONTROTTO_TRADER_"
	PRICER_PREFIX       = "CONTROTTO_PRICER_"
	MARKET_SUFIX        = "_TOKEN"
	PRIVATE_PRICE_SUFIX = "_TOKEN"
	PUBLIC_PRICE_SUFIX  = "_ENABLED"
)

func main() {
	// Load configuration
	cfg, err := config.LoadFromEnv()
	if err != nil {
		slog.Error("Error loading config from env:", "error", err)
		panic("Bye")
	}
	cfgManager := config.NewManager(cfg)

	// Database
	slog.Info("Initiating SQLite path", "path", cfgManager.Get().DBPath)
	sqliteDB, err := sqlite.NewSQLite(cfgManager)
	if err != nil {
		slog.Error("Error creating SQLite:", "error", err)
		panic("Bye")
	}

	// Use Cases
	ac := managing.NewAssetCreator(sqliteDB)
	pc := managing.NewPairManager(cfgManager, sqliteDB, sqliteDB)
	pq := querying.NewPairQuerier(sqliteDB)
	mm := managing.NewMarketManager(marketsConfig)
	ppm := managing.NewPriceProviderManager(pricers)
	aq := querying.NewAssetQuerier(sqliteDB)
	manager := managing.NewService(*ac, *pc, mm, ppm)
	querier := querying.NewService(*aq, *pq)
	tr := trading.NewTradeRecorder(sqliteDB)
	at := trading.NewAssetTrader(mm, sqliteDB)
	trader := trading.NewService(*at, *tr)

	port := cfg.Port
	slog.Info("Initiating server", "port", port)
	rest.Run(cfgManager, &manager, &querier, &trader)
}

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
	cfg := config.Load()
	ac := managing.NewAssetCreator(sqliteDB)
	tpc := managing.NewTradingPairManager(cfg, sqliteDB, sqliteDB)
	mtm := managing.NewMarketManager(cfg)
	mkq := querying.NewPriceQuerier(cfg)
	tpq := querying.NewTradingPairQuerier(cfg, sqliteDB)
	manager := managing.NewService(*ac, *tpc, *mtm)
	aq := querying.NewAssetQuerier(sqliteDB)
	config := config.NewService(cfg)
	querier := querying.NewService(*aq, *mkq, *tpq)

	port := cfg.Port
	slog.Info("Initiating server", "port", port)
	rest.Run(&config, &manager, &querier, port)
}

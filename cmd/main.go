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
	manager := managing.NewService(*ac, *tpc)
	aq := querying.NewAssetQuerier(sqliteDB)
	mkq := querying.NewPriceQuerier(cfg.PriceProviders)
	tpq := querying.NewTradingPairQuerier(sqliteDB, cfg.PriceProviders)
	querier := querying.NewService(*aq, *mkq, *tpq)

	port := cfg.Port
	slog.Info("Initiating server", "port", port)
	rest.Run(cfg, &manager, &querier)
}

package main

import (
	"controtto/src/app/managing"
	"controtto/src/app/querying"
	"controtto/src/gateways/markets"
	"controtto/src/gateways/sqlite"
	"controtto/src/presenters"
	"log/slog"
	"os"
)

func main() {
	port := "8000"
	portEnv := ""
	portEnv = os.Getenv("CONTROTTO_PORT")
	if len(portEnv) > 0 {
		port = portEnv
	}
	dbPath := "pnl.db"
	dbPathEnv := os.Getenv("CONTROTTO_DB_PATH")
	if len(portEnv) > 0 {
		dbPath = dbPathEnv
	}
	sqlite, err := sqlite.NewSQLite(dbPath)
	if err != nil {
		slog.Error("Error creating SQLite:", "error", err)
		panic("Bye")
	}

	binanceAPI := markets.NewBinanceAPI()
	ac := managing.NewAssetCreator(sqlite)
	tpc := managing.NewTradingPairManager(sqlite, sqlite)
	aq := querying.NewAssetQuerier(sqlite)
	mkq := querying.NewMarketQuerier(binanceAPI)
	tpq := querying.NewTradingPairQuerier(sqlite, binanceAPI)
	querier := querying.NewService(*aq, *mkq, *tpq)
	manager := managing.NewService(*ac, *tpc)
	presenters.Run(port, &manager, &querier)
}

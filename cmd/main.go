package main

import (
	"controtto/src/app/managing"
	"controtto/src/app/querying"
	"controtto/src/domain/pnl"
	"controtto/src/gateways/markets"
	"controtto/src/gateways/sqlite"
	"controtto/src/presenters"
	"log/slog"
	"os"
)

func main() {
	// Database
	dbPath := "pnl.db"
	dbPathEnv := os.Getenv("CONTROTTO_DB_PATH")
	if len(dbPathEnv) > 0 {
		dbPath = dbPathEnv
	}
	slog.Info("Initiating SQLite path", "path", dbPath)
	sqlite, err := sqlite.NewSQLite(dbPath)
	if err != nil {
		slog.Error("Error creating SQLite:", "error", err)
		panic("Bye")
	}

	// Markets
	binanceAPI := markets.NewBinanceAPI()
	bingxAPI := markets.NewBingxAPI()
	marketsAPIs := []pnl.Markets{binanceAPI, bingxAPI} // Defines query order
	avToken := os.Getenv("CONTROTTO_AVANTAGE_TOKEN")
	if len(avToken) != 0 {
		avantageAPI := markets.NewAVantageAPI(avToken)
		marketsAPIs = append(marketsAPIs, avantageAPI)
	}
	tToken := os.Getenv("CONTROTTO_TIINGO_TOKEN")
	if len(tToken) != 0 {
		tiingoAPI := markets.NewTiingoAPI(tToken)
		marketsAPIs = append(marketsAPIs, tiingoAPI)
	}
	for _, m := range marketsAPIs {
		slog.Info("Market registered", "market", m.Name())
	}
	ac := managing.NewAssetCreator(sqlite)
	tpc := managing.NewTradingPairManager(sqlite, sqlite)
	aq := querying.NewAssetQuerier(sqlite)
	mkq := querying.NewMarketQuerier(marketsAPIs)
	tpq := querying.NewTradingPairQuerier(sqlite, marketsAPIs)
	querier := querying.NewService(*aq, *mkq, *tpq)
	manager := managing.NewService(*ac, *tpc)

	port := "8000"
	portEnv := os.Getenv("CONTROTTO_PORT")
	if len(portEnv) > 0 {
		port = portEnv
	}
	slog.Info("Initiating server", "port", port)
	presenters.Run(port, &manager, &querier)
}

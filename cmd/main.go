package main

import (
	"controtto/src/app/config"
	"controtto/src/app/managing"
	"controtto/src/app/querying"
	"controtto/src/domain/pnl"
	"controtto/src/gateways/markets"
	"controtto/src/gateways/sqlite"
	"controtto/src/rest"
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
	cfg := config.NewConfig()
	// Markets
	binanceAPI := markets.NewBinanceAPI()
	bingxAPI := markets.NewBingxAPI()
	// marketsAPIs := []pnl.Markets{binanceAPI, bingxAPI} // Defines query order
	marketsAPIs := []pnl.Markets{bingxAPI, binanceAPI} // Defines query order
	if cfg.IsSet().AVantageAPIToken {
		avantageAPI := markets.NewAVantageAPI(cfg.Get().AVantageAPIToken)
		marketsAPIs = append(marketsAPIs, avantageAPI)
	}
	if cfg.IsSet().TiingoAPIToken {
		tiingoAPI := markets.NewTiingoAPI(cfg.Get().TiingoAPIToken)
		marketsAPIs = append(marketsAPIs, tiingoAPI)
	}
	for _, m := range marketsAPIs {
		slog.Info("Market registered", "market", m.Name())
	}
	ac := managing.NewAssetCreator(sqlite)
	tpc := managing.NewTradingPairManager(sqlite, sqlite)
	manager := managing.NewService(*ac, *tpc)
	aq := querying.NewAssetQuerier(sqlite)
	mkq := querying.NewMarketQuerier(marketsAPIs)
	tpq := querying.NewTradingPairQuerier(sqlite, marketsAPIs)
	querier := querying.NewService(*aq, *mkq, *tpq)

	port := cfg.Get().Port
	slog.Info("Initiating server", "port", port)
	rest.Run(cfg, &manager, &querier)
}

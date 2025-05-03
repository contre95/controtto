package rest

import (
	"controtto/src/app/config"
	"controtto/src/app/managing"
	"controtto/src/app/querying"
	"controtto/src/app/trading"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func Run(c *config.Service, m *managing.Service, q *querying.Service, t *trading.Service) {
	engine := html.New("./views", ".html")
	engine.Debug(true)
	app := fiber.New(fiber.Config{
		Views: engine,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		},
	})
	app.Static("/assets", "./public/assets")

	// GET
	app.Get("/", Home)

	// Sidebar
	app.Get("/dashboard", dashboardSection(q.PairQuerier))
	app.Get("/pairs", pairsSection)
	app.Get("/algoTrading", algoTradingSection)
	app.Get("/settings", settingsSection(m.PriceProviderManager, m.MarketManager, c.ConfigManager))

	// Strategies
	app.Get("/ui/algoTrading/spatialArbitrageForm", getSpatialArbitrageForm(q.PairQuerier, m.MarketManager))
	app.Get("/strategy/:id/logs", getRunningJobLogs())
	app.Post("/strategy/:id/start", startStrategyJob())
	app.Post("/strategy/:id/stop", startStrategyJob())

	// Pairs
	app.Get("/ui/pairs/form", newPairForm(q.AssetQuerier))
	app.Get("/ui/pairs/table", pairsTable(q.PairQuerier))
	app.Delete("/pairs/:id", deletePair(m.PairManager))
	app.Post("/assets", newAsset(m.AssetManager))
	app.Post("/pairs", newPair(m.PairManager))

	// Pair
	app.Get("/pairs/:id/", pairSection())
	app.Get("/ui/pairs/:id/cards", pairCards(q.PairQuerier, m.PriceProviderManager))
	app.Get("/ui/pairs/:id/chart", pairChart(q.PairQuerier))
	app.Get("/pairs/:id/market/:mktkey/trades", fetchMarketTrades(t.AssetTrader, m.MarketManager))
	app.Get("/ui/pairs/:id/tape", pairTape(q.PairQuerier))
	app.Get("/pairs/:id/trades/export", tradingExport(q.PairQuerier))
	app.Get("/ui/pairs/:id/trades/table", tradingTable(q.PairQuerier))
	app.Get("/ui/pairs/:id/newTrade/form", newTradeForm(q.PairQuerier, m.PriceProviderManager))
	app.Get("/ui/pairs/:id/market", getMarketAssets(q.PairQuerier, m.MarketManager))
	app.Get("/ui/pairs/:id/newMarketTrade/form", newMarketTradingForm(q.PairQuerier, m.MarketManager))
	app.Post("/pairs/:id/trades", newTrade(t.TradeRecorder))
	app.Post("/pairs/:id/trades/mktImport", importMarketTrades(t.TradeRecorder))
	app.Post("/pairs/:id/trades/csvImport", newTradeImport(t.TradeRecorder))
	app.Delete("/trades/:id", deleteTrade(t.TradeRecorder))

	// Assets
	app.Get("/ui/assets/form", newAssetForm)

	// Misc
	app.Get("/healthcheck/price", checkPrice(m.PriceProviderManager))
	app.Get("/asset/:symbol/color", getAssetColor(q.AssetQuerier))
	app.Get("/calc/price", calculatePrice(m.PriceProviderManager))
	app.Get("/pairs/AvgBuyPrice/:id", avgBuyPrice(q.PairQuerier))
	app.Delete("/empty", empty())

	// Settings
	app.Post("/settings/edit", saveSettingsForm(m.PriceProviderManager, m.MarketManager, c.ConfigManager))

	// Start
	log.Fatal(app.Listen("0.0.0.0" + ":" + c.ConfigManager.GetPort()))
}

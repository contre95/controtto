package rest

import (
	"controtto/src/app/config"
	"controtto/src/app/managing"
	"controtto/src/app/querying"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func Run(c *config.Service, m *managing.Service, q *querying.Service, port string) {
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
	app.Get("/trades", tradesSection)
	app.Get("/ui/trades/table", tradesTable(q.PairQuerier))
	app.Get("/pairs/:id/", pairSection())
	app.Get("/ui/assets/form", newAssetForm)
	app.Get("/ui/pairs/form", newPairForm(q.AssetQuerier))
	app.Get("/healthcheck/price", checkPrice(m.PriceProviderManager))
	app.Get("/calc/price", calculatePrice(m.PriceProviderManager))
	app.Get("/dashboard", dashboardSection(q.PairQuerier))
	app.Get("/ui/pairs/:id/cards", pairCards(q.PairQuerier, m.PriceProviderManager))
	app.Get("/ui/pairs/:id/chart", pairChart(q.PairQuerier))
	app.Get("/ui/pairs/:id/tape", pairTape(q.PairQuerier))
	app.Get("/pairs/AvgBuyPrice/:id", avgBuyPrice(q.PairQuerier))
	app.Get("/pairs/:id/trades/export", tradingExport(q.PairQuerier))
	app.Get("/ui/pairs/:id/trades/table", tradingTable(q.PairQuerier))
	app.Get("/ui/pairs/:id/newTrade/form", newTradeForm(q.PairQuerier, m.PriceProviderManager))
	app.Get("/settings", settingsSection(m.PriceProviderManager, m.MarketManager, c.ConfigManager))
	// app.Get("/settings/anyMarket", marketsSetAPI(c.ConfigManager))
	// app.Get("/markets", marketsSection(m.MarketManager))
	app.Get("/ui/pairs/:id/market", getMarketAssets(q.PairQuerier, m.MarketManager))
	app.Get("/ui/pairs/:id/newMarketTrade/form", newMarketTradingForm(q.PairQuerier, m.MarketManager))
	// DELETE
	app.Delete("/empty", empty())
	app.Delete("/pairs/:id", deletePair(m.PairManager))
	app.Delete("/trades/:id", deleteTrade(m.PairManager))
	// POST
	app.Post("/assets", newAsset(m.AssetManager))
	app.Post("/pairs", newPair(m.PairManager))
	app.Post("/pairs/:id/trades", newTrade(m.PairManager))
	app.Post("/pairs/:id/trades/upload", newTradeImport(m.PairManager))
	app.Post("/settings/edit", saveSettingsForm(m.PriceProviderManager, m.MarketManager, c.ConfigManager))

	log.Fatal(app.Listen("0.0.0.0" + ":" + port))
}

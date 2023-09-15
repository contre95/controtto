package rest

import (
	"controtto/src/app/managing"
	"controtto/src/app/querying"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func Run(port string, m *managing.Service, q *querying.Service) {
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
	app.Get("/pairs/", pairsSection)
	app.Get("/pairs/:id/", pairSection)
	app.Get("/ui/assets/form", newAssetForm)
	app.Get("/ui/pairs/form", newPairForm(q.AssetQuerier))
	app.Get("/healthcheck/price", checkPrice(q.MarketQuerier))
	app.Get("/ui/pairs/table", pairsTable(q.TradingPairQuerier))
	app.Get("/dashboard", dashboardSection(q.TradingPairQuerier))
	app.Get("/ui/pairs/:id/cards", pairCards(q.TradingPairQuerier))
	app.Get("/pairs/AvgBuyPrice/:id", avgBuyPrice(q.TradingPairQuerier))
	app.Get("/pairs/:id/transactions/export", transactionExport(q.TradingPairQuerier))
	app.Get("/ui/pairs/:id/transactions/table", transactionTable(q.TradingPairQuerier))
	app.Get("/ui/pairs/:id/transactions/form", newTransactionForm(q.TradingPairQuerier))
	// DELETE
	app.Delete("/empty", empty())
	app.Delete("/pairs/:id", deleteTradingPair(m.TradingPairManager))
	app.Delete("/transactions/:id", deleteTransaction(m.TradingPairManager))
	// POST
	app.Post("/assets", newAsset(m.AssetCreator))
	app.Post("/pairs", newTradingPair(m.TradingPairManager))
	app.Post("/pairs/:id/transactions", newTransaction(m.TradingPairManager))
	app.Post("/pairs/:id/transactions/upload", newTransactionImport(m.TradingPairManager))

	log.Fatal(app.Listen("0.0.0.0" + ":" + port))
}

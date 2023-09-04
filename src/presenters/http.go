package presenters

import (
	"controtto/src/app/managing"
	"controtto/src/app/querying"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func Run(m *managing.Service, q *querying.Service) {
	engine := html.New("./views", ".html")
	engine.Debug(true)
	app := fiber.New(fiber.Config{
		Views: engine,
		// ViewsLayout: "layouts/main",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		},
	})

	// UI
	app.Get("/ui/pairs/form", newPairForm(q.AssetQuerier))
	app.Get("/ui/pairs/table", pairsTable(q.TradingPairQuerier))
	app.Get("/ui/pairs/:id/transactions/table", transactionTable(q.TradingPairQuerier))
	app.Get("/ui/pairs/:id/cards", pairCards(q.MarketQuerier, q.TradingPairQuerier))
	app.Get("/ui/pairs/:id/transactions/form", newTransactionForm(q.TradingPairQuerier))
	app.Get("/ui/assets/form", newAssetForm)
	app.Get("/pairs/:id/", pairSection)
	app.Get("/pairs/", pairsSection)
	app.Get("/dashboard", dashboardSection)
	// GET
	app.Get("/", Home)
	app.Get("/tables", Tables)
	app.Get("/pairs/AvgButPrice/:id", avgBuyPrice(q.TradingPairQuerier))
	app.Get("/users", TotalUsers)
	// DELETE
	app.Delete("/pairs/:id", deleteTradingPair(m.TradingPairManager))
	app.Delete("/transactions/:id", deleteTransaction(m.TradingPairManager))
	// POST
	app.Post("/pairs/:id/transactions", newTransaction(m.TradingPairManager))
	app.Post("/assets", newAsset(m.AssetCreator))
	app.Post("/pairs", newTradingPair(m.TradingPairManager))
	app.Delete("/empty", func(c *fiber.Ctx) error {
		log.Println("hi")
		return c.SendString("")
	})

	app.Static("/assets", "./public/assets")
	log.Fatal(app.Listen("0.0.0.0:8721"))
}

package rest

import (
	"controtto/src/app/querying"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

func pairsSection(c *fiber.Ctx) error {
	slog.Info("Pairs Section")
	return c.Render("pairsSection", fiber.Map{})
}

func pairSection(c *fiber.Ctx) error {
	slog.Info("Pair Section")
	return c.Render("pairSection", fiber.Map{
		"PairID": c.Params("id"),
	})
}

func newPairForm(aq querying.AssetsQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		req := querying.QueryAssetsReq{}
		resp, err := aq.ListAssets(req)
		if err != nil {
			return c.Render("toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   err,
			})
		}
		return c.Render("newPairForm", fiber.Map{
			"Title":  "New Pair",
			"Assets": resp.Assets,
		})
	}
}

// Tables handler that renderizer the tables view and returns it to the client
func pairsTable(aq querying.TradingPairsQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		slog.Info("Pairs table requested")
		req := querying.ListTradingPairsReq{}
		resp, err := aq.ListTradingPairs(req)
		if err != nil {
			return c.Render("toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   err,
			})
		}
		return c.Render("pairsTable", fiber.Map{
			"Title": "Trading Pairs",
			"Pairs": resp.Pairs,
		})
	}
}

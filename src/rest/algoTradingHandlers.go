package rest

import (
	"controtto/src/app/managing"
	"controtto/src/app/querying"
	"controtto/src/domain/pnl"

	"github.com/gofiber/fiber/v2"
)

func getSpatialArbitrageForm(tpq querying.PairsQuerier, marketManager *managing.MarketManager) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Get available markets
		markets := marketManager.GetMarkets(true)
		marketOptions := make([]pnl.Market, 0, len(markets))
		for _, market := range markets {
			marketOptions = append(marketOptions, *market)
		}
		// Get available trading pairs
		pairsResp, err := tpq.ListPairs(querying.ListPairsReq{})
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to load trading pairs")
		}
		return c.Render("spatialArbitrageForm", fiber.Map{
			"Markets": marketOptions,
			"Pairs":   pairsResp.Pairs,
		})
	}
}

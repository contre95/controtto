package rest

import (
	"controtto/src/app/querying"
	"controtto/src/domain/pnl"
	"slices"

	"github.com/gofiber/fiber/v2"
)

func pairTape(tpq querying.PairsQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		req := querying.GetPairReq{
			TPID:             id,
			WithCalculations: false,
		}
		resp, err := tpq.GetPair(req)
		if err != nil {
			return c.Render("toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   err,
			})
		}
		if resp.Pair.BaseAsset.Type == pnl.Stock {
		}
		return c.Render("pairTape", fiber.Map{
			"Pair": resp.Pair,
		})
	}
}

func pairChart(tpq querying.PairsQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		req := querying.GetPairReq{
			TPID:             id,
			WithCalculations: false,
		}
		resp, err := tpq.GetPair(req)
		if err != nil {
			return c.Render("toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   err,
			})
		}
		slices.Reverse(resp.Pair.Trades)
		return c.Render("pairChart", fiber.Map{
			"Pair": resp.Pair,
		})
	}
}



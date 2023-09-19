package rest

import (
	"controtto/src/app/querying"
	"controtto/src/domain/pnl"
	"log/slog"
	"slices"
	"time"

	"github.com/gofiber/fiber/v2"
)

func empty() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.SendString("")
	}
}

func pairCards(tpq querying.TradingPairsQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		req := querying.GetTradingPairReq{
			TPID:             id,
			WithCalculations: true,
		}
		resp, err := tpq.GetTradingPair(req)
		if err != nil {
			return c.Render("toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   err,
			})
		}
		slices.Reverse(resp.Pair.Transactions)
		return c.Render("pairCards", fiber.Map{
			"Today": time.Now().Format("Mon Jan 02 15:04 2006"),
			"Pair":  resp.Pair,
		})
	}
}

func dashboardSection(aq querying.TradingPairsQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		slog.Info("Dashboard requested")
		req := querying.ListTradingPairsReq{}
		resp, err := aq.ListTradingPairs(req)
		if err != nil {
			return c.Render("toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   err,
			})
		}
		return c.Render("dashboardSection", fiber.Map{
			"Title": "Trading Pairs",
			"Pairs": resp.Pairs,
		})
	}

}

// Home hanlder reders the homescreen
func Home(c *fiber.Ctx) error {
	slog.Info("HOME")
	// render index template
	return c.Render("main", fiber.Map{})
}

func pairChart(tpq querying.TradingPairsQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		req := querying.GetTradingPairReq{
			TPID:             id,
			WithCalculations: true,
		}
		resp, err := tpq.GetTradingPair(req)
		if err != nil {
			return c.Render("toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   err,
			})
		}
		slices.Reverse(resp.Pair.Transactions)
		var tvAssetType string
		switch resp.Pair.BaseAsset.Type {
		case pnl.Stock:
			tvAssetType = "NASDAQ:"
		case pnl.Forex:
			tvAssetType = "FX:"
		case pnl.Crypto:
			tvAssetType = "COINBASE:"
		}
		if resp.Pair.BaseAsset.Type == pnl.Stock {
		}
		return c.Render("pairChart", fiber.Map{
			"Pair":                 resp.Pair,
			"TradingViewAssetType": tvAssetType,
		})
	}
}

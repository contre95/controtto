package rest

import (
	"controtto/src/app/managing"
	"controtto/src/app/querying"
	"log/slog"
	"slices"
	"time"

	"github.com/gofiber/fiber/v2"
)

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
		slices.Reverse(resp.Pair.Trades)
		return c.Render("pairCards", fiber.Map{
			"Today": time.Now().Format("Mon Jan 02 15:04 2006"),
			"Pair":  resp.Pair,
		})
	}
}

func newPairForm(aq querying.AssetsQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		req := querying.ListAssetsReq{}
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
func tradesTable(aq querying.TradingPairsQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		slog.Info("Trades table requested")
		req := querying.ListTradingPairsReq{}
		resp, err := aq.ListTradingPairs(req)
		if err != nil {
			return c.Render("toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   err,
			})
		}
		return c.Render("tradesTable", fiber.Map{
			"Title": "Trading Pairs",
			"Pairs": resp.Pairs,
		})
	}
}

func deleteTradingPair(tpm managing.TradingPairsManager) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		req := managing.DeleteTradingPairReq{
			ID: c.Params("id"),
		}
		slog.Info("Delete", "id", req.ID)
		resp, err := tpm.DeleteTradingPair(req)
		if err != nil {
			return c.Render("toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   err,
			})
		}
		slog.Info("Trading pair deleted", "trading-pair", resp.ID)
		c.Append("HX-Trigger", "newPair")
		return c.Render("toastOk", fiber.Map{
			"Title": "Deleted",
			"Msg":   "Pair deleted",
		})

	}
}

func newTradingPair(tpc managing.TradingPairsManager) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		slog.Info("Creating new pair")
		payload := struct {
			Quote string `form:"quote"`
			Base  string `form:"base"`
		}{}
		if err := c.BodyParser(&payload); err != nil {
			return err
		}
		req := managing.CreateTradingPairReq{
			BaseAssetSymbol:  payload.Base,
			QuoteAssetSymbol: payload.Quote,
		}
		slog.Info("Creating", "req", req)
		resp, err := tpc.Create(req)
		if err != nil {
			return c.Render("toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   err,
			})
		}
		c.Append("HX-Trigger", "newPair")
		return c.Render("toastOk", fiber.Map{
			"Title": "Created",
			"Msg":   resp.Msg,
		})
	}
}

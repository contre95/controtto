package rest

import (
	"controtto/src/app/config"
	"controtto/src/app/querying"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

func getMarketAssets(cfg *config.Config, tpq querying.TradingPairsQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		slog.Info("Create Trade UI requested")
		id := c.Params("id")
		req := querying.GetTradingPairReq{
			TPID:                 id,
			WithCurrentBasePrice: true,
			WithTrades:           false,
			WithCalculations:     false,
		}
		resp, err := tpq.GetTradingPair(req)
		if err != nil {
			return c.Render("toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   err,
			})
		}
		marketData := map[string][3]float64{}
		for name, m := range cfg.GetMarketTraders(false) {
			baseAmount, err := m.MarketAPI.FetchAssetAmount(resp.Pair.BaseAsset.Symbol)
			if err != nil {
				return c.Render("toastErr", fiber.Map{
					"Title": "Error",
					"Msg":   err,
				})
			}
			usdtAmount, err := m.MarketAPI.FetchAssetAmount("USDT")
			usdcAmount, err := m.MarketAPI.FetchAssetAmount("USDC")
			if err != nil {
				return c.Render("toastErr", fiber.Map{
					"Title": "Error",
					"Msg":   err,
				})
			}
			marketData[name] = [3]float64{baseAmount, usdtAmount, usdcAmount}
		}
		return c.Render("marketAssets", fiber.Map{
			"Pair":          resp.Pair,
			"Today":         time.Now().Format("2006-01-02"),
			"MarketTraders": cfg.GetMarketTraders(false),
			"MarketData":    marketData,
		})
	}
}

func newMarketTradingForm(cfg *config.Config, tpq querying.TradingPairsQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		slog.Info("Create Trade UI requested")
		id := c.Params("id")
		req := querying.GetTradingPairReq{
			TPID:                 id,
			WithCurrentBasePrice: true,
			WithTrades:           false,
			WithCalculations:     false,
		}
		resp, err := tpq.GetTradingPair(req)
		if err != nil {
			return c.Render("toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   err,
			})
		}
		return c.Render("marketTradingForm", fiber.Map{
			"Pair":          resp.Pair,
			"Today":         time.Now().Format("2006-01-02"),
			"MarketTraders": cfg.GetMarketTraders(true),
		})
	}
}

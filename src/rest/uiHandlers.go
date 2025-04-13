package rest

import (
	"controtto/src/app/config"
	"controtto/src/app/querying"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

func empty() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.SendString("")
	}
}

// Home hanlder reders the homescreen
func Home(c *fiber.Ctx) error {
	slog.Info("HOME")
	return c.Render("main", fiber.Map{"TradesTrigger": ",revealed"})
}

func dashboardSection(aq querying.TradingPairsQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if c.Get("HX-Request") != "true" {
			return c.Render("main", fiber.Map{
				"DashboardTrigger": ",revealed",
			})
		}
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

func tradesSection(c *fiber.Ctx) error {
	slog.Info("Trades Section")
	if c.Get("HX-Request") != "true" {
		return c.Render("main", fiber.Map{
			"TradesTrigger": ",revealed",
		})
	}
	return c.Render("tradesSection", fiber.Map{
		"Amount": 4,
	})
}

func pairSection(c *fiber.Ctx) error {
	slog.Info("Pair Section")
	return c.Render("pairSection", fiber.Map{
		"PairID": c.Params("id"),
	})
}

func settingsSection(cfg *config.Config) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if c.Get("HX-Request") != "true" {
			return c.Render("main", fiber.Map{
				"SettingsTrigger": ",revealed",
			})
		}
		return c.Render("settingsSection", fiber.Map{
			"Today":          time.Now().Format("Mon Jan 02 15:04 2006"),
			"Uncommon":       cfg.GetUncommonPairs(),
			"PriceProviders": cfg.GetPriceProviders(), // Pass the providers slice to the template
		})
	}
}

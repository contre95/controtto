package rest

import (
	"controtto/src/app/config"
	"controtto/src/app/managing"
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

func dashboardSection(aq querying.PairsQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if c.Get("HX-Request") != "true" {
			return c.Render("main", fiber.Map{
				"DashboardTrigger": ",revealed",
			})
		}
		slog.Info("Dashboard requested")
		req := querying.ListPairsReq{}
		resp, err := aq.ListPairs(req)
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

func pairSection() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		slog.Info("Pair Section")
		if c.Get("HX-Request") != "true" {
			return c.Render("main", fiber.Map{
				"PairID":       c.Params("id"),
				"TradeTrigger": ",revealed",
			})
		}
		return c.Render("pairSection", fiber.Map{
			"PairID": c.Params("id"),
		})
	}
}

func settingsSection(priceProviderManager *managing.PriceProviderManager, marketManager *managing.MarketManager, cfg *config.ConfigManager) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if c.Get("HX-Request") != "true" {
			return c.Render("main", fiber.Map{
				"SettingsTrigger": ",revealed",
			})
		}

		// Get configured price providers and market traders
		providers := priceProviderManager.ListProviders(true) // true = show all providers
		traders := marketManager.ListTraders(true)            // true = show all traders

		return c.Render("settingsSection", fiber.Map{
			"Today":          time.Now().Format("Mon Jan 02 15:04 2006"),
			"Uncommon":       cfg.GetUncommonPairs(),
			"PriceProviders": providers,
			"MarketTraders":  traders,
		})
	}
}

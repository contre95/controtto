package rest

import (
	"controtto/src/app/config"
	"controtto/src/app/managing"
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

//	func marketsSet(marketManager *managing.MarketManager) bool {
//		traders := marketManager.GetMarketTraders(false) // false = only configured traders
//		return len(traders) > 0
//	}
//
//	func marketsSetAPI(marketManager *managing.MarketManager) func(*fiber.Ctx) error {
//		return func(c *fiber.Ctx) error {
//			return c.SendString(fmt.Sprintf("%t", marketsSet(marketManager)))
//		}
//	}
func saveSettingsForm(priceProviderManager *managing.PriceProviderManager, marketManager *managing.MarketManager, cfg *config.ConfigManager) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Update price providers
		providers := priceProviderManager.ListProviders(true) // true = get all providers
		var enable bool
		var token string
		for key, p := range providers {
			if !p.NeedsToken {
				token := c.FormValue(p.ProviderKey)
				enable = token != ""
			} else {
				enable = c.FormValue(p.ProviderKey) != ""
			}
			err := priceProviderManager.UpdateProvider(key, token, enable)
			if err != nil {
				slog.Error("Error updating price provider", "provider", key, "error", err)
				return c.Render("toastErr", fiber.Map{
					"Title": "Error",
					"Msg":   fmt.Sprintf("Failed to update %s: %v", p.ProviderName, err),
				})
			}
		}

		// Update market traders
		traders := marketManager.GetMarkets(true) // true = get all traders
		for key, t := range traders {
			token := c.FormValue(t.MarketKey)
			err := marketManager.UpdateMarket(key, token)
			if err != nil {
				slog.Error("Error updating market trader", "trader", key, "error", err)
				return c.Render("toastErr", fiber.Map{
					"Title": "Error",
					"Msg":   fmt.Sprintf("Failed to update %s: %v", t.MarketName, err),
				})
			}
		}

		// Update uncommon pairs setting
		uncommon := c.FormValue("uncommon_pairs") != ""
		cfg.SetUncommonPairs(uncommon)

		c.Append("HX-Trigger", "reloadSettings")
		slog.Info("Config updated successfully")
		return c.Render("toastOk", fiber.Map{
			"Title": "Success",
			"Msg":   "Configuration updated successfully",
		})
	}
}

package rest

import (
	"controtto/src/app/config"
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

func marketsSet(cfg *config.ConfigManager) bool {
	for _, mkt := range cfg.GetMarketTraders(true) {
		if mkt.IsSet {
			return true
		}
	}
	return false
}

func marketsSetAPI(cfg *config.ConfigManager) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.SendString(fmt.Sprintf("%t", marketsSet(cfg)))
	}
}

func saveSettingsForm(cfg *config.ConfigManager) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		for key, p := range cfg.GetPriceProviders() {
			token := c.FormValue(p.ProviderKey)
			err := cfg.UpdatePriceProvider(key, token, len(token) != 0)
			if err != nil {
				return c.Render("toastErr", fiber.Map{
					"Title": "Error",
					"Msg":   err,
				})
			}
		}
		for key, t := range cfg.GetMarketTraders(true) {
			token := c.FormValue(t.MarketKey)
			err := cfg.UpdateMarketTraderToken(key, token)
			if err != nil {
				return c.Render("toastErr", fiber.Map{
					"Title": "Error",
					"Msg":   err,
				})
			}
		}
		uncommon := c.FormValue("uncommon_pairs") != ""
		cfg.SetUncommonPairs(uncommon)
		c.Append("HX-Trigger", "reloadSettings")
		slog.Info("Config updated")
		return c.Render("toastOk", fiber.Map{
			"Title": "Created",
			"Msg":   "Config updated",
		})
	}
}

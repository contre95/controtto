package rest

import (
	"controtto/src/app/config"
	"fmt"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

func editSettingsForm(cfg *config.Config) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.Render("editSettingsForm", fiber.Map{
			"Today":    time.Now().Format("Mon Jan 02 15:04 2006"),
			"Port":     cfg.Port,
			"Uncommon": cfg.UncommonPairs,
		})
	}
}

func marketsSet(cfg *config.Config) bool {
	for _, mkt := range cfg.GetMarketTraders(true) {
		if mkt.IsSet {
			return true
		}
	}
	return false
}

func marketsSetAPI(cfg *config.Config) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.SendString(fmt.Sprintf("%t", marketsSet(cfg)))
	}
}

func saveSettingsForm(cfg *config.Config) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		for key, p := range cfg.GetPriceProviders() {
			token := c.FormValue(p.ProviderKey)
			cfg.UpdatePriceProvider(key, token, len(token) > 0)
		}
		for key, t := range cfg.GetMarketTraders(true) {
			token := c.FormValue(t.MarketKey)
			cfg.UpdateMarketTraderToken(key, token)
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

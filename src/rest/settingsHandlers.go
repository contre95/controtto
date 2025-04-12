package rest

import (
	"controtto/src/app/config"
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

func saveSettingsForm(cfg *config.Config) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		for key, p := range cfg.GetPriceProviders() {
			token := c.FormValue(p.ProviderInputName)
			cfg.UpdateProviderToken(key, token)
		}
		uncommon := c.FormValue("uncommon_pairs") != ""
		cfg.SetUncommonPairs(uncommon)
		slog.Info("Config updated")
		c.Append("HX-Trigger", "reloadSettings")
		return c.Render("toastOk", fiber.Map{
			"Title": "Created",
			"Msg":   "Config upated",
		})

	}
}

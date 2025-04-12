package rest

import (
	"controtto/src/app/config"
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
		// // If the form checkbox for "Uncommon" is checked, the value might be "on"
		for key := range cfg.GetPriceProviders() {
			token := c.FormValue(key)
			cfg.UpdateProviderToken(key, token)
		}
		return c.Redirect("/settings")
	}
}

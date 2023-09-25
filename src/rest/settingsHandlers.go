package rest

import (
	"controtto/src/app/config"
	"time"

	"github.com/gofiber/fiber/v2"
)

func settingsSection(conf *config.Service) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		cfg := conf.IsSet()
		return c.Render("settingsSection", fiber.Map{
			"Today":             time.Now().Format("Mon Jan 02 15:04 2006"),
			"AlphaVantageToken": cfg.AVantageAPIToken,
			"TiingoToken":       cfg.TiingoAPIToken,
			"Uncommon":          cfg.UncommonPairs,
		})
	}
}

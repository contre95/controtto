package rest

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func settingsSection() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.Render("settingsSection", fiber.Map{
			"Today":             time.Now().Format("Mon Jan 02 15:04 2006"),
			"AlphaVantageToken": "",
			"TiingoToken":       "asljfsd",
			"Uncommon":          true,
		})
	}
}

package presenters

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

func newAssetForm(c *fiber.Ctx) error {
	slog.Info("Create Asset UI requested")
	return c.Render("newAssetForm", fiber.Map{
		"Title": "New Asset",
	})
}

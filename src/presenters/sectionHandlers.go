package presenters

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

func dashboardSection(c *fiber.Ctx) error {
	slog.Info("DASHBOARD")
	return c.Render("dashboardSection", fiber.Map{})
}

func pairsSection(c *fiber.Ctx) error {
	slog.Info("Pairs Section")
	return c.Render("pairsSection", fiber.Map{})
}

func pairSection(c *fiber.Ctx) error {
	slog.Info("Pair Section")
	return c.Render("pairSection", fiber.Map{
		"PairID": c.Params("id"),
	})
}

// Home hanlder reders the homescreen
func Home(c *fiber.Ctx) error {
	slog.Info("HOME")
	// render index template
	return c.Render("main", fiber.Map{})
}

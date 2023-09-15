package rest

import (
	"controtto/src/app/querying"

	"github.com/gofiber/fiber/v2"
)

func avgBuyPrice(tpq querying.TradingPairsQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

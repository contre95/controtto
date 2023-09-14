package presenters

import (
	"controtto/src/app/querying"
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

func marketPrice(aq querying.MarketsQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		slog.Info("Requesting market price", "market", id)
		req := querying.QueryMarketReq{
			AssetSymbolA: id,
		}
		resp, err := aq.GetMarketPrice(req)
		if err != nil {
			return c.SendString("--")
		}
		// render index template
		return c.SendString(fmt.Sprintf("%.2f", resp.Price))
	}
}

func checkPrice(mq querying.MarketsQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		base := c.Query("base")
		quote := c.Query("quote")
		req := querying.QueryMarketReq{
			AssetSymbolA: base,
			AssetSymbolB: quote,
		}
		resp, err := mq.GetMarketPrice(req)
		if err != nil {
			return c.SendString(err.Error())
		}
		return c.SendString(fmt.Sprintf("%f", resp.Price))
	}
}

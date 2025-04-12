package rest

import (
	"controtto/src/app/querying"
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

func marketPrice(aq querying.PriceQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		slog.Info("Requesting market price", "market", id)
		req := querying.QueryPriceReq{
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

func checkPrice(mq querying.PriceQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		base := c.Query("base")
		quote := c.Query("quote")
		req := querying.QueryPriceReq{
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

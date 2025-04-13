package rest

import (
	"controtto/src/app/querying"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func checkPrice(mq querying.PriceQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		base := c.Query("base")
		quote := c.Query("quote")
		req := querying.QueryPriceReq{
			AssetSymbolA: base,
			AssetSymbolB: quote,
		}
		resp, err := mq.GetPrice(req)
		if err != nil {
			return c.SendString(err.Error())
		}
		return c.SendString(fmt.Sprintf("%f", resp.Price))
	}
}

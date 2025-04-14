package rest

import (
	"controtto/src/app/querying"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func avgBuyPrice(tpq querying.TradingPairsQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return nil
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
		resp, err := mq.GetPrice(req)
		if err != nil {
			return c.SendString(err.Error())
		}
		return c.SendString(fmt.Sprintf("%f", resp.Price))
	}
}
func calculatePrice(mq querying.PriceQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		base := c.Query("base")
		quote := c.Query("quote")
		amountStr := c.Query("amount")

		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"code":     1,
				"msg":      "Invalid amount",
				"debugMsg": err.Error(),
				"data":     nil,
			})
		}

		req := querying.QueryPriceReq{
			AssetSymbolA: base,
			AssetSymbolB: quote,
		}
		resp, err := mq.GetPrice(req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":     2,
				"msg":      "Failed to fetch price",
				"debugMsg": err.Error(),
				"data":     nil,
			})
		}

		totalPrice := resp.Price * amount
		return c.SendString(fmt.Sprintf("%.2f", totalPrice))

		// return c.JSON(fiber.Map{
		// 	"code":     0,
		// 	"msg":      "Success",
		// 	"debugMsg": "",
		// 	"data": fiber.Map{
		// 		"base":       base,
		// 		"quote":      quote,
		// 		"amount":     amount,
		// 		"price":      resp.Price,
		// 		"totalPrice": totalPrice,
		// 	},
		// })
	}
}

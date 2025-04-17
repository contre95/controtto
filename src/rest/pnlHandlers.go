package rest

import (
	"controtto/src/app/managing"
	"controtto/src/app/querying"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func avgBuyPrice(tpq querying.TradingPairsQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Implementation would go here
		return nil
	}
}

func checkPrice(priceProviderManager *managing.PriceProviderManager) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		base := c.Query("base")
		quote := c.Query("quote")
		fmt.Println(quote)
		fmt.Println(base)

		if base == "" || quote == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Both base and quote parameters are required")
		}

		// Try all configured providers until we get a price
		var price float64
		var lastError error
		providers := priceProviderManager.ListProviders(false) // false = only configured providers

		for providerKey := range providers {
			req := managing.QueryPriceReq{
				AssetSymbolA: base,
				AssetSymbolB: quote,
			}
			resp, err := priceProviderManager.QueryPrice(req)
			if err == nil {
				price = resp.Price
				break
			}
			lastError = err
			slog.Warn("Failed to get price from provider",
				"provider", providerKey,
				"error", err,
				"pair", fmt.Sprintf("%s/%s", base, quote),
			)
		}

		if price == 0 && lastError != nil {
			return c.Status(fiber.StatusServiceUnavailable).SendString(lastError.Error())
		}

		return c.SendString(fmt.Sprintf("%f", price))
	}
}

func calculatePrice(priceProviderManager *managing.PriceProviderManager) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		base := c.Query("base")
		quote := c.Query("quote")
		amountStr := c.Query("amount")
		fmt.Println(quote, base, amountStr)
		if base == "" || quote == "" {
			return c.SendString("err")
		}

		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			return c.SendString("err")
		}
		// Try all configured providers until we get a price
		var price float64
		var lastError error
		providers := priceProviderManager.ListProviders(false) // false = only configured providers

		for providerKey := range providers {
			req := managing.QueryPriceReq{
				AssetSymbolA: base,
				AssetSymbolB: quote,
			}
			resp, err := priceProviderManager.QueryPrice(req)
			if err == nil {
				price = resp.Price
				break
			}
			lastError = err
			slog.Warn("Failed to get price from provider",
				"provider", providerKey,
				"error", err,
				"pair", fmt.Sprintf("%s/%s", base, quote),
			)
		}

		if price == 0 {
			if lastError != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"code":     3,
					"msg":      "Failed to fetch price from any provider",
					"debugMsg": lastError.Error(),
					"data":     nil,
				})
			}
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"code":     4,
				"msg":      "No price providers available",
				"debugMsg": "",
				"data":     nil,
			})
		}

		totalPrice := price * amount
		// return c.JSON(fiber.Map{
		// 	"code":     0,
		// 	"msg":      "Success",
		// 	"debugMsg": "",
		// 	"data": fiber.Map{
		// 		"base":       base,
		// 		"quote":      quote,
		// 		"amount":     amount,
		// 		"price":      price,
		// 		"totalPrice": totalPrice,
		// 	},
		// })
		return c.SendString(fmt.Sprintf("%2.f", totalPrice))
	}
}

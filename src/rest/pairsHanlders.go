package rest

import (
	"controtto/src/app/managing"
	"controtto/src/app/querying"
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func pairCards(tpq querying.PairsQuerier, pq *managing.PriceProviderManager) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		priceStr := c.Query("price")
		var errMsg string
		var price float64
		name := "ERROR"
		color := "#E0663D"
		id := c.Params("id")
		req := querying.GetPairReq{TPID: id}
		resp, err := tpq.GetPair(req)
		if err != nil {
			return c.Render("pairCards", fiber.Map{
				"Error": err.Error(),
			})
		}
		if priceStr != "" {
			if price, err = strconv.ParseFloat(priceStr, 64); err == nil {
				name = "Custom"
				color = "#67697C"
			} else {
				name = "BadPrice"
				color = "#E0668D"
			}
		} else if priceResp, err := pq.QueryPrice(managing.QueryPriceReq{
			AssetSymbolA: resp.Pair.BaseAsset.Symbol,
			AssetSymbolB: resp.Pair.QuoteAsset.Symbol,
		}); err == nil {
			price = priceResp.Price
			name = priceResp.ProviderName
			color = priceResp.ProviderColor
		}
		fmt.Println("Price", price)
		req.BasePrice = price
		req.WithCalculations = true
		resp, err = tpq.GetPair(req)
		if err != nil {
			return c.Render("pairCards", fiber.Map{
				"Error": err.Error(),
			})
		}
		slices.Reverse(resp.Pair.Trades)
		return c.Render("pairCards", fiber.Map{
			"Today":              time.Now().Format("Mon Jan 02 15:04 2006"),
			"Pair":               resp.Pair,
			"Price":              price,
			"PriceProviderName":  name,
			"PriceProviderColor": color,
			"Error":              errMsg != "",
			"ErrMsg":             errMsg,
		})
	}
}

func newPairForm(aq querying.AssetsQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		req := querying.ListAssetsReq{}
		resp, err := aq.ListAssets(req)
		if err != nil {
			return c.Render("toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   err,
			})
		}
		return c.Render("newPairForm", fiber.Map{
			"Title":  "New Pair",
			"Assets": resp.Assets,
		})
	}
}
func getAssetColor(aq querying.AssetsQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		symbol := c.Params("symbol")
		req := querying.QueryAssetReq{
			Symbol: symbol,
		}
		resp, err := aq.GetAsset(req)
		if err != nil {
			return c.SendString("")
		}
		return c.SendString(resp.Asset.Color)
	}
}

// Tables handler that renderizer the tables view and returns it to the client
func tradesTable(aq querying.PairsQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		slog.Info("Trades table requested")
		req := querying.ListPairsReq{}
		resp, err := aq.ListPairs(req)
		if err != nil {
			return c.Render("toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   err,
			})
		}
		return c.Render("tradesTable", fiber.Map{
			"Title": "Trading Pairs",
			"Pairs": resp.Pairs,
		})
	}
}

func deletePair(tpm managing.PairsManager) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		req := managing.DeletePairReq{
			ID: c.Params("id"),
		}
		slog.Info("Delete", "id", req.ID)
		resp, err := tpm.DeletePair(req)
		if err != nil {
			return c.Render("toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   err,
			})
		}
		slog.Info("Trading pair deleted", "trading-pair", resp.ID)
		c.Append("HX-Trigger", "newPair")
		return c.Render("toastOk", fiber.Map{
			"Title": "Deleted",
			"Msg":   "Pair deleted",
		})

	}
}

func newPair(tpc managing.PairsManager) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		slog.Info("Creating new pair")
		payload := struct {
			Quote string `form:"quote"`
			Base  string `form:"base"`
		}{}
		if err := c.BodyParser(&payload); err != nil {
			return err
		}
		req := managing.CreatePairReq{
			BaseAssetSymbol:  payload.Base,
			QuoteAssetSymbol: payload.Quote,
		}
		slog.Info("Creating", "req", req)
		resp, err := tpc.Create(req)
		if err != nil {
			return c.Render("toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   err,
			})
		}
		c.Append("HX-Trigger", "newPair")
		return c.Render("toastOk", fiber.Map{
			"Title": "Created",
			"Msg":   resp.Msg,
		})
	}
}

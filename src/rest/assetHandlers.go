package rest

import (
	"controtto/src/app/managing"
	"controtto/src/domain/pnl"
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

func newAssetForm(c *fiber.Ctx) error {
	slog.Info("Create Asset UI requested")
	return c.Render("newAssetForm", fiber.Map{
		"Title":      "New Asset",
		"AssetTypes": pnl.GetValidTypes(),
	})
}

func newAsset(ac managing.AssetsCreator) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		slog.Info("Creating new asset")
		payload := struct {
			Symbol string `form:"symbol"`
			Name   string `form:"name"`
			Type   string `form:"atype"`
			Color  string `form:"color"`
		}{}
		if err := c.BodyParser(&payload); err != nil {
			return err
		}
		var assetType pnl.AssetType
		switch payload.Type {
		case "Crypto":
			assetType = pnl.Crypto
		case "Forex":
			assetType = pnl.Forex
		case "Stock":
			assetType = pnl.Stock
		default:
			return c.Render("toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   fmt.Sprintf("invalid asset type: %s", payload.Type),
			})
		}
		req := managing.CreateAssetReq{
			Symbol:      payload.Symbol,
			Color:       payload.Color,
			Type:        assetType,
			Name:        payload.Name,
			CountryCode: "-",
		}
		resp, err := ac.Create(req)
		if err != nil {
			return c.Render("toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   err,
			})
		}
		return c.Render("toastOk", fiber.Map{
			"Title": "Created",
			"Msg":   resp.Msg,
		})
	}
}

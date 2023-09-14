package presenters

import (
	"controtto/src/app/managing"
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

func newAssetForm(c *fiber.Ctx) error {
	slog.Info("Create Asset UI requested")
	return c.Render("newAssetForm", fiber.Map{
		"Title": "New Asset",
	})
}

func newAsset(ac managing.AssetsCreator) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		slog.Info("Creating new asset")
		payload := struct {
			Symbol string `form:"symbol"`
			Name   string `form:"name"`
			Color  string `form:"color"`
		}{}
		if err := c.BodyParser(&payload); err != nil {
			fmt.Println(err)
			return err
		}
		req := managing.CreateAssetReq{
			Symbol:      payload.Symbol,
			Color:       payload.Color,
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

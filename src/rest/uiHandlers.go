package rest

import (
	"controtto/src/app/config"
	"controtto/src/app/querying"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

func empty() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.SendString("")
	}
}

// Home hanlder reders the homescreen
func Home(c *fiber.Ctx) error {
	slog.Info("HOME")
	// render index template
	return c.Render("main", fiber.Map{})
}

func dashboardSection(aq querying.TradingPairsQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		slog.Info("Dashboard requested")
		req := querying.ListTradingPairsReq{}
		resp, err := aq.ListTradingPairs(req)
		if err != nil {
			return c.Render("toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   err,
			})
		}
		return c.Render("dashboardSection", fiber.Map{
			"Title": "Trading Pairs",
			"Pairs": resp.Pairs,
		})
	}
}

func tradesSection(c *fiber.Ctx) error {
	slog.Info("Trades Section")
	if c.Get("HX-Request") != "true" {
		return c.Render("main", fiber.Map{
			"TradesTrigger": "revealed",
		})
	}
	return c.Render("tradesSection", fiber.Map{
		"Amount": 4,
	})
}

func pairSection(c *fiber.Ctx) error {
	slog.Info("Pair Section")
	return c.Render("pairSection", fiber.Map{
		"PairID": c.Params("id"),
	})
}

func settingsSection(cfg *config.Config) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if c.Get("HX-Request") != "true" {
			return c.Render("main", fiber.Map{
				"SettingsTrigger": "revealed",
			})
		}
		// Build providers slice dynamically based on the cfg.Tokens.
		providers := make([]map[string]interface{}, 0)
		for _, token := range cfg.PriceProviders {
			providers = append(providers, map[string]interface{}{
				"TokenSet":       token.TokenSet,
				"TokenTitle":     token.ProviderName,
				"TokenURL":       token.ProviderURL,
				"TokenInputName": token.ProviderInputName,
				"Token":          token.Token,
			})
		}

		// Render the settingsSection partial with the dynamic providers.
		return c.Render("settingsSection", fiber.Map{
			"Today":          time.Now().Format("Mon Jan 02 15:04 2006"),
			"Uncommon":       cfg.UncommonPairs,
			"PriceProviders": providers, // Pass the providers slice to the template
		})
	}
}

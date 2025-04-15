package rest

import (
	"controtto/src/app/managing"
	"controtto/src/app/querying"
	"fmt"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

func getMarketAssets(tpq querying.TradingPairsQuerier, marketManager *managing.MarketManager) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		slog.Info("Market assets request received")

		// Get trading pair details
		id := c.Params("id")
		req := querying.GetTradingPairReq{
			TPID:             id,
			WithTrades:       false,
			WithCalculations: false,
		}

		resp, err := tpq.GetTradingPair(req)
		if err != nil {
			slog.Error("Failed to get trading pair",
				"error", err,
				"pairID", id,
			)
			return c.Render("toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   fmt.Sprintf("Failed to load trading pair: %v", err),
			})
		}

		// Get configured market traders
		marketTraders := marketManager.ListTraders(false) // false = only configured traders

		// Prepare market data (mock data - replace with actual implementation)
		marketData := make(map[string]struct {
			HasError     bool
			ErrorMessage string
			Amounts      [3]float64 // [base, USDT, USDC]
		})

		for marketName := range marketTraders {
			// TODO: Replace with actual market data fetching logic
			marketData[marketName] = struct {
				HasError     bool
				ErrorMessage string
				Amounts      [3]float64
			}{
				HasError:     false,
				ErrorMessage: "",
				Amounts:      [3]float64{1.2345, 100.00, 50.00}, // Example values
			}
		}

		return c.Render("marketAssets", fiber.Map{
			"Pair":          resp.Pair,
			"MarketTraders": marketTraders,
			"MarketData":    marketData,
		})
	}
}

func newMarketTradingForm(tpq querying.TradingPairsQuerier, marketManager *managing.MarketManager) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		slog.Info("Create Market Trading UI requested")

		// Get trading pair details
		id := c.Params("id")
		req := querying.GetTradingPairReq{
			TPID:             id,
			WithTrades:       false,
			WithCalculations: false,
		}

		resp, err := tpq.GetTradingPair(req)
		if err != nil {
			slog.Error("Failed to get trading pair",
				"error", err,
				"pairID", id,
			)
			return c.Render("toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   fmt.Sprintf("Failed to load trading pair: %v", err),
			})
		}

		// Get all market traders (including unconfigured ones)
		marketTraders := marketManager.ListTraders(true) // true = show all traders

		return c.Render("marketTradingForm", fiber.Map{
			"Pair":          resp.Pair,
			"Today":         time.Now().Format("2006-01-02"),
			"MarketTraders": marketTraders,
		})
	}
}

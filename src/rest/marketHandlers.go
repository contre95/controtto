package rest

import (
	"controtto/src/app/managing"
	"controtto/src/app/querying"
	"controtto/src/domain/pnl"
	"fmt"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

func getMarketAssets(tpq querying.PairsQuerier, marketManager *managing.MarketManager) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		slog.Info("Market assets request received")

		// Get trading pair details
		id := c.Params("id")
		req := querying.GetPairReq{
			TPID:             id,
			WithTrades:       false,
			WithCalculations: false,
		}
		resp, err := tpq.GetPair(req)
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
			Amounts      map[string]float64 // [base, USDT, USDC]
		})
		for marketName, trader := range marketTraders {
			var assetAmounts = make(map[string]float64, len(trader.MarketTradingSymbols)+1)
			var err2 error
			var errMsg string
			baseAmt, err1 := marketManager.FetchBalance(marketName, resp.Pair.BaseAsset.Symbol)
			if err1 != nil {
				errMsg += "Base: " + err1.Error() + ". "
			}
			assetAmounts[resp.Pair.BaseAsset.Symbol] = baseAmt
			if trader.Type != pnl.Wallet {
				for _, symbol := range trader.MarketTradingSymbols {
					amount, err2 := marketManager.FetchBalance(marketName, symbol)
					if err2 != nil {
						errMsg += symbol + err2.Error() + ". "
					}
					assetAmounts[symbol] = amount
				}
			}
			hasErr := err1 != nil || err2 != nil
			if hasErr {
				fmt.Println("Error fetching market data:", errMsg)
			}
			marketData[marketName] = struct {
				HasError     bool
				ErrorMessage string
				Amounts      map[string]float64
			}{
				HasError:     hasErr,
				ErrorMessage: errMsg,
				Amounts:      assetAmounts,
			}
		}
		return c.Render("marketAssets", fiber.Map{
			"Pair":          resp.Pair,
			"MarketTraders": marketTraders,
			"MarketData":    marketData,
		})
	}
}

func newMarketTradingForm(tpq querying.PairsQuerier, marketManager *managing.MarketManager) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		slog.Info("Create Market Trading UI requested")

		// Get trading pair details
		id := c.Params("id")
		req := querying.GetPairReq{
			TPID:             id,
			WithTrades:       false,
			WithCalculations: false,
		}

		resp, err := tpq.GetPair(req)
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

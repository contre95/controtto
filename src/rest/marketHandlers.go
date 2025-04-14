package rest

import (
	"controtto/src/app/config"
	"controtto/src/app/querying"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

func getMarketAssets(cfg *config.ConfigManager, tpq querying.TradingPairsQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		traders := cfg.GetMarketTraders(false)
		slog.Info("Create Trade UI requested")
		id := c.Params("id")
		req := querying.GetTradingPairReq{
			TPID:                 id,
			WithCurrentBasePrice: true,
			WithTrades:           false,
			WithCalculations:     false,
		}
		resp, err := tpq.GetTradingPair(req)
		if err != nil {
			return c.Render("toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   err,
			})
		}

		// Changed to a struct that includes both amounts and error status
		type MarketInfo struct {
			Amounts      [3]float64
			HasError     bool
			ErrorMessage string
		}

		marketData := make(map[string]MarketInfo)
		baseSymbol := resp.Pair.BaseAsset.Symbol

		for name, m := range traders {
			var amounts [3]float64
			var hasError bool
			var errMsg string

			baseAmount, err := m.MarketAPI.FetchAssetAmount(baseSymbol)
			if err != nil {
				hasError = true
				errMsg += err.Error() + " "
				slog.Warn("failed to fetch base asset", "market", name, "error", err)
			} else {
				amounts[0] = baseAmount
			}

			usdtAmount, err := m.MarketAPI.FetchAssetAmount("USDT")
			if err != nil {
				hasError = true
				errMsg += err.Error() + " "
				slog.Warn("failed to fetch USDT", "market", name, "error", err)
			} else {
				amounts[1] = usdtAmount
			}

			usdcAmount, err := m.MarketAPI.FetchAssetAmount("USDC")
			if err != nil {
				hasError = true
				errMsg += err.Error() + " "
				slog.Warn("failed to fetch USDC", "market", name, "error", err)
			} else {
				amounts[2] = usdcAmount
			}

			marketData[name] = MarketInfo{
				Amounts:      amounts,
				HasError:     hasError,
				ErrorMessage: errMsg,
			}
		}

		return c.Render("marketAssets", fiber.Map{
			"Pair":          resp.Pair,
			"Today":         time.Now().Format("2006-01-02"),
			"MarketTraders": traders,
			"MarketData":    marketData,
		})
	}
}
func newMarketTradingForm(cfg *config.ConfigManager, tpq querying.TradingPairsQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		slog.Info("Create Trade UI requested")
		id := c.Params("id")
		req := querying.GetTradingPairReq{
			TPID:                 id,
			WithCurrentBasePrice: true,
			WithTrades:           false,
			WithCalculations:     false,
		}
		resp, err := tpq.GetTradingPair(req)
		if err != nil {
			return c.Render("toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   err,
			})
		}
		return c.Render("marketTradingForm", fiber.Map{
			"Pair":          resp.Pair,
			"Today":         time.Now().Format("2006-01-02"),
			"MarketTraders": cfg.GetMarketTraders(true),
		})
	}
}

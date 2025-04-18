package rest

import (
	"controtto/src/app/managing"
	"controtto/src/app/querying"
	"controtto/src/app/trading"
	"controtto/src/domain/pnl"
	"encoding/json"
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

// func fetchMarketTrades(at trading.AssetTrader) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		// Get path parameters
// 		pairID := c.Params("id")
// 		marketKey := c.Params("mtkkey")
// 		fmt.Println("Market Key:", marketKey)
// 		// Parse query parameter
// 		var since time.Time
// 		if sinceStr := c.Query("since"); sinceStr != "" {
// 			parsedSince, err := time.Parse(time.RFC3339, sinceStr)
// 			if err != nil {
// 				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 					"error": "invalid since parameter, must be RFC3339 format",
// 				})
// 			}
// 			since = parsedSince
// 		}
// 		// Create request object
// 		req := trading.FetchTradesReq{
// 			MarketKey: marketKey,
// 			PairID:    pairID,
// 			Since:     since,
// 		}
// 		// Fetch trades using AssetTrader
// 		trades, err := at.FetchTrades(req)
// 		if err != nil {
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 				"error": err.Error(),
// 			})
// 		}
// 		// Return successful response
// 		return c.Status(http.StatusOK).JSON(trades)
// 	}
// }

func fetchMarketTrades(at trading.AssetTrader, mm *managing.MarketManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pairID := c.Params("id")
		marketKey := c.Params("mktkey")
		var since = time.Now().AddDate(0, -6, 0)
		// var since time.Time
		// if sinceStr := c.Query("since"); sinceStr != "" {
		// 	parsedSince, err := time.Parse(time.RFC3339, sinceStr)
		// 	if err != nil {
		// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		// 			"error": "invalid since parameter, must be RFC3339 format",
		// 		})
		// 	}
		// 	since = parsedSince
		// }
		fmt.Println(since.Format(time.RFC3339))
		req := trading.FetchTradesReq{
			MarketKey: marketKey,
			PairID:    pairID,
			Since:     since,
		}
		trades, err := at.FetchTrades(req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		market, err := mm.GetMarket(marketKey)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Market not found",
			})
		}
		return c.Render("marketTradesModal", fiber.Map{
			"Trades":    trades,
			"PairID":    pairID,
			"MarketKey": marketKey,
			"Market":    market,
		})
	}
}

func importMarketTrades(tr trading.TradeRecorder) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pairID := c.Params("id")
		if pairID == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Pair ID is required")
		}
		tradesJSON := c.FormValue("trades")
		if tradesJSON == "" {
			return c.Status(fiber.StatusBadRequest).SendString("No trades data provided")
		}
		var rawTrades struct {
			Trades []struct {
				Timestamp string  `json:"timestamp"`
				Type      string  `json:"type"`
				Price     float64 `json:"price"`
				Amount    float64 `json:"amount"`
				Total     float64 `json:"total"`
				FeeBase   float64 `json:"fee_base"`
				FeeQuote  float64 `json:"fee_quote"`
			} `json:"trades"`
		}
		if err := json.Unmarshal([]byte(tradesJSON), &rawTrades); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Invalid trades data format",
				"details": err.Error(),
			})
		}
		var tradeResponses []trading.RecordTradeResp
		for _, rawTrade := range rawTrades.Trades {
			timestamp, err := time.Parse(time.RFC3339, rawTrade.Timestamp)
			if err != nil {
				return c.Render("toastErr", fiber.Map{
					"Msg": fmt.Sprintf("Could not import trades: %s", err.Error()),
				})
			}
			req := trading.RecordTradeReq{
				PairID:      pairID,
				Timestamp:   timestamp,
				BaseAmount:  rawTrade.Amount,
				QuoteAmount: rawTrade.Total,
				FeeInBase:   rawTrade.FeeBase,
				FeeInQuote:  rawTrade.FeeQuote,
				Type:        rawTrade.Type,
			}
			resp, err := tr.RecordTrade(req)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": fmt.Sprintf("Failed to record trade: %v", err),
				})
			}

			tradeResponses = append(tradeResponses, *resp)
		}
		// Render the updated trades table
		c.Append("HX-Trigger", "newTrade")
		return c.Render("toastOk", fiber.Map{
			"Msg": fmt.Sprintf("Successfully imported %d trades", len(tradeResponses)),
		})
	}
}

package rest

import (
	"controtto/src/app/querying"
	"controtto/src/domain/pnl"
	"encoding/json"
	"log/slog"
	"math/rand"
	"slices"
	"time"

	"github.com/gofiber/fiber/v2"
)

func empty() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.SendString("")
	}
}

func pairCards(tpq querying.TradingPairsQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		req := querying.GetTradingPairReq{
			TPID:             id,
			WithCalculations: true,
		}
		resp, err := tpq.GetTradingPair(req)
		if err != nil {
			return c.Render("toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   err,
			})
		}
		slices.Reverse(resp.Pair.Trades)
		return c.Render("pairCards", fiber.Map{
			"Today": time.Now().Format("Mon Jan 02 15:04 2006"),
			"Pair":  resp.Pair,
		})
	}
}

type PricePoint struct {
	Time  int64   `json:"time"` // Changed to Unix timestamp
	Value float64 `json:"value"`
}

func generateMockPriceData(pairID string, numPoints int) []PricePoint {
	now := time.Now()
	data := make([]PricePoint, numPoints)
	basePrice := 100.0 + rand.Float64()*50 // Initial price
	for i := 0; i < numPoints; i++ {
		time := now.Add(time.Duration(-i) * time.Hour) // Go back in time
		price := basePrice + rand.Float64()*10 - 5  // Price fluctuation
		data[i] = PricePoint{
			Time:  time.Unix(), // Convert to Unix timestamp (seconds)
			Value: price,
		}
	}
	return data
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

		// Calculate Total PNL
		totalPNL := 0.0
		for _, pair := range resp.Pairs {
			totalPNL += pair.Calculations.PNLAmount
		}

		// Generate Mock Chart Data
		chartData := make(map[string][]PricePoint)
		for _, pair := range resp.Pairs {
			chartData[string(pair.ID)] = generateMockPriceData(string(pair.ID), 50) // 50 data points
		}
        

		return c.Render("dashboardSection", fiber.Map{
			"Title":     "Trading Pairs",
			"Pairs":     resp.Pairs,
			"TotalPNL":  totalPNL,
			"ChartData": chartData, // Pass chart data to the template
		})
	}

}

// Home hanlder reders the homescreen
func Home(c *fiber.Ctx) error {
	slog.Info("HOME")
	// render index template
	return c.Render("main", fiber.Map{})
}

func pairChart(tpq querying.TradingPairsQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		req := querying.GetTradingPairReq{
			TPID:             id,
			WithCalculations: false,
		}
		resp, err := tpq.GetTradingPair(req)
		if err != nil {
			return c.Render("toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   err,
			})
		}
		slices.Reverse(resp.Pair.Trades)
		return c.Render("pairChart", fiber.Map{
			"Pair": resp.Pair,
		})
	}
}

func pairTape(tpq querying.TradingPairsQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		req := querying.GetTradingPairReq{
			TPID:             id,
			WithCalculations: false,
		}
		resp, err := tpq.GetTradingPair(req)
		if err != nil {
			return c.Render("toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   err,
			})
		}
		if resp.Pair.BaseAsset.Type == pnl.Stock {
		}
		return c.Render("pairTape", fiber.Map{
			"Pair": resp.Pair,
		})
	}
}

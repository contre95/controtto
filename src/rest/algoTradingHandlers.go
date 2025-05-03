package rest

import (
	"controtto/src/app/managing"
	"controtto/src/app/querying"
	"controtto/src/domain/pnl"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func runSpatialArbitrage() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Set SSE headers
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Transfer-Encoding", "chunked")

		// Get the context to handle client disconnects
		ctx := c.Context()

		// Create a channel to send logs
		logChan := make(chan string)

		// Start a goroutine to simulate log generation
		go generateSpatialArbitrageLogs(logChan, ctx.Done())

		// Stream logs to client
		for {
			select {
			case log := <-logChan:
				// Write SSE formatted message
				_, err := fmt.Fprintf(c, "data: %s\n\n", log)
				if err != nil {
					// Client disconnected
					return nil
				}
				// Flush the response to send immediately
				c.Context().Response.BodyWriter().Flush()
			case <-ctx.Done():
				// Client disconnected
				return nil
			}
		}
	}
}

// generateSpatialArbitrageLogs simulates log
func generateSpatialArbitrageLogs(logChan chan<- string, done <-chan struct{}) {
	defer close(logChan)

	// Example log messages
	messages := []string{
		"Starting spatial arbitrage strategy...",
		"Connected to exchange APIs",
		"Scanning for price discrepancies",
		"Found opportunity: BTC/USD price difference 1.2%",
		"Executing arbitrage trade...",
		"Trade completed successfully",
	}

	for _, msg := range messages {
		select {
		case <-done:
			return
		case logChan <- fmt.Sprintf("[%s] %s", time.Now().Format(time.RFC3339), msg):
			time.Sleep(1 * time.Second) // Simulate delay between logs
		}
	}

	// Keep sending periodic updates
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.Sleep(5 * time.Second):
			logChan <- fmt.Sprintf("[%s] Monitoring markets...", time.Now().Format(time.RFC3339))
		}
	}
}

// Your existing function
func getSpatialArbitrageForm(tpq querying.PairsQuerier, marketManager *managing.MarketManager) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Get available markets
		markets := marketManager.GetMarkets(true)
		marketOptions := make([]pnl.Market, 0, len(markets))
		for _, market := range markets {
			marketOptions = append(marketOptions, *market)
		}
		// Get available trading pairs
		pairsResp, err := tpq.ListPairs(querying.ListPairsReq{})
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to load trading pairs")
		}
		return c.Render("spatialArbitrageForm", fiber.Map{
			"Markets": marketOptions,
			"Pairs":   pairsResp.Pairs,
		})
	}
}

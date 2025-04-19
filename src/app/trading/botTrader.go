package trading

import (
	"controtto/src/app/managing"
	"controtto/src/domain/pnl"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	ErrJobNotFound = errors.New("job not found")
	ErrJobExists   = errors.New("job already exists for this pair")
)

type job struct {
	id        string
	stopChan  chan struct{}
	trades    []pnl.Trade
	logStream chan<- string
}

type TraderBot struct {
	markets     *managing.MarketManager
	pairs       pnl.Pairs
	runningJobs map[string]*job
	mu          sync.RWMutex
}

func NewTraderBot(mm *managing.MarketManager, p pnl.Pairs) *TraderBot {
	return &TraderBot{
		markets:     mm,
		pairs:       p,
		runningJobs: make(map[string]*job),
	}
}

// StartSpatialArbitrage starts a spatial arbitrage job, returns the job id
func (tb *TraderBot) StartSpatialArbitrage(mktkey1, mktkey2, pairID string, logStream chan<- string) (string, error) {
	// Validate markets
	market1, err := tb.markets.GetMarket(mktkey1)
	if err != nil {
		return "", fmt.Errorf("market1: %w", err)
	}
	market2, err := tb.markets.GetMarket(mktkey2)
	if err != nil {
		return "", fmt.Errorf("market2: %w", err)
	}
	pair, err := tb.pairs.GetPair(pairID)
	if err != nil {
		return "", fmt.Errorf("pair: %w", err)
	}
	// Check market health
	if !market1.API.HealthCheck() || !market2.API.HealthCheck() {
		return "", ErrMarketNotHealthy
	}
	fmt.Println("Pair:", pair)
	fmt.Println("Market 1:", market1)
	fmt.Println("Market 2:", market2)
	jobID := uuid.New().String()
	stopChan := make(chan struct{})
	job := &job{
		id:        jobID,
		stopChan:  stopChan,
		trades:    []pnl.Trade{},
		logStream: logStream,
	}
	tb.mu.Lock()
	tb.runningJobs[jobID] = job
	tb.mu.Unlock()

	// Start arbitrage loop in a goroutine
	go func() {
		ticker := time.NewTicker(5 * time.Second) // Check prices every 5 seconds
		defer ticker.Stop()

		for {
			select {
			case <-stopChan:
				return
			case <-ticker.C:
				// Implementation TBD
				logStream <- fmt.Sprintln("Hola")
			}
		}
	}()
	return jobID, nil
}

// Stop stops a running job and returns the trades made by that job
func (tb *TraderBot) Stop(jobID string) ([]pnl.Trade, error) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	job, exists := tb.runningJobs[jobID]
	if !exists {
		return nil, ErrJobNotFound
	}

	// Signal the job to stop
	close(job.stopChan)

	// Collect trades
	trades := job.trades

	// Remove job from running jobs
	delete(tb.runningJobs, jobID)

	// Close log stream if it exists
	if job.logStream != nil {
		close(job.logStream)
	}

	return trades, nil
}

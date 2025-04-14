package managing

import (
	"controtto/src/app/config"
	"controtto/src/domain/pnl"
	"fmt"
	"log/slog"
	"time"
)

// Market Trading Requests and Responses
type MarketBuyReq struct {
	TraderKey     string // Key to identify the MarketTrader
	TradingPairID string
	Amount        float64
	Price         *float64 // nil for market orders
}

type MarketBuyResp struct {
	TradeID string
	Msg     string
}

type MarketSellReq struct {
	TraderKey     string // Key to identify the MarketTrader
	TradingPairID string
	Amount        float64
	Price         *float64 // nil for market orders
}

type MarketSellResp struct {
	TradeID string
	Msg     string
}

type ImportTradesReq struct {
	TraderKey     string // Key to identify the MarketTrader
	TradingPairID string
	Since         time.Time
}

type ImportTradesResp struct {
	Count int
	Msg   string
}

type FetchAssetsReq struct {
	TraderKey string // Key to identify the MarketTrader
}

type MarketTradeManager struct {
	cfg   *config.Config
	pairs pnl.TradingPairs
}

func NewMarketTradeManager(cfg *config.Config, tp pnl.TradingPairs) *MarketTradeManager {
	return &MarketTradeManager{
		cfg:   cfg,
		pairs: tp,
	}
}

func (mtm *MarketTradeManager) MarketBuy(req MarketBuyReq) (*MarketBuyResp, error) {
	trader, exists := mtm.cfg.GetMarketTraders(false)[req.TraderKey]
	if !exists {
		return nil, fmt.Errorf("trader with key %s not found", req.TraderKey)
	}
	tp, err := mtm.pairs.GetTradingPair(req.TradingPairID)
	if err != nil {
		slog.Error("Error getting TradingPair", "error", err)
		return nil, fmt.Errorf("failed to execute buy: %w", err)
	}
	options := pnl.TradeOptions{
		TradingPair:   *tp,
		Amount:        req.Amount,
		Price:         req.Price,
		IsMarketOrder: req.Price == nil,
	}

	trade, err := trader.MarketAPI.Buy(options)
	if err != nil {
		slog.Error("Error executing market buy", "error", err)
		return nil, fmt.Errorf("failed to execute buy: %w", err)
	}

	return &MarketBuyResp{
		TradeID: trade.ID,
		Msg:     "Market buy executed successfully",
	}, nil
}

func (mtm *MarketTradeManager) MarketSell(req MarketSellReq) (*MarketSellResp, error) {
	trader, exists := mtm.cfg.GetMarketTraders(false)[req.TraderKey]
	if !exists {
		return nil, fmt.Errorf("trader with key %s not found", req.TraderKey)
	}
	tp, err := mtm.pairs.GetTradingPair(req.TradingPairID)
	if err != nil {
		slog.
			Error("Error getting TradingPair", "error", err)
		return nil, fmt.Errorf("failed to execute buy: %w", err)
	}
	options := pnl.TradeOptions{
		TradingPair:   *tp,
		Amount:        req.Amount,
		Price:         req.Price,
		IsMarketOrder: req.Price == nil,
	}

	trade, err := trader.MarketAPI.Sell(options)
	if err != nil {
		slog.Error("Error executing market sell", "error", err)
		return nil, fmt.Errorf("failed to execute sell: %w", err)
	}

	return &MarketSellResp{
		TradeID: trade.ID,
		Msg:     "Market sell executed successfully",
	}, nil
}

func (mtm *MarketTradeManager) ImportTrades(req ImportTradesReq) (*ImportTradesResp, error) {
	trader, exists := mtm.cfg.GetMarketTraders(false)[req.TraderKey]
	if !exists {
		return nil, fmt.Errorf("trader with key %s not found", req.TraderKey)
	}
	tp, err := mtm.pairs.GetTradingPair(req.TradingPairID)
	if err != nil {
		slog.Error("Error getting TradingPair", "error", err)
		return nil, fmt.Errorf("failed to execute buy: %w", err)
	}
	trades, err := trader.MarketAPI.ImportTrades(
		*tp,
		req.Since,
	)
	if err != nil {
		slog.Error("Error importing trades", "error", err)
		return nil, fmt.Errorf("failed to import trades: %w", err)
	}

	return &ImportTradesResp{
		Count: len(trades),
		Msg:   fmt.Sprintf("Successfully imported %d trades", len(trades)),
	}, nil
}

package pnl

import (
	"errors"
	"log/slog"
)

func (tp *TradingPair) Calculate() error {
	if err := tp.calculateBuyPrice(); err != nil {
		return err
	}
	if tp.Performance.BaseMarketPrice > 0 {
		if err := tp.calculateProfit(); err != nil {
			return err
		}
	}
	return nil
}

func (tp *TradingPair) calculateProfit() error {
	// Perform any necessary validation or business logic checks here.
	if tp.Performance.BaseMarketPrice == 0 {
		slog.Error("Error calculating P&L", "error", "TradingPair doens't have current base price.")
		return errors.New("Error calculating P&L. No current base price.")
	}
	tp.Performance.PNLAmount = tp.Performance.BaseMarketPrice*tp.Performance.TotalBase - tp.Performance.TotalQuoteSpent
	tp.Performance.CurrentBaseAmountInQuote = tp.Performance.BaseMarketPrice * tp.Performance.TotalBase
	tp.Performance.PNLPercent = (100 * tp.Performance.PNLAmount) / tp.Performance.TotalQuoteSpent
	return nil
}

func (tp *TradingPair) calculateBuyPrice() error {
	// Perform any necessary validation or business logic checks here.
	if len(tp.Trades) == 0 {
		slog.Error("Error calculating P&L", "error", "TradingPair doens't have any trades")
		return errors.New("Please add some trades in order to calculate you profit and loss")
	}
	// tp.Calculations.TotalBase = 0
	// tp.Calculations.TotalQuoteSpent = 0
	for _, t := range tp.Trades {
		if t.TradeType == Buy {
			tp.Performance.TotalBase += t.BaseAmount
			tp.Performance.TotalQuoteSpent += t.QuoteAmount
		}
		if t.TradeType == Sell {
			tp.Performance.TotalBase -= t.BaseAmount
			tp.Performance.TotalQuoteSpent -= t.QuoteAmount
		}
		tp.Performance.TotalFeeInQuote += t.FeeInQuote
		tp.Performance.TotalFeeInBase += t.FeeInBase
		tp.Performance.AvgBuyPrice = float64(tp.Performance.TotalQuoteSpent / tp.Performance.TotalBase)
	}
	tp.Performance.TotalBaseInQuote = tp.Performance.TotalBase * tp.Performance.BaseMarketPrice
	slog.Info("Fields calculated", "base", tp.Performance.TotalBase, "quote", tp.Performance.TotalQuoteSpent, "avg-buy-price", tp.Performance.AvgBuyPrice)
	return nil

}

func (t *Trade) CalculateFields() error {
	t.Price = t.QuoteAmount / t.BaseAmount
	return nil
}

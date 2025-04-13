package pnl

import (
	"errors"
	"log/slog"
)

func (tp *TradingPair) Calculate() error {
	if err := tp.calculateBuyPrice(); err != nil {
		return err
	}
	if tp.Calculations.BasePrice > 0 {
		if err := tp.calculateProfit(); err != nil {
			return err
		}
	}
	return nil
}

func (tp *TradingPair) calculateProfit() error {
	// Perform any necessary validation or business logic checks here.
	if tp.Calculations.BasePrice == 0 {
		slog.Error("Error calculating P&L", "error", "TradingPair doens't have current base price.")
		return errors.New("Error calculating P&L. No current base price.")
	}
	tp.Calculations.PNLAmount = tp.Calculations.BasePrice*tp.Calculations.TotalBase - tp.Calculations.TotalQuoteSpent
	tp.Calculations.CurrentBaseAmountInQuote = tp.Calculations.BasePrice * tp.Calculations.TotalBase
	tp.Calculations.PNLPercent = (100 * tp.Calculations.PNLAmount) / tp.Calculations.TotalQuoteSpent
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
			tp.Calculations.TotalBase += t.BaseAmount
			tp.Calculations.TotalQuoteSpent += t.QuoteAmount
		}
		if t.TradeType == Sell {
			tp.Calculations.TotalBase -= t.BaseAmount
			tp.Calculations.TotalQuoteSpent -= t.QuoteAmount
		}
		tp.Calculations.TotalFeeInQuote += t.FeeInQuote
		tp.Calculations.TotalFeeInBase += t.FeeInBase
		tp.Calculations.AvgBuyPrice = float64(tp.Calculations.TotalQuoteSpent / tp.Calculations.TotalBase)
	}
	tp.Calculations.TotalBaseInQuote = tp.Calculations.TotalBase * tp.Calculations.BasePrice
	slog.Info("Fields calculated", "base", tp.Calculations.TotalBase, "quote", tp.Calculations.TotalQuoteSpent, "avg-buy-price", tp.Calculations.AvgBuyPrice)
	return nil

}

func (t *Trade) CalculateFields() error {
	t.Price = t.QuoteAmount / t.BaseAmount
	return nil
}

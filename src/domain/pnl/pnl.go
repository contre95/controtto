package pnl

import (
	"errors"
	"log/slog"
)

func (tp *Pair) Calculate() error {
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

func (tp *Pair) calculateProfit() error {
	if tp.Calculations.BasePrice == 0 {
		slog.Error("Error calculating P&L", "error", "Pair doesn't have current base price.")
		return errors.New("error calculating P&L: no current base price")
	}
	tp.Calculations.PNLAmount = tp.Calculations.BasePrice*tp.Calculations.TotalBase - tp.Calculations.TotalQuoteSpent
	tp.Calculations.CurrentBaseAmountInQuote = tp.Calculations.BasePrice * tp.Calculations.TotalBase
	if tp.Calculations.TotalQuoteSpent != 0 {
		tp.Calculations.PNLPercent = (100 * tp.Calculations.PNLAmount) / tp.Calculations.TotalQuoteSpent
	}
	return nil
}

func (tp *Pair) calculateBuyPrice() error {
	if len(tp.Trades) == 0 {
		slog.Error("Error calculating P&L", "error", "Pair doesn't have any trades")
		return errors.New("please add some trades to calculate your profit and loss")
	}

	var buyBaseTotal, buyQuoteTotal float64

	for _, t := range tp.Trades {
		if t.TradeType == Buy {
			tp.Calculations.TotalBase += t.BaseAmount
			tp.Calculations.TotalQuoteSpent += t.QuoteAmount

			buyBaseTotal += t.BaseAmount
			buyQuoteTotal += t.QuoteAmount
		}
		if t.TradeType == Sell {
			tp.Calculations.TotalBase -= t.BaseAmount
			tp.Calculations.TotalQuoteSpent -= t.QuoteAmount
		}
		tp.Calculations.TotalFeeInQuote += t.FeeInQuote
		tp.Calculations.TotalFeeInBase += t.FeeInBase
	}

	if buyBaseTotal != 0 {
		tp.Calculations.AvgBuyPrice = buyQuoteTotal / buyBaseTotal
	}

	tp.Calculations.TotalBaseInQuote = tp.Calculations.TotalBase * tp.Calculations.BasePrice

	slog.Info("Fields calculated", "base", tp.Calculations.TotalBase, "quote", tp.Calculations.TotalQuoteSpent, "avg-buy-price", tp.Calculations.AvgBuyPrice)
	return nil
}

func (t *Trade) CalculateFields() error {
	if t.BaseAmount == 0 {
		return errors.New("cannot calculate trade price: base amount is zero")
	}
	t.Price = t.QuoteAmount / t.BaseAmount
	return nil
}

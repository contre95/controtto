package pnl

import "log/slog"

func (tp *TradingPair) Calculate() {
	// Perform any necessary validation or business logic checks here.
	tp.Calculations.TotalBase = 0
	tp.Calculations.TotalQuoteSpent = 0
	for _, t := range tp.Transactions {
		if t.TransactionType == Buy {
			tp.Calculations.TotalBase += t.BaseAmount
			tp.Calculations.TotalQuoteSpent += t.QuoteAmount
		}
		if t.TransactionType == Sell {
			tp.Calculations.TotalBase -= t.BaseAmount
			tp.Calculations.TotalQuoteSpent -= t.QuoteAmount
		}
		tp.Calculations.TotalTradingFeeSpent += t.TradingFee * t.QuoteAmount
		tp.Calculations.TotalWithdrawalFeeSpent += t.WithdrawalFee * t.QuoteAmount
	}
	tp.Calculations.AvgBuyPrice = float64(tp.Calculations.TotalQuoteSpent / tp.Calculations.TotalBase)
	slog.Info("Fields calculated", "base", tp.Calculations.TotalBase, "quote", tp.Calculations.TotalQuoteSpent, "avg-buy-price", tp.Calculations.AvgBuyPrice)
}

func (t *Transaction) CalculateFields() error {
	t.Price = t.QuoteAmount / t.BaseAmount
	return nil
}

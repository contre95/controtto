package trading

type Service struct {
	AssetTrader   AssetTrader
	TradeRecorder TradeRecorder
	TraderBot     *TraderBot
}

func NewService(at AssetTrader, tr TradeRecorder, tb *TraderBot) Service {
	return Service{at, tr, tb}
}

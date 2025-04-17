package trading

type Service struct {
	AssetTrader   AssetTrader
	TradeRecorder TradeRecorder
}

func NewService(at AssetTrader, tr TradeRecorder) Service {
	return Service{at, tr}
}

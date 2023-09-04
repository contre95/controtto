package querying

type Service struct {
	AssetQuerier       AssetsQuerier
	MarketQuerier      MarketsQuerier
	TradingPairQuerier TradingPairsQuerier
}

func NewService(aq AssetsQuerier, mq MarketsQuerier, tpq TradingPairsQuerier) Service {
	return Service{aq, mq,tpq}
}

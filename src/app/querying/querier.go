package querying

// TODO: Maybe all functions just need to bing into this interactor instead of each individual use case (queriers int his case)
// TODO: No, the trading pair querier will receive as well the MarketQuerier use case bundle
type Service struct {
	AssetQuerier       AssetsQuerier
	MarketQuerier      MarketsQuerier
	TradingPairQuerier TradingPairsQuerier
}

func NewService(aq AssetsQuerier, mq MarketsQuerier, tpq TradingPairsQuerier) Service {
	return Service{aq, mq, tpq}
}

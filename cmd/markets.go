package main

import (
	"controtto/src/domain/pnl"
	"controtto/src/gateways/markets"

	"os"
	"strings"
)

var (
	tiingo     = "tiingo"
	bitcoin    = "bitcoin"
	avantage   = "avantage"
	coinbase   = "coinbase"
	binance    = "binance"
	bingx      = "bingx"
	pancake    = "pancake"
	trading212 = "trading212"
)

var marketsConfig = map[string]pnl.Market{
	bitcoin: {
		IsSet:                os.Getenv(PREFIX+strings.ToUpper(bitcoin)+MARKET_SUFIX) != "",
		Env:                  PREFIX + strings.ToUpper(bitcoin) + MARKET_SUFIX,
		MarketName:           "Bitcoin",
		MarketKey:            bitcoin + "_trader",
		Color:                "#FBD8AE",
		Type:                 pnl.Wallet,
		MarketTradingSymbols: []string{"BTC"},
		Details:              "",
		Token:                os.Getenv(PREFIX + strings.ToUpper(bitcoin) + MARKET_SUFIX),
		ProviderURL:          "https://bitcoin-api.github.io/docs/#/en-us/swapV2/introduce",
		MarketLogo:           "/assets/img/" + bitcoin + ".png",
		Init:                 markets.NewMockMarketAPI,
	},
	bingx: {
		IsSet:                os.Getenv(PREFIX+strings.ToUpper(bingx)+MARKET_SUFIX) != "",
		Env:                  PREFIX + strings.ToUpper(bingx) + MARKET_SUFIX,
		MarketName:           "BingX",
		MarketKey:            bingx + "_trader",
		MarketTradingSymbols: []string{"USDC", "USDT"},
		Color:                "#0F5FFF",
		Type:                 pnl.Exchange,
		Details:              "",
		Token:                os.Getenv(PREFIX + strings.ToUpper(bingx) + MARKET_SUFIX),
		ProviderURL:          "https://bingx-api.github.io/docs/#/en-us/swapV2/introduce",
		MarketLogo:           "/assets/img/" + bingx + ".png",
		Init:                 markets.NewBingXAPI,
	},
	pancake: {
		IsSet:                os.Getenv(PREFIX+strings.ToUpper(pancake)+MARKET_SUFIX) != "",
		Env:                  PREFIX + strings.ToUpper(pancake) + MARKET_SUFIX,
		MarketName:           "Pancake",
		MarketTradingSymbols: []string{"USDC", "USDT"},
		MarketKey:            pancake + "_trader",
		Color:                "#23CAD5",
		Details:              "",
		Type:                 pnl.DEX,
		Token:                os.Getenv(PREFIX + strings.ToUpper(pancake) + MARKET_SUFIX),
		ProviderURL:          "https://docs.pancakeswap.finance/developers/api",
		MarketLogo:           "/assets/img/" + pancake + ".png",
		Init:                 markets.NewMockMarketAPI,
	},
	binance: {
		IsSet:                os.Getenv(PREFIX+strings.ToUpper(binance)+MARKET_SUFIX) != "",
		Env:                  PREFIX + strings.ToUpper(binance) + MARKET_SUFIX,
		MarketName:           "Binance",
		MarketKey:            binance + "_trader",
		Details:              "",
		Color:                "#F0B90D",
		MarketTradingSymbols: []string{"USDC", "USDT"},
		Type:                 pnl.Exchange,
		Token:                os.Getenv(PREFIX + strings.ToUpper(binance) + MARKET_SUFIX),
		ProviderURL:          "https://developers.binance.com/en",
		MarketLogo:           "/assets/img/" + binance + ".png",
		Init:                 markets.NewBinanceAPI,
	},
	trading212: {
		IsSet:                os.Getenv(PREFIX+strings.ToUpper(trading212)+MARKET_SUFIX) != "",
		Env:                  PREFIX + strings.ToUpper(trading212) + MARKET_SUFIX,
		MarketName:           "Trading212",
		MarketKey:            trading212 + "_trader",
		Details:              "",
		Color:                "#00AAE4",
		MarketTradingSymbols: []string{"USD", "EUR"},
		Type:                 pnl.Broker,
		Token:                os.Getenv(PREFIX + strings.ToUpper(trading212) + MARKET_SUFIX),
		ProviderURL:          "https://helpcentre.trading212.com/hc/en-us/articles/14584770928157-How-can-I-generate-an-API-key",
		MarketLogo:           "/assets/img/" + trading212 + ".png",
		Init:                 markets.NewTrading212API,
	},
}

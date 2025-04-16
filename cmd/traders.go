package main

import (
	"controtto/src/domain/pnl"
	"controtto/src/gateways/marketTraders"
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

var traders = map[string]pnl.MarketTrader{
	bitcoin: {
		IsSet:       os.Getenv(PREFIX+strings.ToUpper(bitcoin)+TRADER_SUFIX) != "",
		Env:         PREFIX + strings.ToUpper(bitcoin) + TRADER_SUFIX,
		MarketName:  "Bitcoin",
		MarketKey:   bitcoin + "_trader",
		Color:       "#FBD8AE",
		Type:        pnl.Wallet,
		Details:     "",
		Token:       os.Getenv(PREFIX + strings.ToUpper(bitcoin) + TRADER_SUFIX),
		ProviderURL: "https://bitcoin-api.github.io/docs/#/en-us/swapV2/introduce",
		MarketLogo:  "/assets/img/" + bitcoin + ".png",
		Init:        marketTraders.NewBitcoinWalletAPI,
	},
	bingx: {
		IsSet:       os.Getenv(PREFIX+strings.ToUpper(bingx)+TRADER_SUFIX) != "",
		Env:         PREFIX + strings.ToUpper(bingx) + TRADER_SUFIX,
		MarketName:  "BingX",
		MarketKey:   bingx + "_trader",
		Color:       "#0F5FFF",
		Type:        pnl.Exchange,
		Details:     "",
		Token:       os.Getenv(PREFIX + strings.ToUpper(bingx) + TRADER_SUFIX),
		ProviderURL: "https://bingx-api.github.io/docs/#/en-us/swapV2/introduce",
		MarketLogo:  "/assets/img/" + bingx + ".png",
		Init:        marketTraders.NewBingXAPI,
	},
	pancake: {
		IsSet:       os.Getenv(PREFIX+strings.ToUpper(pancake)+TRADER_SUFIX) != "",
		Env:         PREFIX + strings.ToUpper(pancake) + TRADER_SUFIX,
		MarketName:  "Pancake",
		MarketKey:   pancake + "_trader",
		Color:       "#23CAD5",
		Details:     "",
		Type:        pnl.DEX,
		Token:       os.Getenv(PREFIX + strings.ToUpper(pancake) + TRADER_SUFIX),
		ProviderURL: "https://docs.pancakeswap.finance/developers/api",
		MarketLogo:  "/assets/img/" + pancake + ".png",
		Init:        marketTraders.NewMockMarketAPI,
	},
	binance: {
		IsSet:       os.Getenv(PREFIX+strings.ToUpper(binance)+TRADER_SUFIX) != "",
		Env:         PREFIX + strings.ToUpper(binance) + TRADER_SUFIX,
		MarketName:  "Binance",
		MarketKey:   binance + "_trader",
		Details:     "",
		Color:       "#F0B90D",
		Type:        pnl.Exchange,
		Token:       os.Getenv(PREFIX + strings.ToUpper(binance) + TRADER_SUFIX),
		ProviderURL: "https://developers.binance.com/en",
		MarketLogo:  "/assets/img/" + binance + ".png",
		Init:        marketTraders.NewBinanceAPI,
	},
	trading212: {
		IsSet:       os.Getenv(PREFIX+strings.ToUpper(trading212)+TRADER_SUFIX) != "",
		Env:         PREFIX + strings.ToUpper(trading212) + TRADER_SUFIX,
		MarketName:  "Trading212",
		MarketKey:   trading212 + "_trader",
		Details:     "",
		Color:       "#00AAE4",
		Type:        pnl.Broker,
		Token:       os.Getenv(PREFIX + strings.ToUpper(trading212) + TRADER_SUFIX),
		ProviderURL: "https://helpcentre.trading212.com/hc/en-us/articles/14584770928157-How-can-I-generate-an-API-key",
		MarketLogo:  "/assets/img/" + trading212 + ".png",
		Init:        marketTraders.NewTrading212API,
	},
}

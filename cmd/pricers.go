package main

import (
	"controtto/src/domain/pnl"
	"controtto/src/gateways/priceProviders"
	"os"
	"strings"
)

var pricers = map[string]pnl.PriceProvider{
	bingx: {
		IsSet:        os.Getenv(PRICER_PREFIX+strings.ToUpper(bingx)+PUBLIC_PRICE_SUFIX) == "", // Enable by default
		Env:          PRICER_PREFIX + strings.ToUpper(bingx) + PUBLIC_PRICE_SUFIX,
		ProviderName: "BingX",
		ProviderURL:  "https://bingx.com",
		NeedsToken:   false,
		ProviderLogo: "/assets/img/" + bingx + ".png",
		Color:        "#2954FE",
		ProviderKey:  bingx,
		Init:         priceProviders.NewBingxAPI,
	},
	binance: {
		IsSet:        os.Getenv(PRICER_PREFIX+strings.ToUpper(binance)+PUBLIC_PRICE_SUFIX) == "", // Enable by default
		Env:          PRICER_PREFIX + strings.ToUpper(binance) + PUBLIC_PRICE_SUFIX,
		ProviderName: "Binance",
		ProviderURL:  "https://www.binance.com",
		ProviderKey:  binance,
		NeedsToken:   false,
		Color:        "#F0B90D",
		ProviderLogo: "/assets/img/" + binance + ".png",
		Init:         priceProviders.NewBinanceAPI,
	},
	coinbase: {
		IsSet:        os.Getenv(PRICER_PREFIX+strings.ToUpper(coinbase)+PUBLIC_PRICE_SUFIX) == "", // Enable by default
		Env:          PRICER_PREFIX + strings.ToUpper(coinbase) + PUBLIC_PRICE_SUFIX,
		ProviderName: "Coinbase",
		ProviderURL:  "https://docs.cdp.coinbase.com/",
		ProviderKey:  coinbase,
		NeedsToken:   false,
		ProviderLogo: "/assets/img/" + coinbase + ".png",
		Color:        "#0052FF",
		Init:         priceProviders.NewCoinbaseAPI,
	},
	avantage: {
		IsSet:        os.Getenv(PRICER_PREFIX+strings.ToUpper(avantage)+PRIVATE_PRICE_SUFIX) != "",
		Env:          PRICER_PREFIX + strings.ToUpper(avantage) + PRIVATE_PRICE_SUFIX,
		ProviderName: "Alpha Vantage",
		NeedsToken:   true,
		ProviderURL:  "https://www.avantagecapital.com",
		ProviderKey:  avantage,
		ProviderLogo: "/assets/img/" + avantage + ".png",
		Token:        os.Getenv(PRICER_PREFIX + strings.ToUpper(avantage) + PRIVATE_PRICE_SUFIX),
		Color:        "#C2F4E1",
		Init:         priceProviders.NewAVantageAPI,
	},
	tiingo: {
		IsSet:        os.Getenv(PRICER_PREFIX+strings.ToUpper(tiingo)+PRIVATE_PRICE_SUFIX) != "",
		Env:          PRICER_PREFIX + strings.ToUpper(tiingo) + PRIVATE_PRICE_SUFIX,
		NeedsToken:   true,
		ProviderName: "Tiingo",
		ProviderURL:  "https://www.tiingo.com",
		ProviderKey:  tiingo,
		ProviderLogo: "/assets/img/" + tiingo + ".png",
		Token:        os.Getenv(PRICER_PREFIX + strings.ToUpper(tiingo) + PRIVATE_PRICE_SUFIX),
		Color:        "#6E1BD8",
		Init:         priceProviders.NewTiingoAPI,
	},
}

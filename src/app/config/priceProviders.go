package config

import (
	"controtto/src/domain/pnl"
	"controtto/src/gateways/priceProviders"
	"fmt"
	"maps"
	"os"
	"strconv"
	"strings"
)

func (c *ConfigManager) GetPriceProviders() map[string]pnl.PriceProvider {
	c.mu.RLock()
	defer c.mu.RUnlock()
	// Return a copy to prevent external modification
	providers := make(map[string]pnl.PriceProvider, len(c.priceProviders))
	maps.Copy(providers, c.priceProviders)
	return providers
}

func (c *ConfigManager) UpdatePriceProvider(key, token string, toggle bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	provider, ok := c.priceProviders[key]
	if !ok {
		return fmt.Errorf("price provider %q not found", key)
	}
	provider.IsSet = toggle
	if provider.NeedsToken {
		provider.Token = token
		os.Setenv(PREFIX+strings.ToUpper(key)+PRIVATE_PRICE_SUFIX, token)
	} else {
		os.Setenv(PREFIX+strings.ToUpper(key)+PUBLIC_PRICE_SUFIX, strconv.FormatBool(toggle))
	}
	switch key {
	case "tiingo":
		provider.PriceAPI = priceProviders.NewTiingoAPI(token)
	case "avantage":
		provider.PriceAPI = priceProviders.NewAVantageAPI(token)
	case "coinbase":
		provider.PriceAPI = priceProviders.NewCoinbaseAPI()
	case "binance":
		provider.PriceAPI = priceProviders.NewBinanceAPI()
	case "bingx":
		provider.PriceAPI = priceProviders.NewBinanceAPI()
	}
	c.priceProviders[key] = provider
	return nil
}

func loadProviders() map[string]pnl.PriceProvider {
	tiingo := "tiingo"
	avantage := "avantage"
	coinbase := "coinbase"
	binance := "binance"
	bingx := "bingx"
	return map[string]pnl.PriceProvider{
		bingx: {
			IsSet:        os.Getenv(PREFIX+strings.ToUpper(bingx)+PUBLIC_PRICE_SUFIX) == "", // Enable by default
			Env:          PREFIX + strings.ToUpper(bingx) + PUBLIC_PRICE_SUFIX,
			ProviderName: "BingX",
			ProviderURL:  "https://bingx.com",
			NeedsToken:   false,
			ProviderLogo: "/assets/img/" + bingx + ".png",
			Color:        "#2954FE",
			ProviderKey:  bingx + "_toggle",
			PriceAPI:     priceProviders.NewBingxAPI(),
		},
		binance: {
			IsSet:        os.Getenv(PREFIX+strings.ToUpper(binance)+PUBLIC_PRICE_SUFIX) == "", // Enable by default
			Env:          PREFIX + strings.ToUpper(binance) + PUBLIC_PRICE_SUFIX,
			ProviderName: "Binance",
			ProviderURL:  "https://www.binance.com",
			ProviderKey:  binance + "_toggle",
			NeedsToken:   false,
			Color:        "#F0B90D",
			ProviderLogo: "/assets/img/" + binance + ".png",
			PriceAPI:     priceProviders.NewBinanceAPI(),
		},
		coinbase: {
			IsSet:        os.Getenv(PREFIX+strings.ToUpper(coinbase)+PUBLIC_PRICE_SUFIX) == "", // Enable by default
			Env:          PREFIX + strings.ToUpper(coinbase) + PUBLIC_PRICE_SUFIX,
			ProviderName: "Coinbase",
			ProviderURL:  "https://docs.cdp.coinbase.com/",
			ProviderKey:  coinbase + "_toggle",
			NeedsToken:   false,
			ProviderLogo: "/assets/img/" + coinbase + ".png",
			Color:        "#0052FF",
			PriceAPI:     priceProviders.NewCoinbaseAPI(),
		},
		avantage: {
			IsSet:        os.Getenv(PREFIX+strings.ToUpper(avantage)+PRIVATE_PRICE_SUFIX) != "",
			Env:          PREFIX + strings.ToUpper(avantage) + PRIVATE_PRICE_SUFIX,
			ProviderName: "Alpha Vantage",
			NeedsToken:   true,
			ProviderURL:  "https://www.avantagecapital.com",
			ProviderKey:  avantage + "_price_token",
			ProviderLogo: "/assets/img/" + avantage + ".png",
			Token:        os.Getenv(PREFIX + strings.ToUpper(avantage) + PRIVATE_PRICE_SUFIX),
			Color:        "#C2F4E1",
			PriceAPI:     priceProviders.NewAVantageAPI(os.Getenv(PREFIX + strings.ToUpper(avantage) + PRIVATE_PRICE_SUFIX)),
		},
		tiingo: {
			IsSet:        os.Getenv(PREFIX+strings.ToUpper(tiingo)+PRIVATE_PRICE_SUFIX) != "",
			Env:          PREFIX + strings.ToUpper(tiingo) + PRIVATE_PRICE_SUFIX,
			NeedsToken:   true,
			ProviderName: "Tiingo",
			ProviderURL:  "https://www.tiingo.com",
			ProviderKey:  tiingo + "_price_token",
			ProviderLogo: "/assets/img/" + tiingo + ".png",
			Token:        os.Getenv(PREFIX + strings.ToUpper(tiingo) + PRIVATE_PRICE_SUFIX),
			Color:        "#6E1BD8",
			PriceAPI:     priceProviders.NewTiingoAPI(os.Getenv(PREFIX + strings.ToUpper(tiingo) + PRIVATE_PRICE_SUFIX)),
		},
	}
}

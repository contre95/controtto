package config

import (
	"controtto/src/domain/pnl"
	"controtto/src/gateways/marketTraders"
	"fmt"
	"os"
	"strings"
)

func (c *Config) UpdateMarketTraderToken(key, token string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	trader, ok := c.MarketTraders[key]
	if !ok {
		return fmt.Errorf("market trader %q not found", key)
	}
	trader.Token = token
	trader.IsSet = token != ""
	os.Setenv(PREFIX+strings.ToUpper(key)+TRADER_SUFIX, token)
	var api pnl.MarketAPI
	switch key {
	case "t212":
		api = marketTraders.NewTrading212API(token)
	case "bingx":
		api = marketTraders.NewBingXAPI(token)
	case "pancake":
		api = marketTraders.NewMockMarketAPI(token)
	case "binance":
		api = marketTraders.NewMockMarketAPI(token)
	default:
		return fmt.Errorf("market trader %q not found", key)
	}

	trader.MarketAPI = api
	if trader.IsSet {
		if !api.HealthCheck() {
			return fmt.Errorf("Couldn't connect to %s", key)
		}
		details, err := api.AccountDetails()
		fmt.Println("details", details)
		fmt.Println("details", details)
		if err != nil {
			fmt.Println("details", key, err)
			trader.Details = "no details"
		}
		trader.Details = details
	}
	c.MarketTraders[key] = trader
	return nil
}

func (c *Config) GetMarketTraders(all bool) map[string]pnl.MarketTrader {
	c.mu.RLock()
	defer c.mu.RUnlock()

	filtered := make(map[string]pnl.MarketTrader)
	for k, v := range c.MarketTraders {
		if all || v.IsSet {
			filtered[k] = v
		}
	}
	return filtered
}

func loadMarketTraders() map[string]pnl.MarketTrader {
	trading212 := "t212"
	pancake := "pancake"
	bingx := "bingx"
	binance := "binance"
	return map[string]pnl.MarketTrader{
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
			MarketAPI:   marketTraders.NewBingXAPI(os.Getenv(PREFIX + strings.ToUpper(bingx) + TRADER_SUFIX)),
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
			MarketAPI:   marketTraders.NewMockMarketAPI(os.Getenv(PREFIX + strings.ToUpper(pancake) + TRADER_SUFIX)),
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
			MarketAPI:   marketTraders.NewMockMarketAPI(os.Getenv(PREFIX + strings.ToUpper(binance) + TRADER_SUFIX)),
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
			MarketAPI:   marketTraders.NewTrading212API(os.Getenv(PREFIX + strings.ToUpper(trading212) + TRADER_SUFIX)),
		},
	}
}

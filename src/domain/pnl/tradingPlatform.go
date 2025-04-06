package pnl

import "fmt"

type PlatformType int

const (
	PlatformTypeDEX      PlatformType = iota // 0
	PlatformTypeExchange                     // 1
	PlatformTypeBroker                       // 2
)

// String returns the string representation of the PlatformType
func (p PlatformType) String() string {
	return [...]string{"Broker", "Exchange", "DEX"}[p]
}

type TradingPlatform interface {
	// GetTrades returns the list of trades between assetA and assetB
	GetTrades(assetA, assetB string) ([]Trade, error)
	// GetTradedAssets returns the list of assets traded in the platform
	GetTradedAssets() ([]Asset, error)
	// Name returns the name of the trading platform
	Name() string
	// Type returns the type of the trading platform (Broker, Exchange, DEX, etc.)
	Type() PlatformType
}

func ValidatePlatform(p TradingPlatform) error {
	switch p.Type() {
	case PlatformTypeBroker, PlatformTypeExchange, PlatformTypeDEX:
		return nil
	default:
		return fmt.Errorf("invalid platform type: %v", p.Type())
	}
}

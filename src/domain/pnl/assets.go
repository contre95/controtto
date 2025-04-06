package pnl

import (
	"fmt"
	"slices"
)

type InvalidAsset error

// Validate checks the invariants of the Asset
func (a Asset) Validate() (*Asset, error) {
	// Validate Symbol
	if a.Symbol == "" {
		return nil, fmt.Errorf("symbol cannot be empty")
	}

	// Validate Name
	if a.Name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}

	// Validate AssetType
	validTypes := GetValidTypes()
	isValidType := slices.Contains(validTypes, a.Type)
	if !isValidType {
		return nil, fmt.Errorf("invalid asset type: %v", a.Type)
	}

	// Validate CountryCode (optional, depending on asset type)
	if a.Type == Forex && a.CountryCode == "" {
		return nil, fmt.Errorf("country code is required for Forex assets")
	}

	return &a, nil
}

// NewAsset reutrns a new Asset and validates it's invariants
func NewAsset(symbol, color, name, countryCode string, assetType AssetType) (*Asset, error) {
	a := Asset{
		Symbol:      symbol,
		Color:       color,
		Name:        name,
		Type:        assetType,
		CountryCode: countryCode,
	}
	return a.Validate()
}

package pnl

import (
	"errors"
	"regexp"
	"slices"
)

type InvalidAsset error

func (a *Asset) Validate() (*Asset, error) {
	hexColorPattern := `^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`
	re := regexp.MustCompile(hexColorPattern)
	if !re.MatchString(a.Color) {
		return nil, InvalidAsset(errors.New("Wrong color string"))
	}
	if slices.Contains(GetValidTypes(), a.Type) {
		return nil, InvalidAsset(errors.New("Invalid Asset Symbol"))
	}
	if len(a.Symbol) < 3 || len(a.Symbol) > 8 {
		return nil, InvalidAsset(errors.New("Invalid Asset Symbol"))
	}
	if len(a.Name) < 1 || len(a.Name) > 24 {
		return nil, InvalidAsset(errors.New("Invalid Asset Name"))
	}
	return a, nil
}

// TODO: Define invariants
// NewAsset reutrns a new Asset and validates it's invariants
func NewAsset(symbol, color, name, countryCode, assetType string) (*Asset, error) {
	a := Asset{
		Symbol:      symbol,
		Color:       color,
		Name:        name,
		CountryCode: countryCode,
	}
	return a.Validate()
}

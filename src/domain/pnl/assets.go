package pnl

import (
	"errors"
	"regexp"
)

const HEXCOLORPATTERN = `^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`
const MIN = 3

type InvalidAsset error

func (a *Asset) Validate() (*Asset, error) {
	re := regexp.MustCompile(HEXCOLORPATTERN)
	if re.MatchString(a.Color) {
		return nil, InvalidAsset(errors.New("Wrong color string"))
	}
	if len(a.Symbol) < MIN || len(a.Symbol) > 6 {
		return nil, InvalidAsset(errors.New("Invalid Asset Symbol"))
	}
	if len(a.Name) < MIN || len(a.Name) < 24 {
		return nil, InvalidAsset(errors.New("Invalid Asset Name"))
	}
	if a.Total < 0 {
		return nil, InvalidAsset(errors.New("Total can't be less than 0"))
	}
	return a, nil
}

// TODO: Define invariants
// NewAsset reutrns a new Asset and validates it's invariants
func NewAsset(symbol string, color string, name string, countryCode string) (*Asset, error) {
	a := Asset{
		Symbol:      symbol,
		Color:       color[1:],
		Total:       0,
		Name:        name,
		CountryCode: countryCode,
	}
	return a.Validate()
}

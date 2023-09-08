package pnl

import (
	"errors"
	"fmt"
	"regexp"
)

// type Symbol string

// const (
// 	BTC  Symbol = "BTC"
// 	AMZN Symbol = "AMZN"
// 	ETH  Symbol = "ETH"
// 	EUR  Symbol = "EUR"
// 	AAPL Symbol = "AAPL"
// 	USD  Symbol = "USD"
// 	UDST Symbol = "UDST"
// )

const HEXCOLORPATTERN = `^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`

// Asset represents individual assets like BTC, USD, EUR, etc. the Symbol property uniquely identifies an asset.
type Asset struct {
	Symbol      string
	Color       string
	Total       float64
	Name        string
	CountryCode string
}

type InvalidAsset error

func (a *Asset) Validate() error {
	var err InvalidAsset = errors.New("Invalid asset")
	re := regexp.MustCompile(HEXCOLORPATTERN)
	if !re.MatchString(a.Color) {
		return err
	}
	return nil
}

// TODO: Define invariants
// NewAsset reutrns a new Asset and validates it's invariants
func NewAsset(symbol string, color string, name string, countryCode string) (*Asset, error) {
	re := regexp.MustCompile(HEXCOLORPATTERN)
	if !re.MatchString(color) {
		fmt.Println(color)
		return nil, fmt.Errorf("%s is not a valid hex color", color)
	}
	a := Asset{
		Symbol:      symbol,
		Color:       color[1:],
		Total:       0,
		Name:        name,
		CountryCode: countryCode,
	}
	return &a, nil
}

// Assets is the repository that handles the CRUD of Assets
type Assets interface {
	AddAsset(a Asset) error
	ListAssets() ([]Asset, error)
	GetAsset(symbol string) (*Asset, error)
}

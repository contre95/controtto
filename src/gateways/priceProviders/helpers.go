package priceProviders

import (
	"strconv"
)

// Helper function to convert a string to a float64.
func stringToFloat64(s string) (float64, error) {
	price, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return price, nil
}

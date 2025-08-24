package ui

import (
	"fmt"
	"strings"
)

func DisplayPrice(v float64) string {
	s := fmt.Sprintf("%.2f", v)
	n := len(s)
	if n <= 6 { // e.g. 999.99
		return s
	}
	// Insert commas for thousands
	parts := strings.Split(s, ".")
	intPart := parts[0]
	decPart := parts[1]
	var out []byte
	for i, c := range intPart {
		if i != 0 && (len(intPart)-i)%3 == 0 {
			out = append(out, ',')
		}
		out = append(out, byte(c))
	}
	return string(out) + "." + decPart
}

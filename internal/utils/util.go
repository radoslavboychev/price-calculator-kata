package utils

import "github.com/radoslavboychev/price-calculator-kata/internal/utils/format"

// AmountFromPercentage calculates the absolute amount from percentage and returns the result as a float with 4 decimal precision
func AmountFromPercentage(percentage uint16, price float64) (result float64) {
	return format.ToDecimal(((float64(percentage) / 100) * price), 4)
}

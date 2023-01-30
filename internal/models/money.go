package models

import (
	"github.com/radoslavboychev/price-calculator-kata/internal/currency"
	"github.com/radoslavboychev/price-calculator-kata/internal/utils/format"
)

// Money struct represents currency with amount and sign
type Money struct {
	Currency currency.CurrencyCode
	Value    float64
}

// NewMoney constructor function for Money types
func NewMoney(currency currency.CurrencyCode, value float64) Money {
	if value < 0 {
		value = 0
	}

	return Money{
		Currency: currency,
		Value:    format.ToDecimal(value, 4),
	}
}

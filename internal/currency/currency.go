package currency

import "github.com/radoslavboychev/price-calculator-kata/config"

// Enum for various currency code types
const (
	USD CurrencyCode = iota
	GBP
	JPY
	EUR
)

// CurrencyCode defines an enum for the different currencies
type CurrencyCode uint16

// Currency stores the ISO-3 currency code
type Currency struct {
	Code CurrencyCode
}

func NewCurrency(code uint16) *Currency {
	return &Currency{
		Code: CurrencyCode(code),
	}
}

// LoadCurrency reads the currency from the config and loads it
func LoadCurrency() *Currency {

	conf := config.LoadConfig()
	if conf.Currency > 3 {
		conf.Currency = 0
	}

	return &Currency{
		Code: CurrencyCode(conf.Currency),
	}
}

// String repreents a currency as string for printing purposes
func (c CurrencyCode) String() string {
	switch c {
	case USD:
		return "USD"
	case GBP:
		return "GBP"
	case JPY:
		return "JPY"
	case EUR:
		return "EUR"
	}
	return "UNKNOWN CURRENCY"
}

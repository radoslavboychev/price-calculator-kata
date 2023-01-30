package cap

import (
	"github.com/radoslavboychev/price-calculator-kata/config"
	"github.com/radoslavboychev/price-calculator-kata/internal/models"
)

// Cap interface defines behavior for all types which implement it
type DiscountCap interface {
	CalculateCap(startingPrice models.Money, discount float64) float64
}

// CapAbsolute represents discount cap based on absolute value
type capAbsolute struct {
	Value models.Money
}

// CapPercentage represents discount cap based on percentage values
type capPercentage struct {
	Value float64
}

// CalculateCap calculates the cap amount for absolute cap values
func (c *capAbsolute) CalculateCap(startingPrice models.Money, discount float64) float64 {

	if discount > c.Value.Value {
		discount = c.Value.Value
	}
	return discount
}

// CalculateCap calculates the cap amount for percentage-based cap values
func (c *capPercentage) CalculateCap(startingPrice models.Money, discount float64) float64 {
	if c.Value == 0 {
		c.Value = 100
	}
	capAmount := (c.Value / 100) * startingPrice.Value
	if discount > capAmount {
		discount = capAmount
	}
	return discount
}

// newCapPercentage constructor function for percentage based discount caps
func newCapPercentage(value float64) *capPercentage {
	// if cap is 0%, set it to 100% to basically remove it
	if value == 0 {
		value = 100
	}

	if value < 0 {
		value = 100
	}

	return &capPercentage{
		Value: value,
	}
}

// newCapAbsolute constructor function for absolute value based discount caps
func newCapAbsolute(value float64) *capAbsolute {
	// if the cap is 0, it is not valid, set it very high to basically remove it
	if value == 0 {
		value = 1000000
	}

	//  if the cap is a negative number, it is not valid, set it very high to basically remove it
	if value < 0 {
		value = 1000000
	}

	return &capAbsolute{
		Value: models.Money{
			Value: value,
		},
	}
}

// NewDiscountCap checks the type of discount cap defined in the config (absolute or percentage) and returns a new instance of the cap
// with the values from config
// if an invalid value for the cap is set, returns a new cap that is set to 100% of the product price,
// meaning a cap would practically not exist
func NewDiscountCap(value float64) DiscountCap {
	conf := config.LoadConfig()

	switch conf.CapType {
	case 1:
		return newCapPercentage(value)
	case 2:
		return newCapAbsolute(value)
	default:
		return newCapPercentage(100)
	}
}

// NewDiscountCapTesting slightly different method of generating the cap for testing purposes.
// Instead of reading
func NewDiscountCapTesting(capType uint16, value float64) DiscountCap {
	switch capType {
	case 1:
		return newCapAbsolute(value)
	case 2:
		return newCapPercentage(value)
	default:
		return newCapPercentage(100)
	}
}

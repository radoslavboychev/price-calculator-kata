package models

// Tax contains a tax rate and a tax amount
type Tax struct {
	rate   uint16
	Amount Money
}

// NewTax constructor function for Tax types
func NewTax(rate uint16) *Tax {
	if rate > 100 {
		rate = 100
	}

	return &Tax{
		rate: rate,
	}
}

// Rate returns the tax rate
func (t *Tax) Rate() uint16 {
	return t.rate
}

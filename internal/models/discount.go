package models

// Enum for the types of discount precedence
const (
	NoPrecedence TakesPrecedence = iota
	PrecedenceUniversal
	PrecedenceSpecial
)

// Enum to indicate if a discount takes precedence over tax
type TakesPrecedence uint16

// Discount contains the 2 different types of discounts
type Discount struct {
	UniversalDiscount universalDiscount
	SpecialDiscount   specialDiscount
	TakesPrecedence   TakesPrecedence
}

// Universal Discount that can apply to all products
type universalDiscount struct {
	rate   uint16
	Amount Money
}

// Special Discount that applies to products with specified UPC
type specialDiscount struct {
	upc    int
	rate   uint16
	Amount Money
}

// NewUniversalDiscount constructor function for universal discounts
func NewUniversalDiscount(rate uint16, amount Money) *universalDiscount {
	if rate > 100 {
		rate = 100
	}

	return &universalDiscount{
		rate:   rate,
		Amount: amount,
	}
}

// NewSpecialDiscount constructor function for special discounts
func NewSpecialDiscount(upc int, rate uint16, amount Money) *specialDiscount {
	if upc < 0 {
		upc = 0
	}

	if rate > 100 {
		rate = 100
	}

	return &specialDiscount{
		upc:    upc,
		rate:   rate,
		Amount: amount,
	}
}

// NewDiscount constructor function for new discounts
func NewDiscount(universal universalDiscount, special specialDiscount, precedence TakesPrecedence) *Discount {

	return &Discount{
		UniversalDiscount: universal,
		SpecialDiscount:   special,
		TakesPrecedence:   precedence,
	}
}

// Rate returns the discount rate for universal discounts
func (d *universalDiscount) Rate() uint16 {
	return d.rate
}

// Rate returns the discount rate for special discounts
func (d *specialDiscount) Rate() uint16 {
	return d.rate
}

// UPC returns a special discount's UPC
func (s *specialDiscount) UPC() int {
	return s.upc
}

package calculator

import (
	"github.com/radoslavboychev/price-calculator-kata/internal/cap"
	"github.com/radoslavboychev/price-calculator-kata/internal/combining"
	"github.com/radoslavboychev/price-calculator-kata/internal/models"
	"github.com/radoslavboychev/price-calculator-kata/internal/utils"
	"github.com/radoslavboychev/price-calculator-kata/internal/utils/format"
	"github.com/radoslavboychev/price-calculator-kata/pkg/result"
)

// calculator struct that will store the types needed for performing calculations
type calculator struct {
	tax         models.Tax
	discount    models.Discount
	combineType combining.CombType
	cap         cap.DiscountCap
}

// NewCalculator constructor returns a new calculator and initializes the values
func NewCalculator(tax models.Tax, discount models.Discount, combineType combining.CombType, discountCap cap.DiscountCap) *calculator {

	return &calculator{
		tax:         tax,
		discount:    discount,
		combineType: combineType,
		cap:         discountCap,
	}
}

// Calculate runs the calculations for a specific product depending on the various conditions that could be met, and reports the results.
func (c *calculator) Calculate(p *models.Product) *result.Result {
	startingPrice := p.Price()
	productPrice := p.Price()

	switch c.discount.TakesPrecedence {
	case 1:
		c.discount.UniversalDiscount.Amount.Value = utils.AmountFromPercentage(c.discount.UniversalDiscount.Rate(), productPrice.Value)
		c.tax.Amount.Value = utils.AmountFromPercentage(c.tax.Rate(), productPrice.Value-c.discount.UniversalDiscount.Amount.Value)
		if c.discount.SpecialDiscount.UPC() == p.UPC() {
			c.discount.SpecialDiscount.Amount.Value = utils.AmountFromPercentage(c.discount.SpecialDiscount.Rate(), productPrice.Value-c.discount.UniversalDiscount.Amount.Value)
		}
	case 2:
		if c.discount.SpecialDiscount.UPC() == p.UPC() {
			c.discount.SpecialDiscount.Amount.Value = utils.AmountFromPercentage(c.discount.SpecialDiscount.Rate(), productPrice.Value)
		}
		c.tax.Amount.Value = utils.AmountFromPercentage(c.tax.Rate(), productPrice.Value-c.discount.SpecialDiscount.Amount.Value)
		c.discount.UniversalDiscount.Amount.Value = utils.AmountFromPercentage(c.discount.UniversalDiscount.Rate(), productPrice.Value-c.discount.SpecialDiscount.Amount.Value)

	default:
		c.tax.Amount.Value = utils.AmountFromPercentage(c.tax.Rate(), productPrice.Value)
		c.discount.UniversalDiscount.Amount.Value = utils.AmountFromPercentage(c.discount.UniversalDiscount.Rate(), productPrice.Value)
		if c.discount.SpecialDiscount.UPC() == p.UPC() {
			c.discount.SpecialDiscount.Amount.Value = utils.AmountFromPercentage(c.discount.SpecialDiscount.Rate(), productPrice.Value)
		}
	}

	var sumDiscount float64

	if c.combineType == combining.TypeAdditive {
		sumDiscount = c.cap.CalculateCap(startingPrice, (c.discount.SpecialDiscount.Amount.Value + c.discount.UniversalDiscount.Amount.Value))
		productPrice.Value = (startingPrice.Value + c.tax.Amount.Value) - sumDiscount
	} else {
		if c.discount.SpecialDiscount.UPC() == p.UPC() {
			c.discount.SpecialDiscount.Amount.Value = utils.AmountFromPercentage(c.discount.SpecialDiscount.Rate(), startingPrice.Value-c.discount.UniversalDiscount.Amount.Value)
		}
		sumDiscount = c.cap.CalculateCap(startingPrice, c.discount.UniversalDiscount.Amount.Value+c.discount.SpecialDiscount.Amount.Value)
		productPrice.Value = (startingPrice.Value + c.tax.Amount.Value) - sumDiscount
	}

	costs := calculateCosts(p.Cost(), startingPrice)

	productPrice.Value += costs

	resCurrency := p.Price().Currency
	res := result.NewResult(

		models.NewMoney(resCurrency, format.ToDecimal(startingPrice.Value, 2)),
		models.NewMoney(resCurrency, format.ToDecimal(c.tax.Amount.Value, 2)),
		models.NewMoney(resCurrency, format.ToDecimal(sumDiscount, 2)),
		models.NewMoney(resCurrency, format.ToDecimal(costs, 2)),
		models.NewMoney(resCurrency, format.ToDecimal(productPrice.Value, 2)),
		p.Cost(),
	)

	return res
}

// calculateCosts calculates and returns a sum of all expenses
func calculateCosts(costs models.Costs, startingPrice models.Money) float64 {
	var sum float64

	for _, cost := range costs.Expenses {
		val := cost.CalculateExpense(startingPrice)

		sum += val.Value
	}
	return format.ToDecimal(sum, 4)
}

package calculator

import (
	"log"
	"testing"

	"github.com/radoslavboychev/price-calculator-kata/internal/cap"
	"github.com/radoslavboychev/price-calculator-kata/internal/combining"
	"github.com/radoslavboychev/price-calculator-kata/internal/currency"
	"github.com/radoslavboychev/price-calculator-kata/internal/models"
	"github.com/radoslavboychev/price-calculator-kata/internal/utils"
	"github.com/radoslavboychev/price-calculator-kata/internal/utils/format"
	"github.com/radoslavboychev/price-calculator-kata/pkg/result"

	"github.com/stretchr/testify/assert"
)

func TestCalculate(t *testing.T) {
	// Test case for when the product UPC is too short and invalid. Checks if the UPC generation method is correct
	t.Run("TEST_PRODUCT_UPC_TOO_SHORT", func(t *testing.T) {

		// Arrange
		var expectedUPC int = 123
		var randomUPC int = 123

		// Act
		p1 := models.NewProduct("Random", randomUPC, models.Money{}, models.NewCosts())
		log.Printf("Generated UPC: %v", p1.UPC())

		// Assert
		assert.NotEqual(t, expectedUPC, p1.UPC())

	})

	// Test case for when a product with empty values is passed into
	t.Run("TEST_PRODUCT_VALUES_EMPTY", func(t *testing.T) {

		// Arrange
		p := models.NewProduct("", 0, models.Money{}, models.NewCosts())

		// Assert
		assert.NotEmpty(t, p.Name(), p.Cost(), p.Price(), p.UPC(), p.Price().Value, p.Price().Currency)
	})

	// Testing if the limits in tax rate and discount rates are calculated correctly if invalid amounts have been inserted initially (too high)
	t.Run("TEST_CALCULATOR_LIMITS", func(t *testing.T) {

		tax := *models.NewTax(500)
		discount := *models.NewDiscount(
			*models.NewUniversalDiscount(300, models.Money{}),
			*models.NewSpecialDiscount(0, 5000, models.Money{}),
			models.NoPrecedence,
		)

		// Arrange
		var expectedTaxRate uint16 = 100
		var expectedDiscountRate uint16 = 100
		var expectedSpecialDiscountRate uint16 = 100

		// Act
		calc := NewCalculator(tax, discount, combining.TypeAdditive, cap.NewDiscountCapTesting(0, 100))

		// Assert
		assert.Equal(t, expectedTaxRate, calc.tax.Rate())
		assert.Equal(t, expectedDiscountRate, calc.discount.UniversalDiscount.Rate())
		assert.Equal(t, expectedSpecialDiscountRate, calc.discount.SpecialDiscount.Rate())
	})

	t.Run("TEST_CALCULATOR_NIL_VALUES", func(t *testing.T) {

		tax := *models.NewTax(0)
		discount := *models.NewDiscount(*models.NewUniversalDiscount(0, models.Money{}),
			*models.NewSpecialDiscount(0, 0, models.Money{}),
			models.NoPrecedence,
		)

		// Arrange
		var expectedTaxRate uint16 = 0
		var expectedDiscountRate uint16 = 0
		var expectedSpecialDiscountRate uint16 = 0

		// Act
		calc := NewCalculator(tax, discount, combining.TypeAdditive, cap.NewDiscountCapTesting(0, 100))

		// Assert
		assert.Equal(t, expectedTaxRate, calc.tax.Rate())
		assert.Equal(t, expectedDiscountRate, calc.discount.UniversalDiscount.Rate())
		assert.Equal(t, expectedSpecialDiscountRate, calc.discount.SpecialDiscount.Rate())
	})

	// Testing the TAX requirement - tax calculation
	t.Run("TEST_TAX_REQUIREMENT", func(t *testing.T) {

		tax := *models.NewTax(20)
		discount := models.Discount{
			UniversalDiscount: *models.NewUniversalDiscount(0, models.Money{}),
			SpecialDiscount:   *models.NewSpecialDiscount(0, 0, models.Money{}),
			TakesPrecedence:   models.NoPrecedence,
		}

		p := models.NewProduct("The Little Prince", 123456, models.NewMoney(0, 20.25), models.NewCosts())

		calc := NewCalculator(tax, discount, combining.TypeAdditive, cap.NewDiscountCapTesting(0, 100))

		// Arrange
		expectedTotal := 24.30

		// Act
		res := calc.Calculate(&p)

		// Assert
		assert.Equal(t, expectedTotal, res.TotalPrice().Value)
	})

	// Tests DISCOUNT requirement - calculating and applying discount
	t.Run("TEST_DISCOUNT_REQUIREMENT", func(t *testing.T) {

		p := models.NewProduct("The Little Prince", 123456, models.NewMoney(0, 20.25), models.NewCosts())

		tax := *models.NewTax(20)
		discount := *models.NewDiscount(
			*models.NewUniversalDiscount(15, models.Money{}),
			*models.NewSpecialDiscount(0, 0, models.Money{}),
			models.NoPrecedence,
		)

		calc := NewCalculator(tax, discount, combining.TypeAdditive, cap.NewDiscountCapTesting(0, 100))

		// Arrange
		expectedTax := 4.05
		expectedDiscount := 3.04
		expectedTotal := 21.26

		// Act
		res := calc.Calculate(&p)

		// Assert
		assert.Equal(t, expectedDiscount, res.TotalDiscount().Value)
		assert.Equal(t, expectedTax, res.TaxAmount().Value)
		assert.Equal(t, expectedTotal, res.TotalPrice().Value)

	})

	// Tests REPORT requirement for printing different reports in different conditions
	t.Run("TEST_REPORT_REQUIREMENT", func(t *testing.T) {

		p := models.NewProduct("The Little Prince", 123456, models.NewMoney(0, 20.25), models.NewCosts())

		tax1 := *models.NewTax(20)
		discount1 := *models.NewDiscount(
			*models.NewUniversalDiscount(15, models.Money{}),
			*models.NewSpecialDiscount(0, 0, models.Money{}),
			models.NoPrecedence,
		)

		tax2 := *models.NewTax(20)
		discount2 := *models.NewDiscount(
			*models.NewUniversalDiscount(0, models.Money{}),
			*models.NewSpecialDiscount(0, 0, models.Money{}),
			models.NoPrecedence,
		)

		case1 := NewCalculator(tax1, discount1, combining.TypeAdditive, cap.NewDiscountCapTesting(0, 100))
		case2 := NewCalculator(tax2, discount2, combining.TypeAdditive, cap.NewDiscountCapTesting(0, 100))

		// Arrange
		resCase1 := case1.Calculate(&p)
		resCase2 := case2.Calculate(&p)

		// Act
		report1 := resCase1.Report()
		report2 := resCase2.Report()

		// Assert
		assert.Contains(t, report1, "Discounts")
		assert.NotContains(t, report2, "Discounts")

	})

	// Tests SELECTIVE requirement for special UPC discounts
	t.Run("TEST_SELECTIVE_REQUIREMENT", func(t *testing.T) {

		p := models.NewProduct("The Little Prince", 123456, models.NewMoney(0, 20.25), models.NewCosts())

		tax := *models.NewTax(20)
		discount := *models.NewDiscount(
			*models.NewUniversalDiscount(15, models.Money{}),
			*models.NewSpecialDiscount(123456, 7, models.Money{}),
			models.NoPrecedence,
		)

		calc := NewCalculator(tax, discount, combining.TypeAdditive, cap.NewDiscountCapTesting(0, 100))

		// Arrange
		expectedTax := 4.05
		expectedDiscounts := 4.46
		expectedTotal := 19.85

		// Act
		res := calc.Calculate(&p)

		// Assert
		assert.Equal(t, expectedTax, res.TaxAmount().Value)
		assert.Equal(t, expectedDiscounts, res.TotalDiscount().Value)
		assert.Equal(t, expectedTotal, res.TotalPrice().Value)
	})

	// Tests PRECEDENCE requirement for applying tax before or after specific discounts. Special discount applies before tax
	t.Run("TEST_PRECEDENCE_REQUIREMENT", func(t *testing.T) {

		p := models.NewProduct("The Little Prince", 123456, models.NewMoney(0, 20.25), models.NewCosts())

		tax := *models.NewTax(20)
		discount := *models.NewDiscount(
			*models.NewUniversalDiscount(15, models.Money{}),
			*models.NewSpecialDiscount(123456, 7, models.Money{}),
			models.PrecedenceSpecial,
		)

		calc := NewCalculator(tax, discount, combining.TypeAdditive, cap.NewDiscountCapTesting(0, 100))

		// Arrange
		expectedUPCDiscount := 1.42
		expectedUniversalDiscount := 2.82
		expectedTotalDiscount := expectedUPCDiscount + expectedUniversalDiscount
		expectedFinalPrice := 19.77
		expectedTax := 3.77

		// Act
		res := calc.Calculate(&p)

		// Assert
		assert.Equal(t, expectedTotalDiscount, res.TotalDiscount().Value)
		assert.Equal(t, expectedTax, res.TaxAmount().Value)
		assert.Equal(t, expectedFinalPrice, res.TotalPrice().Value)
	})

	// Tests the case where universal discount takes precedence over tax
	t.Run("TEST_PRECEDENCE_REQUIREMENT_UNIVERSAL_DISCOUNT_TAKES_PRECEDENCE", func(t *testing.T) {

		p := models.NewProduct("The Little Prince", 123456, models.NewMoney(0, 20.25), models.NewCosts())

		tax := *models.NewTax(20)
		discount := *models.NewDiscount(
			*models.NewUniversalDiscount(15, models.Money{}),
			*models.NewSpecialDiscount(123456, 7, models.Money{}),
			models.PrecedenceUniversal,
		)

		calc := NewCalculator(tax, discount, combining.TypeAdditive, cap.NewDiscountCapTesting(0, 100))

		// Arrange
		expectedUPCDiscount := 1.42
		expectedUniversalDiscount := 2.82
		expectedTotalDiscount := expectedUPCDiscount + expectedUniversalDiscount
		expectedFinalPrice := 19.45
		expectedTax := 3.44

		// Act
		res := calc.Calculate(&p)

		// Assert
		assert.Equal(t, expectedTotalDiscount, res.TotalDiscount().Value)
		assert.Equal(t, expectedTax, res.TaxAmount().Value)
		assert.Equal(t, expectedFinalPrice, res.TotalPrice().Value)
	})

	// Tests EXPENSE requirement where costs can be applied onto a price
	t.Run("TEST_EXPENSE_REQUIREMENT", func(t *testing.T) {

		// Arrange
		costAbsolute := models.NewExpenseAbsolute("Transport", 2.2)
		costPercentage := models.NewExpensePercentage("Packaging", 1)

		costs := models.NewCosts(costAbsolute, costPercentage)

		p := models.NewProduct("The Little Prince", 123456, models.NewMoney(0, 20.25), models.NewCosts(costs))

		tax := *models.NewTax(21)
		discount := *models.NewDiscount(
			*models.NewUniversalDiscount(15, models.Money{}),
			*models.NewSpecialDiscount(123456, 7, models.Money{}),
			models.NoPrecedence,
		)

		calc := NewCalculator(tax, discount, combining.TypeAdditive, cap.NewDiscountCapTesting(0, 100))

		// Arrange
		expectedTax := 4.25
		expectedDiscount := 4.46
		expectedTotal := 22.45

		// Act
		res := calc.Calculate(&p)

		// Assert
		assert.Equal(t, expectedTax, res.TaxAmount().Value)
		assert.Equal(t, expectedDiscount, res.TotalDiscount().Value)
		assert.Equal(t, expectedTotal, res.TotalPrice().Value)
	})

	// Tests COMBINING requirement where there are two different methods of combining discounts
	t.Run("TEST_COMBINING_REQUIREMENT_ADDITIVE", func(t *testing.T) {

		// Arrange
		costAbsolute := models.NewExpenseAbsolute("Transport", 2.2)
		costPercentage := models.NewExpensePercentage("Packaging", 1)

		costs := models.NewCosts(costAbsolute, costPercentage)

		p := models.NewProduct("The Little Prince", 123456, models.NewMoney(0, 20.25), models.NewCosts(costs))

		tax := *models.NewTax(21)

		discount := *models.NewDiscount(
			*models.NewUniversalDiscount(15, models.Money{}),
			*models.NewSpecialDiscount(123456, 7, models.Money{}),
			models.NoPrecedence,
		)

		calc := NewCalculator(tax, discount, combining.TypeAdditive, cap.NewDiscountCapTesting(0, 100))

		// Arrange
		expectedTax := 4.25
		expectedDiscountsAdditive := 4.46
		expectedTotalAdditive := 22.45

		// Act
		resAdditive := calc.Calculate(&p)

		// Assert
		assert.Equal(t, expectedTax, resAdditive.TaxAmount().Value)
		assert.Equal(t, expectedDiscountsAdditive, resAdditive.TotalDiscount().Value)
		assert.Equal(t, expectedTotalAdditive, resAdditive.TotalPrice().Value)
	})

	t.Run("TEST_COMBINING_REQUIREMENT_MULTIPLICATIVE", func(t *testing.T) {
		// Arrange
		costAbsolute := models.NewExpenseAbsolute("Transport", 2.2)
		costPercentage := models.NewExpensePercentage("Packaging", 1)

		costs := models.NewCosts(costAbsolute, costPercentage)

		p := models.NewProduct("The Little Prince", 123456, models.NewMoney(0, 20.25), models.NewCosts(costs))

		tax := *models.NewTax(21)

		discount := *models.NewDiscount(
			*models.NewUniversalDiscount(15, models.Money{}),
			*models.NewSpecialDiscount(123456, 7, models.Money{}),
			models.NoPrecedence,
		)

		calc := NewCalculator(tax, discount, combining.TypeMultiplicative, cap.NewDiscountCapTesting(0, 100))

		// Arrange
		expectedTax := 4.25
		expectedDiscountsMultiplicative := 4.24
		expectedTotalMultiplicative := 22.66

		// Act
		resMultiplicative := calc.Calculate(&p)

		// Assert
		assert.Equal(t, expectedTax, resMultiplicative.TaxAmount().Value)
		assert.Equal(t, expectedDiscountsMultiplicative, resMultiplicative.TotalDiscount().Value)
		assert.Equal(t, expectedTotalMultiplicative, resMultiplicative.TotalPrice().Value)
	})

	// Tests CURRENCY requirement where different currencies can be used. Checking default case for USD
	t.Run("TEST_CURRENCY_USD", func(t *testing.T) {
		p := models.NewProduct("The Little Prince", 123456, models.NewMoney(0, 20.25), models.NewCosts())

		tax := *models.NewTax(20)
		discount := *models.NewDiscount(
			*models.NewUniversalDiscount(0, models.Money{}),
			*models.NewSpecialDiscount(0, 0, models.Money{}),
			models.NoPrecedence)

		calc := NewCalculator(tax, discount, combining.TypeAdditive, cap.NewDiscountCapTesting(0, 100))

		// Arrange
		expectedTax := 4.05
		expectedTotal := 24.30
		expectedCurrencyTax := currency.USD
		expectedCurrencyTotal := currency.USD
		expectedCurrencyStarting := currency.USD

		// Act
		res := calc.Calculate(&p)

		// Assert
		assert.Equal(t, expectedTax, res.TaxAmount().Value)
		assert.Equal(t, expectedCurrencyTotal, res.TotalPrice().Currency)
		assert.Equal(t, expectedTotal, res.TotalPrice().Value)
		assert.Equal(t, expectedCurrencyStarting, res.StartingPrice().Currency)
		assert.Equal(t, expectedCurrencyTax, res.TaxAmount().Currency)

	})

	// Tests CURRENCY requirement where different currencies can be used. Checking default case for GBP
	t.Run("TEST_CURRENCY_GBP", func(t *testing.T) {
		p := models.NewProduct("The Little Prince", 123456, models.NewMoney(1, 17.76), models.NewCosts())

		tax := *models.NewTax(20)
		discount := *models.NewDiscount(
			*models.NewUniversalDiscount(0, models.Money{}),
			*models.NewSpecialDiscount(0, 0, models.Money{}),
			models.NoPrecedence)

		calc := NewCalculator(tax, discount, combining.TypeAdditive, cap.NewDiscountCapTesting(0, 100))

		// Arrange
		expectedTax := 3.55
		expectedTotal := 21.31
		expectedCurrencyTax := currency.GBP
		expectedCurrencyTotal := currency.GBP
		expectedCurrencyStarting := currency.GBP

		// Act
		res := calc.Calculate(&p)

		// Assert
		assert.Equal(t, expectedTax, res.TaxAmount().Value)
		assert.Equal(t, expectedCurrencyTotal, res.TotalPrice().Currency)
		assert.Equal(t, expectedTotal, res.TotalPrice().Value)
		assert.Equal(t, expectedCurrencyStarting, res.StartingPrice().Currency)
		assert.Equal(t, expectedCurrencyTax, res.TaxAmount().Currency)

	})

	// Tests PRECISION requirement where calculations need to be performed with 4 decimal places precision but reported with 2 decimal places
	t.Run("TEST_PRECISION", func(t *testing.T) {
		// Arrange
		costPercentage := models.NewExpensePercentage("Packaging", 3)

		costs := models.NewCosts(costPercentage)

		p := models.NewProduct("The Little Prince", 123456, models.NewMoney(0, 20.25), models.NewCosts(costs))

		tax := *models.NewTax(21)
		discount := models.NewDiscount(*models.NewUniversalDiscount(15, models.Money{}),
			*models.NewSpecialDiscount(123456, 7, models.Money{}),
			models.NoPrecedence)

		calc := NewCalculator(tax, *discount, combining.TypeMultiplicative, cap.NewDiscountCapTesting(0, 100))

		// Arrange

		// Calculated high-precision numbers (4 decimals)
		expectedTaxAmountPrecise := 4.2525
		expectedUniversalDiscountPrecise := 3.0375
		expectedSpecialDiscountPrecise := 1.2049
		expectedTotalDiscountPrecise := 4.2424
		expectedCostsPrecise := 0.6075

		// Resulting values (2 decimals)
		expectedStartingPrice := 20.25
		expectedTax := 4.25
		expectedDiscounts := 4.24
		expectedTotal := 20.87

		// Act
		res, taxAmountPrecise, universalDiscountPrecise, specialDiscountPrecise, totalDiscountPrecise, costsPrecise := calc.calculatePrecision(&p)

		// Assert

		// Calculated high-precision numbers (4 decimals)
		assert.Equal(t, expectedTaxAmountPrecise, taxAmountPrecise)
		assert.Equal(t, expectedUniversalDiscountPrecise, universalDiscountPrecise)
		assert.Equal(t, expectedSpecialDiscountPrecise, specialDiscountPrecise)
		assert.Equal(t, expectedTotalDiscountPrecise, totalDiscountPrecise)
		assert.Equal(t, expectedCostsPrecise, costsPrecise)

		// Resulting values (2 decimals)
		assert.Equal(t, expectedStartingPrice, res.StartingPrice().Value)
		assert.Equal(t, expectedTax, res.TaxAmount().Value)
		assert.Equal(t, expectedDiscounts, res.TotalDiscount().Value)
		assert.Equal(t, expectedTotal, res.TotalPrice().Value)

	})

	t.Run("TEST_PRECISION_ADDITIVE_COMBINING", func(t *testing.T) {

		costPercentage := models.NewExpensePercentage("Packaging", 3)

		costs := models.NewCosts(costPercentage)

		p := models.NewProduct("The Little Prince", 123456, models.NewMoney(0, 20.25), models.NewCosts(costs))

		tax := *models.NewTax(21)
		discounts := *models.NewDiscount(
			*models.NewUniversalDiscount(15, models.Money{}),
			*models.NewSpecialDiscount(123456, 7, models.Money{}),
			models.NoPrecedence)

		calc := NewCalculator(tax, discounts, combining.TypeAdditive, cap.NewDiscountCapTesting(0, 100))

		// Arrange

		// Calculated high-precision numbers (4 decimals)
		expectedTaxAmountPrecise := 4.2525
		expectedUniversalDiscountPrecise := 3.0375
		expectedSpecialDiscountPrecise := 1.4175
		expectedTotalDiscountPrecise := 4.455
		expectedCostsPrecise := 0.6075

		// Resulting values (2 decimals)
		expectedStartingPrice := 20.25
		expectedTax := 4.25
		expectedDiscounts := 4.46
		expectedTotal := 20.66

		// Act
		res, taxAmountPrecise, universalDiscountPrecise, specialDiscountPrecise, totalDiscountPrecise, costsPrecise := calc.calculatePrecision(&p)

		// Assert

		// Calculated high-precision numbers (4 decimals)
		assert.Equal(t, expectedTaxAmountPrecise, taxAmountPrecise)
		assert.Equal(t, expectedUniversalDiscountPrecise, universalDiscountPrecise)
		assert.Equal(t, expectedSpecialDiscountPrecise, specialDiscountPrecise)
		assert.Equal(t, expectedTotalDiscountPrecise, totalDiscountPrecise)
		assert.Equal(t, expectedCostsPrecise, costsPrecise)

		// Resulting values (2 decimals)
		assert.Equal(t, expectedStartingPrice, res.StartingPrice().Value)
		assert.Equal(t, expectedTax, res.TaxAmount().Value)
		assert.Equal(t, expectedDiscounts, res.TotalDiscount().Value)
		assert.Equal(t, expectedTotal, res.TotalPrice().Value)

	})

	t.Run("TEST_PRECISION_UNIVERSAL_DISCOUNT_PRECEDENCE", func(t *testing.T) {

		costPercentage := models.NewExpensePercentage("Packaging", 3)

		costs := models.NewCosts(costPercentage)

		p := models.NewProduct("The Little Prince", 123456, models.NewMoney(0, 20.25), models.NewCosts(costs))

		tax := *models.NewTax(21)
		discount := *models.NewDiscount(
			*models.NewUniversalDiscount(15, models.Money{}),
			*models.NewSpecialDiscount(123456, 7, models.Money{}),
			models.PrecedenceUniversal)

		calc := NewCalculator(tax, discount, combining.TypeAdditive, cap.NewDiscountCapTesting(0, 100))

		// Arrange

		// Calculated high-precision numbers (4 decimals)
		expectedTaxAmountPrecise := 3.6146
		expectedUniversalDiscountPrecise := 3.0375
		expectedSpecialDiscountPrecise := 1.2049
		expectedTotalDiscountPrecise := 4.2424
		expectedCostsPrecise := 0.6075

		// Resulting values (2 decimals)
		expectedStartingPrice := 20.25
		expectedTax := 3.61
		expectedDiscounts := 4.24
		expectedTotal := 20.23

		// Act
		res, taxAmountPrecise, universalDiscountPrecise, specialDiscountPrecise, totalDiscountPrecise, costsPrecise := calc.calculatePrecision(&p)

		// Assert

		// Calculated high-precision numbers (4 decimals)
		assert.Equal(t, expectedTaxAmountPrecise, taxAmountPrecise)
		assert.Equal(t, expectedUniversalDiscountPrecise, universalDiscountPrecise)
		assert.Equal(t, expectedSpecialDiscountPrecise, specialDiscountPrecise)
		assert.Equal(t, expectedTotalDiscountPrecise, totalDiscountPrecise)
		assert.Equal(t, expectedCostsPrecise, costsPrecise)

		// Resulting values (2 decimals)
		assert.Equal(t, expectedStartingPrice, res.StartingPrice().Value)
		assert.Equal(t, expectedTax, res.TaxAmount().Value)
		assert.Equal(t, expectedDiscounts, res.TotalDiscount().Value)
		assert.Equal(t, expectedTotal, res.TotalPrice().Value)

	})

	t.Run("TEST_PRECISION_SPECIAL_DISCOUNT_PRECEDENCE", func(t *testing.T) {

		costPercentage := models.NewExpensePercentage("Packaging", 3)

		costs := models.NewCosts(costPercentage)

		p := models.NewProduct("The Little Prince", 123456, models.NewMoney(0, 20.25), models.NewCosts(costs))

		tax := *models.NewTax(21)
		discount := *models.NewDiscount(
			*models.NewUniversalDiscount(15, models.Money{}),
			*models.NewSpecialDiscount(123456, 7, models.Money{}),
			models.PrecedenceSpecial,
		)

		calc := NewCalculator(tax, discount, combining.TypeAdditive, cap.NewDiscountCapTesting(0, 100))

		// Arrange

		// Calculated high-precision numbers (4 decimals)
		expectedTaxAmountPrecise := 3.9548
		expectedUniversalDiscountPrecise := 2.8249
		expectedSpecialDiscountPrecise := 1.4175
		expectedTotalDiscountPrecise := 4.2424
		expectedCostsPrecise := 0.6075

		// Resulting values (2 decimals)
		expectedStartingPrice := 20.25
		expectedTax := 3.95
		expectedDiscounts := 4.24
		expectedTotal := 20.57

		// Act
		res, taxAmountPrecise, universalDiscountPrecise, specialDiscountPrecise, totalDiscountPrecise, costsPrecise := calc.calculatePrecision(&p)

		// Assert

		// Calculated high-precision numbers (4 decimals)
		assert.Equal(t, expectedTaxAmountPrecise, taxAmountPrecise)
		assert.Equal(t, expectedUniversalDiscountPrecise, universalDiscountPrecise)
		assert.Equal(t, expectedSpecialDiscountPrecise, specialDiscountPrecise)
		assert.Equal(t, expectedTotalDiscountPrecise, totalDiscountPrecise)
		assert.Equal(t, expectedCostsPrecise, costsPrecise)

		// Resulting values (2 decimals)
		assert.Equal(t, expectedStartingPrice, res.StartingPrice().Value)
		assert.Equal(t, expectedTax, res.TaxAmount().Value)
		assert.Equal(t, expectedDiscounts, res.TotalDiscount().Value)
		assert.Equal(t, expectedTotal, res.TotalPrice().Value)

	})
}

func TestCalculateCosts(t *testing.T) {

	t.Run("TEST_CALCULATE_COSTS_NEGATIVE_PRICE", func(t *testing.T) {
		// Arrange
		costAbsolute := models.NewExpenseAbsolute("Transport", 2.2)
		costPercentage := models.NewExpensePercentage("Packaging", 1)
		allExpenses := models.NewCosts(costAbsolute, costPercentage)

		// Act
		res := calculateCosts(allExpenses, models.NewMoney(currency.USD, -25))
		expectedResult := 2.2

		// Assert
		assert.Equal(t, res, expectedResult)
	})

	t.Run("TEST_CALCULATE_COSTS_ONLY_ABSOLUTE", func(t *testing.T) {
		// Arrange
		costAbsolute := models.NewExpenseAbsolute("Transport", 2.2)
		allExpenses := models.NewCosts(costAbsolute)

		// Act
		res := calculateCosts(allExpenses, models.NewMoney(currency.USD, 5))
		expectedResult := 2.2

		// Assert
		assert.Equal(t, res, expectedResult)
	})

	t.Run("TEST_CALCULATE_COSTS_ONLY_PERCENTAGE", func(t *testing.T) {
		// Arrange
		costPercentage := models.NewExpensePercentage("Packaging", 10)
		allExpenses := models.NewCosts(costPercentage)

		// Act
		res := calculateCosts(allExpenses, models.NewMoney(currency.USD, 50))
		var expectedResult float64 = 5

		// Assert
		assert.Equal(t, res, expectedResult)
	})
}

func TestDiscountCap(t *testing.T) {

	// Tests CAP requirement where discounts can have a specific cap. Testing case where it's a percentage-based cap
	t.Run("TEST_CAP_REQUIREMENT_PERCENTAGE", func(t *testing.T) {
		p := models.NewProduct("The Little Prince", 123456, models.NewMoney(0, 20.25), models.NewCosts())

		tax := *models.NewTax(21)
		discount := *models.NewDiscount(
			*models.NewUniversalDiscount(15, models.Money{}),
			*models.NewSpecialDiscount(123456, 7, models.Money{}),
			models.NoPrecedence,
		)

		calc := NewCalculator(tax, discount, combining.TypeAdditive, cap.NewDiscountCapTesting(2, 20))

		// Arrange
		expectedTax := 4.25
		expectedDiscounts := 4.05
		expectedTotal := 20.45

		// Act
		res := calc.Calculate(&p)

		// Assert
		assert.Equal(t, expectedTax, res.TaxAmount().Value)
		assert.Equal(t, expectedDiscounts, res.TotalDiscount().Value)
		assert.Equal(t, expectedTotal, res.TotalPrice().Value)
	})

	// Tests CAP requirement where discounts can have a specific cap. Testing case where it's an absolute amount
	t.Run("TEST_CAP_REQUIREMENT_ABSOLUTE", func(t *testing.T) {
		p := models.NewProduct("The Little Prince", 123456, models.NewMoney(0, 20.25), models.NewCosts())

		tax := *models.NewTax(21)
		discount := *models.NewDiscount(
			*models.NewUniversalDiscount(15, models.Money{}),
			*models.NewSpecialDiscount(123456, 7, models.Money{}),
			models.NoPrecedence,
		)

		calc := NewCalculator(tax, discount, combining.TypeAdditive, cap.NewDiscountCapTesting(1, 4))

		// Arrange
		expectedTax := 4.25
		expectedDiscounts := 4.00
		expectedTotal := 20.50

		// Act
		res := calc.Calculate(&p)

		// Assert
		assert.Equal(t, expectedTax, res.TaxAmount().Value)
		assert.Equal(t, expectedDiscounts, res.TotalDiscount().Value)
		assert.Equal(t, expectedTotal, res.TotalPrice().Value)
	})

	// Tests CAP requirement where discounts can have a specific cap. Testing case where it's a percentage-based cap with the second set of parameters from the example
	t.Run("TEST_CAP_REQUIREMENT_PERCENTAGE_SECOND", func(t *testing.T) {
		p := models.NewProduct("The Little Prince", 123456, models.NewMoney(0, 20.25), models.NewCosts())

		tax := *models.NewTax(21)
		discount := *models.NewDiscount(
			*models.NewUniversalDiscount(15, models.Money{}),
			*models.NewSpecialDiscount(123456, 7, models.Money{}),
			models.NoPrecedence,
		)

		calc := NewCalculator(tax, discount, combining.TypeAdditive, cap.NewDiscountCapTesting(2, 30))

		// Arrange
		expectedTax := 4.25
		expectedDiscounts := 4.46
		expectedTotal := 20.05

		// Act
		res := calc.Calculate(&p)

		// Assert
		assert.Equal(t, expectedTax, res.TaxAmount().Value)
		assert.Equal(t, expectedDiscounts, res.TotalDiscount().Value)
		assert.Equal(t, expectedTotal, res.TotalPrice().Value)
	})
}

// calculatePrecision functions the same as the regular Calculate() method but returns amounts with 4 decimal precision for testing purposes
func (c *calculator) calculatePrecision(p *models.Product) (res *result.Result, taxAmountPrecise, universalDiscountPrecise, specialDiscountPrecise, totalDiscountPrecise, costsPrecise float64) {
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

	taxAmountPrecise = c.tax.Amount.Value
	universalDiscountPrecise = c.discount.UniversalDiscount.Amount.Value
	specialDiscountPrecise = c.discount.SpecialDiscount.Amount.Value
	totalDiscountPrecise = sumDiscount
	costsPrecise = costs

	resCurrency := p.Price().Currency
	res = result.NewResult(

		models.NewMoney(resCurrency, format.ToDecimal(startingPrice.Value, 2)),
		models.NewMoney(resCurrency, format.ToDecimal(c.tax.Amount.Value, 2)),
		models.NewMoney(resCurrency, format.ToDecimal(sumDiscount, 2)),
		models.NewMoney(resCurrency, format.ToDecimal(costs, 2)),
		models.NewMoney(resCurrency, format.ToDecimal(productPrice.Value, 2)),
		p.Cost(),
	)

	res.Report()

	// for testing the PRECISION requirement
	return res, taxAmountPrecise, universalDiscountPrecise, specialDiscountPrecise, totalDiscountPrecise, costsPrecise
}

func BenchmarkCalculate(b *testing.B) {

	b.Run("BENCHMARK_TAX_REQUIREMENT", func(b *testing.B) {
		tax := *models.NewTax(20)
		discount := models.Discount{
			UniversalDiscount: *models.NewUniversalDiscount(0, models.Money{}),
			SpecialDiscount:   *models.NewSpecialDiscount(0, 0, models.Money{}),
			TakesPrecedence:   models.NoPrecedence,
		}

		p := models.NewProduct("The Little Prince", 123456, models.NewMoney(0, 20.25), models.NewCosts())

		calc := NewCalculator(tax, discount, combining.TypeAdditive, cap.NewDiscountCapTesting(0, 100))

		// Act
		calc.Calculate(&p)

	})

	b.Run("BENCHMARK_DISCOUNT_REQUIREMENT", func(b *testing.B) {
		p := models.NewProduct("The Little Prince", 123456, models.NewMoney(0, 20.25), models.NewCosts())

		tax := *models.NewTax(20)
		discount := *models.NewDiscount(
			*models.NewUniversalDiscount(15, models.Money{}),
			*models.NewSpecialDiscount(0, 0, models.Money{}),
			models.NoPrecedence,
		)

		calc := NewCalculator(tax, discount, combining.TypeAdditive, cap.NewDiscountCapTesting(0, 100))

		calc.Calculate(&p)
	})

	b.Run("BENCHMARK_REPORT_REQUIREMENT", func(b *testing.B) {
		p := models.NewProduct("The Little Prince", 123456, models.NewMoney(0, 20.25), models.NewCosts())

		tax1 := *models.NewTax(20)
		discount1 := *models.NewDiscount(
			*models.NewUniversalDiscount(15, models.Money{}),
			*models.NewSpecialDiscount(0, 0, models.Money{}),
			models.NoPrecedence,
		)

		case1 := NewCalculator(tax1, discount1, combining.TypeAdditive, cap.NewDiscountCapTesting(0, 100))

		// Arrange
		case1.Calculate(&p)

	})

	b.Run("BENCHMARK_SELECTIVE_REQUIREMENT", func(b *testing.B) {
		p := models.NewProduct("The Little Prince", 123456, models.NewMoney(0, 20.25), models.NewCosts())

		tax := *models.NewTax(20)
		discount := *models.NewDiscount(
			*models.NewUniversalDiscount(15, models.Money{}),
			*models.NewSpecialDiscount(123456, 7, models.Money{}),
			models.NoPrecedence,
		)

		calc := NewCalculator(tax, discount, combining.TypeAdditive, cap.NewDiscountCapTesting(0, 100))

		// Act
		calc.Calculate(&p)

	})

	b.Run("BENCHMARK_PRECEDENCE_REQUIREMENT", func(b *testing.B) {
		p := models.NewProduct("The Little Prince", 123456, models.NewMoney(0, 20.25), models.NewCosts())

		tax := *models.NewTax(20)
		discount := *models.NewDiscount(
			*models.NewUniversalDiscount(15, models.Money{}),
			*models.NewSpecialDiscount(123456, 7, models.Money{}),
			models.PrecedenceSpecial,
		)

		calc := NewCalculator(tax, discount, combining.TypeAdditive, cap.NewDiscountCapTesting(0, 100))

		// Act
		calc.Calculate(&p)

	})

	b.Run("BENCHMARK_EXPENSE_REQUIREMENT", func(b *testing.B) {
		// Arrange
		costAbsolute := models.NewExpenseAbsolute("Transport", 2.2)
		costPercentage := models.NewExpensePercentage("Packaging", 1)

		costs := models.NewCosts(costAbsolute, costPercentage)

		p := models.NewProduct("The Little Prince", 123456, models.NewMoney(0, 20.25), models.NewCosts(costs))

		tax := *models.NewTax(21)
		discount := *models.NewDiscount(
			*models.NewUniversalDiscount(15, models.Money{}),
			*models.NewSpecialDiscount(123456, 7, models.Money{}),
			models.NoPrecedence,
		)

		calc := NewCalculator(tax, discount, combining.TypeAdditive, cap.NewDiscountCapTesting(0, 100))

		// Act
		calc.Calculate(&p)

	})

	b.Run("BENCHMARK_COMBINING_REQUIREMENT", func(b *testing.B) {
		// Arrange
		costAbsolute := models.NewExpenseAbsolute("Transport", 2.2)
		costPercentage := models.NewExpensePercentage("Packaging", 1)

		costs := models.NewCosts(costAbsolute, costPercentage)

		p := models.NewProduct("The Little Prince", 123456, models.NewMoney(0, 20.25), models.NewCosts(costs))

		tax := *models.NewTax(21)

		discount := *models.NewDiscount(
			*models.NewUniversalDiscount(15, models.Money{}),
			*models.NewSpecialDiscount(123456, 7, models.Money{}),
			models.NoPrecedence,
		)

		calc := NewCalculator(tax, discount, combining.TypeAdditive, cap.NewDiscountCapTesting(0, 100))

		// Act
		calc.Calculate(&p)

	})

	b.Run("BENCHMARK_CURRENCY_USD", func(b *testing.B) {
		p := models.NewProduct("The Little Prince", 123456, models.NewMoney(0, 20.25), models.NewCosts())

		tax := *models.NewTax(20)
		discount := *models.NewDiscount(
			*models.NewUniversalDiscount(0, models.Money{}),
			*models.NewSpecialDiscount(0, 0, models.Money{}),
			models.NoPrecedence)

		calc := NewCalculator(tax, discount, combining.TypeAdditive, cap.NewDiscountCapTesting(0, 100))

		// Act
		calc.Calculate(&p)
	})

	b.Run("BENCHMARK_PRECISION_REQUIREMENT", func(b *testing.B) {

		// Arrange
		costPercentage := models.NewExpensePercentage("Packaging", 3)

		costs := models.NewCosts(costPercentage)

		p := models.NewProduct("The Little Prince", 123456, models.NewMoney(0, 20.25), models.NewCosts(costs))

		tax := *models.NewTax(21)
		discount := models.NewDiscount(*models.NewUniversalDiscount(15, models.Money{}),
			*models.NewSpecialDiscount(123456, 7, models.Money{}),
			models.NoPrecedence)

		calc := NewCalculator(tax, *discount, combining.TypeMultiplicative, cap.NewDiscountCapTesting(0, 100))

		// Act
		calc.calculatePrecision(&p)

	})
}

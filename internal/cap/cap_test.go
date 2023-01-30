package cap

import (
	"testing"

	"github.com/radoslavboychev/price-calculator-kata/internal/currency"
	"github.com/radoslavboychev/price-calculator-kata/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestCalculateCap(t *testing.T) {
	// Case for calculating the discount cap from absolute amount
	t.Run("CALCULATE_CAP_ABSOLUTE_VALUE", func(t *testing.T) {
		// Arrange
		cap := newCapAbsolute(2)

		p := models.NewProduct("The Little Prince", 123456, models.NewMoney(0, 20.25), models.NewCosts())
		discount := models.NewDiscount(*models.NewUniversalDiscount(20, models.NewMoney(currency.USD, 5)),
			*models.NewSpecialDiscount(p.UPC(), 0, models.NewMoney(currency.USD, 0)), models.NoPrecedence)

		var expectedResult float64 = 2

		// Act
		res := cap.CalculateCap(p.Price(), discount.UniversalDiscount.Amount.Value)

		// Assert
		assert.Equal(t, expectedResult, res)

	})

	// Case for calculating the discount cap from percentage
	t.Run("CALCULATE_CAP_PERCENTAGE", func(t *testing.T) {
		// Arrange
		cap := newCapPercentage(10)
		p := models.NewProduct("The Little Prince", 123456, models.NewMoney(0, 20.25), models.NewCosts())
		discount := models.NewDiscount(*models.NewUniversalDiscount(20, models.NewMoney(currency.USD, 5)),
			*models.NewSpecialDiscount(p.UPC(), 0, models.NewMoney(currency.USD, 0)), models.NoPrecedence)

		var expectedResult float64 = 2.025

		// Act
		res := cap.CalculateCap(p.Price(), discount.UniversalDiscount.Amount.Value)

		// Assert
		assert.Equal(t, expectedResult, res)

	})

	// Case for calculating the discount cap from percentage when percentage is set to zero
	// cap set to zero will be changed to a value of 100, therefore removing it, there will be no cap, same applies to negative cap
	t.Run("CALCULATE_CAP_PERCENTAGE_CAP_IS_ZERO", func(t *testing.T) {
		// Arrange
		cap := newCapPercentage(0)
		p := models.NewProduct("The Little Prince", 123456, models.NewMoney(0, 20.25), models.NewCosts())
		discount := models.NewDiscount(*models.NewUniversalDiscount(20, models.NewMoney(currency.USD, 5)),
			*models.NewSpecialDiscount(p.UPC(), 0, models.NewMoney(currency.USD, 0)), models.NoPrecedence)

		var expectedResult float64 = 5

		// Act
		res := cap.CalculateCap(p.Price(), discount.UniversalDiscount.Amount.Value)

		// Assert
		assert.Equal(t, expectedResult, res)

	})

	// Case for calculating the discount cap from absolute amount when cap is set to zero
	// Cap set to zero will be changed to a value of 100, therefore removing it, same applies with negative cap
	t.Run("CALCULATE_CAP_ABSOLUTE_VALUE_CAP_IS_ZERO", func(t *testing.T) {
		// Arrange
		cap := newCapAbsolute(0)

		p := models.NewProduct("The Little Prince", 123456, models.NewMoney(0, 20.25), models.NewCosts())
		discount := models.NewDiscount(*models.NewUniversalDiscount(20, models.NewMoney(currency.USD, 5)),
			*models.NewSpecialDiscount(p.UPC(), 0, models.NewMoney(currency.USD, 0)), models.NoPrecedence)

		var expectedResult float64 = 5

		// Act
		res := cap.CalculateCap(p.Price(), discount.UniversalDiscount.Amount.Value)

		// Assert
		assert.Equal(t, expectedResult, res)

	})

	// Case for when discount is zero
	t.Run("CALCULATE_CAP_ABSOLUTE_DISCOUNT_IS_ZERO", func(t *testing.T) {
		// Arrange
		cap := newCapAbsolute(5)

		p := models.NewProduct("The Little Prince", 123456, models.NewMoney(0, 20.25), models.NewCosts())

		var expectedResult float64 = 0

		// Act
		res := cap.CalculateCap(p.Price(), 0)

		// Assert
		assert.Equal(t, expectedResult, res)

	})
}

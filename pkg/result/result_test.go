package result

import (
	"testing"

	"github.com/radoslavboychev/price-calculator-kata/internal/currency"
	"github.com/radoslavboychev/price-calculator-kata/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestReport(t *testing.T) {

	// Case when all entries are present
	t.Run("TEST_REPORT_ALL_PRESENT", func(t *testing.T) {
		// Arrange
		startingPrice := models.NewMoney(currency.USD, 20.25)
		taxAmount := models.NewMoney(currency.USD, 3)
		totalDiscount := models.NewMoney(currency.USD, 2)
		totalExpenses := models.NewMoney(currency.USD, 2)
		totalPrice := models.NewMoney(currency.USD, 10)
		costs := models.NewCosts()

		// Act
		r := NewResult(startingPrice, taxAmount, totalDiscount, totalExpenses, totalPrice, costs)
		str := r.Report()

		// Assert
		assert.Contains(t, str, "Cost")
		assert.Contains(t, str, "Tax")
		assert.Contains(t, str, "Discounts")
		assert.Contains(t, str, "TOTAL")
	})

	// Case when tax is nil
	t.Run("TEST_REPORT_NO_TAX", func(t *testing.T) {
		// Arrange

		// tax is nil
		taxAmount := models.Money{}
		startingPrice := models.NewMoney(currency.USD, 20.25)
		totalDiscount := models.NewMoney(currency.USD, 2)
		totalExpenses := models.NewMoney(currency.USD, 2)
		totalPrice := models.NewMoney(currency.USD, 10)
		costs := models.NewCosts()

		// Act
		r := NewResult(startingPrice, taxAmount, totalDiscount, totalExpenses, totalPrice, costs)
		str := r.Report()

		// Assert
		assert.Contains(t, str, "Cost")
		assert.NotContains(t, str, "Tax")
		assert.Contains(t, str, "Discounts")
		assert.Contains(t, str, "TOTAL")
	})

	// Case when discounts are nil
	t.Run("TEST_REPORT_NO_DISCOUNT", func(t *testing.T) {
		// Arrange
		totalDiscount := models.Money{}
		startingPrice := models.NewMoney(currency.USD, 20.25)
		taxAmount := models.NewMoney(currency.USD, 3)
		totalExpenses := models.NewMoney(currency.USD, 2)
		totalPrice := models.NewMoney(currency.USD, 10)
		costs := models.NewCosts()

		// Act
		r := NewResult(startingPrice, taxAmount, totalDiscount, totalExpenses, totalPrice, costs)
		str := r.Report()

		// Assert
		assert.Contains(t, str, "Cost")
		assert.Contains(t, str, "Tax")
		assert.NotContains(t, str, "Discounts")
		assert.Contains(t, str, "TOTAL")
	})

	t.Run("TEST_REPORT_ALL_IS_NIL", func(t *testing.T) {
		// Arrange
		startingPrice := models.Money{}
		taxAmount := models.Money{}
		totalDiscount := models.Money{}
		totalExpenses := models.Money{}
		totalPrice := models.Money{}
		costs := models.NewCosts()

		// Act
		r := NewResult(startingPrice, taxAmount, totalDiscount, totalExpenses, totalPrice, costs)
		str := r.Report()

		// Assert
		// "Cost" and "TOTAL" will always be a part of the report even if nil
		assert.Contains(t, str, "Cost")
		assert.NotContains(t, str, "Tax")
		assert.NotContains(t, str, "Discounts")
		assert.Contains(t, str, "TOTAL")
	})

	t.Run("TEST_REPORT_", func(t *testing.T) {

	})
}

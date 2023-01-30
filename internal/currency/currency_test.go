package currency

import (
	"testing"

	"github.com/radoslavboychev/price-calculator-kata/config"
	"github.com/stretchr/testify/assert"
)

func TestLoadCurrency(t *testing.T) {
	t.Run("LOAD_CURRENCY_DEFAULT_CASE", func(t *testing.T) {

		// Arrange
		conf := config.LoadConfig()
		conf.Currency = 0

		expectedCurrencyCode := conf.Currency
		expectedResult := NewCurrency(expectedCurrencyCode)

		// Act
		res := LoadCurrency()
		res.Code = 0

		// Assert
		assert.Equal(t, expectedResult, res)
	})

	t.Run("LOAD_CURRENCY_INVALID_CODE", func(t *testing.T) {

		// Arrange
		conf := config.LoadConfig()
		conf.Currency = 0

		expectedCurrencyCode := conf.Currency
		expectedResult := NewCurrency(expectedCurrencyCode)

		// Act
		res := loadCurrencyTest(7)

		// Assert
		assert.Equal(t, expectedResult, res)
	})
}

// loadCurrencyTest loads a currency with a specific currency code
// if the code is bigger than the final iota, it gets set to 0 (default)
func loadCurrencyTest(code uint16) *Currency {

	if code > 3 {
		code = 0
	}

	return &Currency{
		Code: CurrencyCode(code),
	}
}

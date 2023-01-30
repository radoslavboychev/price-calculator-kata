package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAmountFromPercentage(t *testing.T) {
	// Test case when percentage is zero, should return 0
	t.Run("AMOUNT_FROM_PERCENTAGE_ZERO_PERCENT", func(t *testing.T) {
		var percentage uint16 = 0
		var price float64 = 20

		var expectedResult float64 = 0.00

		res := AmountFromPercentage(percentage, price)

		assert.Equal(t, expectedResult, res)
	})

	// Test case when percentage is properly set
	t.Run("AMOUNT_FROM_PERCENTAGE_DEFAULT", func(t *testing.T) {
		var percentage uint16 = 10
		var price float64 = 20

		var expectedResult float64 = 2

		res := AmountFromPercentage(percentage, price)

		assert.Equal(t, expectedResult, res)
	})
}

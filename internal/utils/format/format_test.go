package format

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToDecimal(t *testing.T) {
	t.Run("TEST_TODECIMAL_CONVERSION", func(t *testing.T) {
		// Arrange
		var input float64 = 20.2564353
		var places int = 4

		expectedResult := 20.2564

		// Act
		res := ToDecimal(input, places)

		// Assert
		assert.Equal(t, expectedResult, res)
	})

	t.Run("TEST_TODECIMAL_NIL_INPUT", func(t *testing.T) {
		// Arrange
		var input float64 = 0
		var places int = 4

		var expectedResult float64 = 0

		// Act
		res := ToDecimal(input, places)

		// Assert
		assert.Equal(t, expectedResult, res)
	})

	t.Run("TEST_TODECIMAL_NEGATIVE_INPUT", func(t *testing.T) {
		// Arrange
		var input float64 = -5.048372
		var places int = 4

		var expectedResult float64 = -5.0484

		// Act
		res := ToDecimal(input, places)

		// Assert
		assert.Equal(t, expectedResult, res)
	})

	t.Run("TEST_TODECIMAL_NEGATIVE_PLACES", func(t *testing.T) {
		// Arrange
		var input float64 = 5.048372

		// if places < 0, they're set to 0, removing the decimal spaces
		var places int = -5

		var expectedResult float64 = 5

		// Act
		res := ToDecimal(input, places)

		// Assert
		assert.Equal(t, expectedResult, res)
	})

	t.Run("TEST_TODECIMAL_ZERO_PLACES", func(t *testing.T) {
		// Arrange
		var input float64 = 5.048372
		var places int = 0

		var expectedResult float64 = 5

		// Act
		res := ToDecimal(input, places)

		// Assert
		assert.Equal(t, expectedResult, res)
	})
}

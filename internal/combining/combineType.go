package combining

import "github.com/radoslavboychev/price-calculator-kata/config"

// Enum for the various discount combination types
const (
	TypeAdditive CombType = iota
	TypeMultiplicative
)

// CombineType defines an enum for discount combination types
type CombType uint16

// NewCombineType reads the defined discount combination type from the config
// and returns the proper combination type
// if an invalid value is set in the config, Additive method is used by default
func NewCombineTypeFromConfig() CombType {
	conf := config.LoadConfig()

	switch conf.CombinationType {
	case 0:
		return TypeAdditive
	case 1:
		return TypeMultiplicative
	default:
		return TypeAdditive
	}
}

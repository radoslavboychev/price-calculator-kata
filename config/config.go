package config

import (
	"log"

	"github.com/spf13/viper"
)

// config struct contains all configurable variables for the calculator application
type config struct {
	Tax                     uint16  `mapstructure:"TAX_RATE"`
	UniversalDiscountRate   uint16  `mapstructure:"UNIVERSAL_DISCOUNT_RATE"`
	SpecialDiscountRate     uint16  `mapstructure:"SPECIAL_DISCOUNT_RATE"`
	SpecialDiscountUPC      int     `mapstructure:"SPECIAL_DISCOUNT_UPC"`
	DiscountTakesPrecedence uint16  `mapstructure:"DISCOUNT_TAKES_PRECEDENCE"`
	CapType                 uint16  `mapstructure:"DISCOUNT_CAP_TYPE"`
	CapValue                float64 `mapstructure:"CAP_VALUE"`
	Currency                uint16  `mapstructure:"CURRENCY"`
	CombinationType         uint16  `mapstructure:"COMBINE_TYPE"`
	CostPercentage          float64 `mapstructure:"COST_PERCENTAGE"`
	CostAbsolute            float64 `mapstructure:"COST_ABSOLUTE"`
}

// variable to unmarshal the config in
var conf *config

// LoadConfig loads the environmental variables from a file
func LoadConfig() config {

	setDefaultConfigValues()

	viper.AddConfigPath("../.././config")

	viper.SetConfigName("config")

	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = viper.Unmarshal(&conf)
	if err != nil {
		log.Fatal(err)
	}

	return *conf
}

// setDefaultConfigValues defines some default values in the configuration in case they've not been set
func setDefaultConfigValues() {
	viper.SetDefault("TAX_RATE", 20)
	viper.SetDefault("UNIVERSAL_DISCOUNT_RATE", 0)
	viper.SetDefault("SPECIAL_DISCOUNT_UPC", 0)
	viper.SetDefault("SPECIAL_DISCOUNT_RATE", 0)
	viper.SetDefault("DISCOUNT_TAKES_PRECEDENCE", 0)
	viper.SetDefault("DISCOUNT_CAP_TYPE", 0)
	viper.SetDefault("CAP_VALUE", 0)
	viper.SetDefault("CURRENCY", 0)
	viper.SetDefault("COMBINE_TYPE", 0)
	viper.SetDefault("COST_PERCENTAGE", 0)
	viper.SetDefault("COST_PERCENTAGE", 0)
}

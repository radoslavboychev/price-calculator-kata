package main

import (
	"log"

	"github.com/radoslavboychev/price-calculator-kata/config"
	"github.com/radoslavboychev/price-calculator-kata/internal/cap"
	"github.com/radoslavboychev/price-calculator-kata/internal/combining"
	"github.com/radoslavboychev/price-calculator-kata/internal/currency"
	"github.com/radoslavboychev/price-calculator-kata/internal/models"
	"github.com/radoslavboychev/price-calculator-kata/pkg/calculator"
)

func main() {

	// CONFIGURATOR
	conf := config.LoadConfig()
	logConfig()

	// CURRENCY
	defaultCurrency := currency.LoadCurrency()

	// TAX
	tax := *models.NewTax(conf.Tax)

	// DISCOUNT
	universalDiscount := models.NewUniversalDiscount(conf.UniversalDiscountRate, models.NewMoney(defaultCurrency.Code, 0))
	specialDiscount := models.NewSpecialDiscount(conf.SpecialDiscountUPC, conf.SpecialDiscountRate, models.NewMoney(defaultCurrency.Code, 0))
	discount := *models.NewDiscount(*universalDiscount, *specialDiscount, models.TakesPrecedence(conf.DiscountTakesPrecedence))

	// EXPENSE
	expenseAbsolute := models.NewExpenseAbsolute("Transport", conf.CostAbsolute)
	expensePercentage := models.NewExpensePercentage("Packaging", conf.CostPercentage)
	productCosts := models.NewCosts(expenseAbsolute, expensePercentage)

	// COMBINING
	combineType := combining.NewCombineTypeFromConfig()

	// CAP
	discountCap := cap.NewDiscountCap(conf.CapValue)

	// create an object
	p := models.NewProduct("The Little Prince", 123456, models.NewMoney(defaultCurrency.Code, 20.25), productCosts)

	// create the calculator object
	calc := calculator.NewCalculator(tax, discount, combineType, discountCap)

	// conduct all calculations for the specific product
	res := calc.Calculate(&p)
	res.Report()
}

// logConfig prints the currently loaded config
func logConfig() {
	conf := config.LoadConfig()

	log.Println("CONFIGURATION")
	log.Println("###########")
	log.Printf("Tax Rate: %v\n", conf.Tax)
	log.Printf("Universal Discount Rate: %v \n", conf.UniversalDiscountRate)
	log.Printf("Special Discount: Rate - %v%%; UPC - %v \n", conf.SpecialDiscountRate, conf.SpecialDiscountUPC)

	switch conf.DiscountTakesPrecedence {
	case 1:
		log.Println("Universal discount takes precedence over tax!")
	case 2:
		log.Println("Special discount takes precedence over tax!")
	default:
		log.Println("Discount does not take precedence over tax!")
	}

	switch conf.CapType {
	case 1:
		log.Printf("Discount cap: Percentage-based, Value: %v%%\n", conf.CapValue)
	case 2:
		log.Printf("Discount cap: Absolute, Value: %v\n", conf.CapValue)
	default:
		log.Println("No discount cap has been set!")
	}

	log.Printf("Currency: %v\n", currency.LoadCurrency().Code)

	switch conf.CombinationType {
	case 0:
		log.Printf("Discount combination type: Additive!")
	case 1:
		log.Printf("Discount combination type: Multiplicative!")
	default:
		log.Printf("Invalid combination type")
	}

	log.Println("Executing calculations...")
	log.Println("###########")
}

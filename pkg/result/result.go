package result

import (
	"fmt"

	"github.com/radoslavboychev/price-calculator-kata/internal/models"
)

// Result stores calculator results
type Result struct {
	startingPrice models.Money
	taxAmount     models.Money
	totalDiscount models.Money
	totalExpenses models.Money
	totalPrice    models.Money
	costs         models.Costs
}

// NewResult constructor
func NewResult(startingPrice, taxAmount, totalDiscount, totalExpenses, totalPrice models.Money, costs models.Costs) *Result {
	return &Result{
		startingPrice: startingPrice,
		taxAmount:     taxAmount,
		totalDiscount: totalDiscount,
		totalExpenses: totalExpenses,
		totalPrice:    totalPrice,
		costs:         costs,
	}
}

// Report prints a report of all relevant results. Does not print amounts with null or zero values
func (r *Result) Report() string {

	// Starting price will be reported
	starting := fmt.Sprintf("Cost = %.2f %v\n", r.StartingPrice().Value, r.StartingPrice().Currency.String())
	fmt.Print(starting)

	// if the tax exists it will get reported
	var tax string
	if r.TaxAmount().Value != 0 {
		tax = fmt.Sprintf("Tax = %.2f %v\n", r.TaxAmount().Value, r.TaxAmount().Currency.String())
		fmt.Print(tax)
	}

	// if discounts exist they will be reported
	var totalDiscount string
	if r.TotalDiscount().Value != 0 {
		totalDiscount = fmt.Sprintf("Discounts = %.2f %v\n", r.TotalDiscount().Value, r.TotalDiscount().Currency.String())
		fmt.Print(totalDiscount)
	}

	// if expenses exist they will be reported one by one
	if r.TotalExpenses().Value != 0 {
		for _, c := range r.Costs().Expenses {
			if c != nil {
				c.ReportExpense(r.StartingPrice())
			}
		}
	}

	// the total price will be reported
	total := fmt.Sprintf("TOTAL = %.2f %v\n", r.TotalPrice().Value, r.TotalPrice().Currency.String())
	fmt.Print(total)

	// concatenate all strings and return them (for test cases)
	report := starting + tax + totalDiscount + total
	return report
}

// StartingPrice returns a result's starting price
func (r *Result) StartingPrice() models.Money {
	return r.startingPrice
}

// TaxAmount returns a result's tax amount
func (r *Result) TaxAmount() models.Money {
	return r.taxAmount
}

// TaxDiscount returns a result's tax amount
func (r *Result) TotalDiscount() models.Money {
	return r.totalDiscount
}

// TotalExpenses returns a result's total expenses
func (r *Result) TotalExpenses() models.Money {
	return r.totalExpenses
}

// TotalPrice returns a result's total price
func (r *Result) TotalPrice() models.Money {
	return r.totalPrice
}

// Costs returns a result's costs
func (r *Result) Costs() models.Costs {
	return r.costs
}

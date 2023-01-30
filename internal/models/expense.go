package models

import (
	"github.com/radoslavboychev/price-calculator-kata/internal/currency"

	"fmt"
)

// Expense interface defines behaviour for all Expense types that implement it
type expense interface {
	CalculateExpense(amount Money) Money
	ReportExpense(startingPrice Money)
}

// expensePercentage represents percentage-based expenses
type expensePercentage struct {
	Description string
	Amount      float64
	Currency    currency.CurrencyCode
}

// expenseAbsolute represents expenses with absolute values
type expenseAbsolute struct {
	Description string
	Amount      Money
}

// Costs represents all of the expenses for a product
type Costs struct {
	Expenses []expense
}

// NewCosts constructor that takes in any amount of product costs
func NewCosts(e ...expense) Costs {
	var costs = []expense{}
	costs = append(costs, e...)

	return Costs{
		Expenses: costs,
	}
}

// NewExpensePercentage constructor for percentage based expenses
func NewExpensePercentage(description string, amount float64) *expensePercentage {
	if amount < 0 {
		amount = 0
	}

	return &expensePercentage{
		Description: description,
		Amount:      amount,
	}
}

// NewExpenseAbsolute constructor for absolute value expenses
func NewExpenseAbsolute(description string, amount float64) *expenseAbsolute {
	if amount < 0 {
		amount = 0
	}

	return &expenseAbsolute{
		Description: description,
		Amount: Money{
			Value: amount,
		},
	}
}

// CalculateExpense calculates the exact amount of expense from a percentage
func (e *expensePercentage) CalculateExpense(startingPrice Money) Money {

	expenseAmount := (e.Amount / 100) * startingPrice.Value

	return Money{
		Value:    expenseAmount,
		Currency: startingPrice.Currency,
	}
}

// CalculateExpense for absolute value expenses returns the amount of expense for absolute amount expenses
func (e *expenseAbsolute) CalculateExpense(startingPrice Money) Money {
	return Money{
		Value:    e.Amount.Value,
		Currency: startingPrice.Currency,
	}
}

// CalculateExpense iterates through a list of costs, calculates their expenses and returns the sum of costs
func (e Costs) CalculateExpense(startingPrice Money) Money {
	var sum float64 = 0
	for _, v := range e.Expenses {
		sum += v.CalculateExpense(startingPrice).Value
	}
	return Money{
		Currency: startingPrice.Currency,
		Value:    sum,
	}
}

// ReportExpense iterates through all costs and reports them unless they're nil
func (e Costs) ReportExpense(startingPrice Money) {
	for _, v := range e.Expenses {
		v.ReportExpense(startingPrice)
	}
}

// ToString method to report the value of absolute expense costs
func (e *expenseAbsolute) ReportExpense(startingPrice Money) {

	if e.Amount.Value != 0 {
		str := fmt.Sprintf("%v =  %.2f %v", e.Description, e.Amount.Value, e.Amount.Currency)
		fmt.Println(str)
	}
}

// ToString method to calculate and report the value of percentage expense costs
func (e *expensePercentage) ReportExpense(startingPrice Money) {

	amount := e.CalculateExpense(startingPrice)
	if amount.Value != 0 {
		formattedAmount := fmt.Sprintf("%.2f %v", amount.Value, amount.Currency.String())
		str := fmt.Sprintf("%v = %v", e.Description, formattedAmount)

		fmt.Println(str)
	}

}

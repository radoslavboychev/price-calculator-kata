package models

import (
	"crypto/rand"
	"math"
	"math/big"
	"strconv"
)

// Product struct represents a product
type Product struct {
	name  string
	upc   int
	price Money
	cost  Costs
}

// NewProduct creates an instance of a new Product with the parameters set
func NewProduct(name string, upc int, price Money, cost Costs) Product {
	if price.Value < 0 {
		price.Value = 0
	}

	if name == "" {
		name = "Unnamed Product"
	}

	// if default UPC is too short, generate new one of length 6
	if len(strconv.Itoa(upc)) < 6 {
		upc, _ = generateUPC()
	}

	return Product{
		name:  name,
		upc:   upc,
		price: price,
		cost:  cost,
	}
}

// generateUPC creates a new, randomized UPC that is 6 digits long
func generateUPC() (int, error) {
	maxLimit := int64(int(math.Pow10(6)) - 1)
	lowLimit := int(math.Pow10(5))

	randomNumber, err := rand.Int(rand.Reader, big.NewInt(maxLimit))
	if err != nil {
		return 0, err
	}
	randomNumberInt := int(randomNumber.Int64())

	// Handling integers between 0, 10^(n-1) .. for n=4, handling cases between (0, 999)
	if randomNumberInt <= lowLimit {
		randomNumberInt += lowLimit
	}

	// Never likely to occur, kust for safe side.
	if randomNumberInt > int(maxLimit) {
		randomNumberInt = int(maxLimit)
	}
	return randomNumberInt, nil
}

// Returns product name
func (p Product) Name() string {
	return p.name
}

// Returns product UPC
func (p Product) UPC() int {
	return p.upc
}

// Returns product price
func (p Product) Price() Money {
	return p.price
}

// Returns product costs
func (p Product) Cost() Costs {
	return p.cost
}

package common

import (
	"github.com/quick-trade/xoney/common/data"
	"github.com/quick-trade/xoney/errors"
	"github.com/quick-trade/xoney/internal"
)

// BaseDistribution represents a user's balance on an exchange in the base currency.
// It maps each currency to the corresponding quantity, allowing for the representation
// of the entire portfolio in a unified manner.
type BaseDistribution map[data.Currency]float64

// Portfolio holds the financial position of a user on exchange
// and the main currency of valuation. It provides methods to manipulate and
// evaluate the asset distribution.
type Portfolio struct {
	// assets holds the distribution of different currencies and their quantities.
	assets BaseDistribution
	// mainCurrency is the primary currency used for valuation of the portfolio.
	mainCurrency data.Currency
}

// Calculates the total value of the portfolio in the main currency.
// It takes a map of prices where each currency is mapped to its current price in the main currency.
// The function returns the total as a float64 and an error if any currency in the portfolio doesn't have a corresponding price.
func (p Portfolio) Total(prices map[data.Currency]float64) (float64, error) {
	total := 0.0
	err := errors.NewMissingCurrencyError(internal.DefaultCapacity)
	success := true

	for currency, quantity := range p.assets {
		price, ok := prices[currency]
		if !ok {
			if currency.Asset == p.mainCurrency.Asset {
				price = 1
			} else {
				success = false
				err.Add(currency.String())
			}
		}

		total += quantity * price
	}

	if success {
		return total, nil
	}

	return total, err
}

// Balance returns the amount of the specified currency held in the portfolio.
func (p Portfolio) Balance(currency data.Currency) float64 {
	return p.assets[currency]
}

// Assets returns a reference to the current assets held in the portfolio.
// Be cautious when using this reference directly as it can alter the portfolio's state.
// Consider using .Copy() method for a safe, mutable copy if needed.
func (p Portfolio) Assets() BaseDistribution { return p.assets }

// Set assigns the specified quantity of a currency to the portfolio.
func (p *Portfolio) Set(currency data.Currency, quantity float64) {
	p.assets[currency] = quantity
}

// Increase adds the specified quantity of a currency to the portfolio.
func (p *Portfolio) Increase(currency data.Currency, quantity float64) {
	p.assets[currency] += quantity
}

// Decrease subtracts the specified quantity of a currency from the portfolio.
func (p *Portfolio) Decrease(currency data.Currency, quantity float64) {
	p.assets[currency] -= quantity
}

// Returns the main currency used for the valuation of the portfolio.
func (p Portfolio) MainCurrency() data.Currency { return p.mainCurrency }

// Creates a deep copy of the portfolio.
func (p Portfolio) Copy() Portfolio {
	return Portfolio{
		assets:       internal.MapCopy(p.assets),
		mainCurrency: p.mainCurrency,
	}
}

// NewPortfolio creates a new Portfolio with the specified main currency.
func NewPortfolio(mainCurrency data.Currency) Portfolio {
	return Portfolio{
		assets:       make(map[data.Currency]float64, internal.DefaultCapacity),
		mainCurrency: mainCurrency,
	}
}

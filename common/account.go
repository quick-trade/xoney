package common

import (
	"xoney/common/data"
	"xoney/errors"
	"xoney/internal"
)

type Portfolio struct {
	assets       map[data.Currency]float64
	mainCurrency data.Currency
}

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
				err.Add(currency.Asset)
			}
		}

		total += quantity * price
	}

	if success {
		return total, nil
	}

	return total, err
}

func (p Portfolio) Balance(currency data.Currency) float64 {
	return p.assets[currency]
}
func (p Portfolio) Assets() map[data.Currency]float64 { return p.assets }

func (p *Portfolio) Set(currency data.Currency, quantity float64) {
	p.assets[currency] = quantity
}

func (p *Portfolio) Increase(currency data.Currency, quantity float64) {
	p.assets[currency] += quantity
}

func (p *Portfolio) Decrease(currency data.Currency, quantity float64) {
	p.assets[currency] -= quantity
}
func (p Portfolio) MainCurrency() data.Currency { return p.mainCurrency }
func NewPortfolio(mainCurrency data.Currency) Portfolio {
	return Portfolio{
		assets:       make(map[data.Currency]float64, internal.DefaultCapacity),
		mainCurrency: mainCurrency,
	}
}

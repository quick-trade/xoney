package toolkit

import (
	"fmt"
	"math"
	"xoney/common"
	"xoney/common/data"
	"xoney/errors"
	"xoney/events"
	"xoney/exchange"
	"xoney/internal"

	st "xoney/strategy"
)

type BaseWeight float64
type PortfolioWeights map[data.Currency]BaseWeight
func NewPortfolioWeights(distribution map[data.Currency]BaseWeight) (*PortfolioWeights, error) {
	weights := PortfolioWeights(distribution)
	if err := weights.isValid(); err != nil {
		return nil, err
	}
	return &weights, nil
}
func (f PortfolioWeights) isValid() error {
	sumWeights := 0.0

	for _, weight := range f {
		sumWeights += math.Abs(float64(weight))
	}

	if sumWeights != 1 {
		return fmt.Errorf("invalid portfolio weights: sum of abs(weights): %f", sumWeights)
	}

	return nil
}

func (pw PortfolioWeights) Synchronize(
	current common.BaseDistribution,
	prices map[data.Currency]float64,
) (target common.BaseDistribution, err error) {
	totalQuote := 0.0
	totalQuoteWeight := 0.0
	missingCurrencyErr := errors.NewMissingCurrencyError(internal.DefaultCapacity)
	success := true

	// Calculate the total value of quote distribution
	for currency, amount := range current {
		price, ok := prices[currency]
		if !ok {
			missingCurrencyErr.Add(currency.String())
			success = false
		} else {
			totalQuote += amount * price
			totalQuoteWeight += float64(pw[currency]) * price
		}
	}

	if !success {
		return nil, missingCurrencyErr
	}

	// Calculate the target base distribution based on weights
	target = make(common.BaseDistribution)

	for currency, weight := range pw {
		target[currency] = float64(weight) * totalQuote / totalQuoteWeight
	}

	if !success {
		return nil, missingCurrencyErr
	}

	return target, nil
}


type CapitalAllocator interface {
	Start(charts data.ChartContainer) error
	Next(candle data.InstrumentCandle) (PortfolioWeights, error)
	MinDurations() st.Durations
}

type RebalancePortfolio struct {
	weights PortfolioWeights
}
func (r *RebalancePortfolio) Occur(connector exchange.Connector) error {
	//baseDistribution := connector.Portfolio().Assets()
	//var prices
	panic("TODO: Implement")
}
func NewRebalancePortfolio(weights PortfolioWeights) *RebalancePortfolio {
	return &RebalancePortfolio{weights: weights}
}

type CapitalAllocationBot struct {
	allocator CapitalAllocator
}
func (c *CapitalAllocationBot) MinDurations() st.Durations {
	return c.allocator.MinDurations()
}

func (c *CapitalAllocationBot) Start(charts data.ChartContainer) error {
	return c.allocator.Start(charts)
}

func (c *CapitalAllocationBot) Next(candle data.InstrumentCandle) ([]events.Event, error) {
	weights, err := c.allocator.Next(candle)
	if err != nil {
		return nil, err
	}

	return []events.Event{NewRebalancePortfolio(weights)}, nil
}

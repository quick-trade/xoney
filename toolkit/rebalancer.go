package toolkit

import (
	"fmt"
	"math"
	"xoney/common"
	"xoney/common/data"
	"xoney/events"
	"xoney/exchange"

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
func (f PortfolioWeights) Synchronize(
	current common.BaseDistribution,
	prices map[data.Currency]float64,
) (target common.BaseDistribution, err error) {
	totalQuote := 0.0

	// Calculate the total value of quote distribution
	for currency, amount := range current {
		price, ok := prices[currency]
		if !ok {
			return nil, fmt.Errorf("missing price for currency: %v", currency)
		}
		totalQuote += amount * price
	}

	// Calculate the target base distribution based on weights
	target = make(common.BaseDistribution)

	for currency, weight := range f {
		price, ok := prices[currency]
		if !ok {
			return nil, fmt.Errorf("missing price for currency: %v", currency)
		}

		target[currency] = float64(weight) * totalQuote / price
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

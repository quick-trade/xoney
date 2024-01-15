package toolkit

import (
	"xoney/common/data"
	"xoney/common"
	"xoney/events"
	"xoney/exchange"

	st "xoney/strategy"
)

type PortfolioWeight struct {
	quoteWeight float64
	buyPrice float64
}
type PortfolioWeights map[data.Currency]PortfolioWeight
func (f PortfolioWeights) isValid() bool {
	sumWeights := 0.0

	for _, weight := range f {
		sumWeights += weight.quoteWeight
	}

	return sumWeights <= 1
}
func (f PortfolioWeights) Synchronize(current common.BaseDistribution) error {
	panic("TODO: Implement")
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

package toolkit

import (
	"xoney/common/data"
	"xoney/events"
	"xoney/exchange"

	st "xoney/strategy"
)

type PortfolioWeights map[data.Currency]float64

type CapitalAllocator interface {
	Start(charts data.ChartContainer) error
	Next(candle data.InstrumentCandle) (PortfolioWeights, error)
	MinDurations() st.Durations
}

type RebalancePortfolio struct {
	weights PortfolioWeights
}
func (r *RebalancePortfolio) Occur(connector exchange.Connector) error {
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

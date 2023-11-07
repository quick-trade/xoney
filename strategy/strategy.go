package strategy

import (
	"time"
	"xoney/common/data"
	"xoney/events"
)

type Durations map[data.Instrument]time.Duration
func (d Durations) Max() time.Duration {
	var maxDur time.Duration
	for _, duration := range d {
		if duration > maxDur {
			maxDur = duration
		}
	}
	return maxDur
}

type Tradable interface {
	Start(charts data.ChartContainer)
	Next(candle data.Candle) []events.Event
	MinDurations() Durations
}

type VectorizedTradable interface {
	Tradable
	Backtest(
		commission float64,
		initialDepo float64,
		charts data.ChartContainer,
	) (data.Equity, error)
}

type Optimizable interface {
	Tradable
	Parameters() []Parameter
}

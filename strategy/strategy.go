package strategy

import (
	"time"

	"xoney/common/data"
	"xoney/events"
	"xoney/exchange"
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
	Start(charts data.ChartContainer) error
	Next(candle data.InstrumentCandle) (events.Event, error)
	MinDurations() Durations
}

type VectorizedTradable interface {
	Tradable
	Backtest(
		simulator exchange.Simulator,
		charts data.ChartContainer,
	) (data.Equity, error)
}

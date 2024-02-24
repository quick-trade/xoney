package strategy

import (
	"time"

	"github.com/quick-trade/xoney/common/data"
	"github.com/quick-trade/xoney/events"
	"github.com/quick-trade/xoney/exchange"
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

package strategy

import (
	"time"

	"xoney/common/data"
	"xoney/events"
)

type Tradable interface {
	Start(charts data.ChartContainer)
	Next(candle data.Candle) []events.Event
	MinDuration() time.Duration
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

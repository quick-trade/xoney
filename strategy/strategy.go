package strategy

import (
	"time"

	"xoney/common/data"
)

type Tradable interface {
	FetchEvents(charts data.ChartContainer)
	MinDuration() time.Duration
}

type VectorizedTradable interface {
	Tradable
	Backtest(commission float64) (data.Equity, error)
}

type Optimizable interface {
	Tradable
	Parameters() []Parameter
}

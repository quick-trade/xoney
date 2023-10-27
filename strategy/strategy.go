package strategy

import "xoney/common/data"

type Tradable interface {
	FetchEvents(charts data.ChartContainer)
	MinCandles() int
}

type VectorizedTradable interface {
	Tradable
	Backtest(commission float64)
}

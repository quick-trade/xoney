package optimize

import (
	bt "xoney/backtest"
	"xoney/common/data"
	"xoney/strategy"
)

type Optimizer interface {
	Optimize(system strategy.Optimizable, charts data.ChartContainer)
	GetBests(n int) []strategy.Optimizable
	SetMetric(metric bt.Metric)
}

type RandomOptimizer struct {
	backtester bt.Backtester
	trials     []strategy.Optimizable
}

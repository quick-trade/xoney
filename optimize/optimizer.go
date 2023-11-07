package optimize

import (
	bt "xoney/backtest"
	"xoney/common/data"
	st "xoney/strategy"
)

type Optimizer interface {
	Optimize(system st.Optimizable, charts data.ChartContainer) error
	GetBests(n int) []*st.Optimizable
	SetMetrics(metrics []bt.Metric)
	Metrics() []bt.Metric
}

type RandomOptimizer struct {
	backtester bt.Backtester
	trials     []*st.Optimizable
	metrics    []bt.Metric
}

func (r *RandomOptimizer) GetBests(n int) []*st.Optimizable {
	panic("TODO")
}

func (r *RandomOptimizer) Optimize(system *st.Optimizable, charts data.ChartContainer) error {
	panic("TODO")
}

func (r *RandomOptimizer) SetMetrics(metrics []bt.Metric) {
	r.metrics = metrics
}

func (r *RandomOptimizer) Metrics() []bt.Metric {
	return r.metrics
}

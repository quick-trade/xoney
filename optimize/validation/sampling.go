package validation

import (
	"fmt"
	"xoney/common"
	"xoney/common/data"

	bt "xoney/backtest"

	opt "xoney/optimize"
	st "xoney/strategy"
)

type EquityResult common.Result[data.Equity]

func newEquityResult(equity data.Equity, err error) EquityResult {
	return EquityResult{Data: equity, Error: err}
}

// important function for getting unbiased strategy estimations.
func shiftedCharts(
	charts data.ChartContainer,
	period common.Period,
	system *st.Optimizable,
) data.ChartContainer {
	period = period.ShiftedStart(-(*system).MinDuration())

	result := charts.ChartsByPeriod(period)

	return result
}

type InSample struct {
	optimizer opt.Optimizer
	period    common.Period
	charts    data.ChartContainer
}

func (i *InSample) Optimize(system *st.Optimizable) error {
	charts := shiftedCharts(i.charts, i.period, system)

	if err := i.optimizer.Optimize(system, charts); err != nil {
		return fmt.Errorf("failed to optimize the system: %w", err)
	}

	return nil
}

func (i *InSample) BestSystem() *st.Optimizable {
	return i.optimizer.GetBests(1)[0]
}

func NewInSample(
	charts data.ChartContainer,
	period common.Period,
	optimizer opt.Optimizer,
) *InSample {
	return &InSample{
		optimizer: optimizer,
		period:    period,
		charts:    charts,
	}
}

type OutOfSample struct {
	charts     data.ChartContainer
	backtester bt.Backtester
	period     common.Period
}

func (o *OutOfSample) Backtest(system *st.Optimizable) (data.Equity, error) {
	var tr st.Tradable = *system

	charts := shiftedCharts(o.charts, o.period, system)
	return o.backtester.Backtest(charts, &tr)
}

func NewOutOfSample(
	charts data.ChartContainer,
	period common.Period,
	backtester bt.Backtester,
) *OutOfSample {
	return &OutOfSample{
		charts:     charts,
		backtester: backtester,
		period:     period,
	}
}

type SamplePair struct {
	IS  *InSample
	OOS *OutOfSample
}

func (s *SamplePair) Test(system st.Optimizable) (data.Equity, error) {
	err := s.IS.Optimize(&system)
	if err != nil {
		return data.Equity{}, fmt.Errorf("error during optimization: %w", err)
	}

	best := s.IS.BestSystem()
	equity, err := s.OOS.Backtest(best)
	if err != nil {
		return data.Equity{}, fmt.Errorf("error during backtesting: %w", err)
	}

	return equity, nil
}

func (s *SamplePair) test(system st.Optimizable) EquityResult {
	equity, err := s.Test(system)

	return newEquityResult(equity, err)
}

type Sampler interface {
	Samples(data data.ChartContainer) ([]SamplePair, error)
}

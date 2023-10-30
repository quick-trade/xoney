package validation

import (
	"fmt"

	bt "xoney/backtest"
	"xoney/common"
	"xoney/common/data"
	opt "xoney/optimize"
	st "xoney/strategy"
)

type EquityResult common.Result[data.Equity]

func newEquityResult(equity data.Equity, err error) EquityResult {
	return EquityResult{Data: equity, Error: err}
}
func shiftedCharts(
	charts data.ChartContainer,
	period common.Period,
	system *st.Optimizable,
) (data.ChartContainer, error) {
	period = period.ShiftedStart(-(*system).MinDuration())

	result, err := charts.ChartsByPeriod(period)
	if err != nil {
		return data.ChartContainer{}, fmt.Errorf("failed to retrieve charts by period: %w", err)
	}

	return result, nil
}

type InSample struct {
	optimizer opt.Optimizer
	period    common.Period
	charts    data.ChartContainer
}

func (i *InSample) Optimize(system *st.Optimizable) error {
	charts, err := shiftedCharts(i.charts, i.period, system)
	if err != nil {
		return err
	}

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

func (o *OutOfSample) Backtest(system *st.Optimizable) data.Equity {
	var tr st.Tradable = *system

	return o.backtester.Backtest(o.charts, &tr)
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
		return data.Equity{}, err
	}

	best := s.IS.BestSystem()
	equity := s.OOS.Backtest(best)

	return equity, nil
}

func (s *SamplePair) test(system st.Optimizable) EquityResult {
	equity, err := s.Test(system)

	return newEquityResult(equity, err)
}

type Sampler interface {
	Samples(data data.ChartContainer) ([]SamplePair, error)
}

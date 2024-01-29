package backtest

import (
	"fmt"
	"time"

	"xoney/common/data"
	"xoney/exchange"
	"xoney/internal"
	exec "xoney/internal/executing"
	st "xoney/strategy"
)

type StepByStepBacktester struct {
	system    st.Tradable
	equity    data.Equity
	simulator exchange.Simulator
}

func NewStepByStepBacktester(simulator exchange.Simulator) *StepByStepBacktester {
	return &StepByStepBacktester{
		equity:    data.Equity{},
		simulator: simulator,
	}
}

func (b *StepByStepBacktester) Start(charts data.ChartContainer, system st.Tradable) error {
	err := b.setup(charts, system)
	if err != nil {
		return fmt.Errorf("error during backtest setup: %w", err)
	}

	return nil
}

func (b *StepByStepBacktester) Next(candle data.InstrumentCandle) error {
	if err := b.updatePrices(candle); err != nil {
		return err
	}

	timestamp := candle.TimeClose
	if err := b.updateBalance(timestamp); err != nil {
		return err
	}

	event, err := b.system.Next(candle)
	if err != nil {
		return err
	}

	if err = exec.ProcessEvent(b.simulator, event); err != nil {
		return err
	}

	return nil
}

func (b *StepByStepBacktester) GetEquity() data.Equity {
	return b.equity
}

func (b *StepByStepBacktester) setup(
	charts data.ChartContainer,
	system st.Tradable,
) error {
	err := b.cleanup()
	if err != nil {
		return err
	}

	b.system = system

	b.equity = *generateStartEquity(charts)

	durations := system.MinDurations()
	maxDuration := durations.Max()

	strategyCharts := lastByDuration(charts, maxDuration)
	err = system.Start(strategyCharts)

	return err
}

func (b *StepByStepBacktester) cleanup() error {
	err := b.simulator.Cleanup()
	if err != nil {
		return fmt.Errorf("failed to cleanup: %w", err)
	}

	return nil
}

func (b *StepByStepBacktester) updatePrices(candle data.InstrumentCandle) error {
	return b.simulator.UpdatePrice(candle)
}

func (b *StepByStepBacktester) updateBalance(timestamp time.Time) error {
	totalBalance, err := b.simulator.Total()
	if err != nil {
		return fmt.Errorf("error getting total balance: %w", err)
	}

	b.equity.AddValue(totalBalance, timestamp)

	b.equity.AddPortfolio(b.simulator.Portfolio().Assets())

	return nil
}

type Backtester struct {
	simulator exchange.Simulator
}

func NewBacktester(simulator exchange.Simulator) *Backtester {
	return &Backtester{
		simulator: simulator,
	}
}

func (b *Backtester) Backtest(
	charts data.ChartContainer,
	system st.Tradable,
) (data.Equity, error) {
	if vecTradable, ok := system.(st.VectorizedTradable); ok {
		return vecTradable.Backtest(b.simulator, charts)
	}

	equity, err := b.runTest(charts, system) // TODO: BUGFIX: charts here is not corrected by MinDurations
	if err != nil {
		return equity, fmt.Errorf("error during backtest: %w", err)
	}

	return data.Equity{}, nil
}

func (b *Backtester) runTest(
	charts data.ChartContainer,
	system st.Tradable,
) (data.Equity, error) {
	bt := NewStepByStepBacktester(b.simulator)

	startCharts := firstByDuration(charts, system.MinDurations().Max())
	bt.Start(startCharts, system)

	for _, candle := range charts.Candles() {
		bt.Next(candle)
	}

	return bt.GetEquity(), nil
}

func maxTimeFrame(charts data.ChartContainer) data.TimeFrame {
	var tf data.TimeFrame

	for _, chart := range charts {
		chartTimeframe := chart.Timestamp.Timeframe()
		if tf.Duration < chartTimeframe.Duration {
			tf = chartTimeframe
		}
	}

	return tf
}

func generateStartEquity(
	charts data.ChartContainer,
) *data.Equity {
	timeframe := maxTimeFrame(charts)

	return data.NewEquity(timeframe, internal.DefaultCapacity)
}

func firstByDuration(charts data.ChartContainer, maxDuration time.Duration) data.ChartContainer {
	start := charts.FirstStart()
	stop := start.Add(maxDuration)

	period := data.NewPeriod(start, stop)

	return charts.ChartsByPeriod(period)
}

func lastByDuration(
	charts data.ChartContainer,
	maxDuration time.Duration,
) data.ChartContainer {
	end := charts.LastEnd()
	start := end.Add(-maxDuration)

	return charts.ChartsByPeriod(data.NewPeriod(start, end))
}

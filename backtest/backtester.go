package backtest

import (
	"fmt"
	"time"
	"xoney/common/data"
	exec "xoney/internal/executing"
	"xoney/exchange"

	st "xoney/strategy"
)

type Backtester struct {
	equity    data.Equity
	simulator exchange.Simulator
}

func NewBacktester(simulator exchange.Simulator) *Backtester {
	return &Backtester{
		equity:    data.Equity{},
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

	err := b.setup(charts, system)
	if err != nil {
		return b.equity, fmt.Errorf("error during backtest setup: %w", err)
	}

	err = b.runTest(charts, system) // TODO: BUGFIX: charts here is not corrected by MinDurations
	if err != nil {
		return b.equity, fmt.Errorf("error during backtest: %w", err)
	}

	return b.equity, nil
}

func (b *Backtester) setup(
	charts data.ChartContainer,
	system st.Tradable,
) error {
	err := b.cleanup()
	if err != nil {
		return err
	}

	durations := system.MinDurations()
	maxDuration := durations.Max()
	period := equityPeriod(charts, durations)

	b.equity = *generateEquity(charts, period, maxDuration)

	startPeriod := setupPeriod(charts, maxDuration)
	strategyCharts := charts.ChartsByPeriod(startPeriod)
	err = system.Start(strategyCharts)

	return err
}

func (b *Backtester) runTest(
	charts data.ChartContainer,
	system st.Tradable,
) error {
	for _, candle := range charts.Candles() {
		if err := b.updatePrices(candle); err != nil {
			return err
		}

		timestamp := candle.TimeClose
		if err := b.updateBalance(timestamp); err != nil {
			return err
		}

		event, err := system.Next(candle)
		if err != nil {
			return err
		}

		if err = exec.ProcessEvent(b.simulator, event); err != nil {
			return err
		}
	}

	return nil
}

func (b *Backtester) cleanup() error {
	err := b.simulator.Cleanup()
	if err != nil {
		return fmt.Errorf("failed to cleanup: %w", err)
	}

	return nil
}

func (b *Backtester) updatePrices(candle data.InstrumentCandle) error {
	return b.simulator.UpdatePrice(candle)
}

func (b *Backtester) updateBalance(timestamp time.Time) error {
	totalBalance, err := b.simulator.Total()
	if err != nil {
		return fmt.Errorf("error getting total balance: %w", err)
	}

	b.equity.AddValue(totalBalance, timestamp)

	b.equity.AddPortfolio(b.simulator.Portfolio().Assets())

	return nil
}

func equityPeriod(
	charts data.ChartContainer,
	durations st.Durations,
) data.Period {
	var firstStart time.Time

	var latestEnd time.Time

	for inst, chart := range charts {
		chartStart := chart.Timestamp.Start()
		instMinDuration := durations[inst]
		instStart := chartStart.Add(instMinDuration)

		if firstStart.Before(instStart) {
			firstStart = chartStart
		}

		chartEnd := chart.Timestamp.End()

		if latestEnd.Before(chartEnd) {
			latestEnd = chartEnd
		}
	}

	return data.NewPeriod(firstStart, latestEnd)
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

func generateEquity(
	charts data.ChartContainer,
	period data.Period,
	maxDuration time.Duration,
) *data.Equity {
	period = period.ShiftedStart(-maxDuration)

	timeframe := maxTimeFrame(charts)
	duration := period.End.Sub(period.Start)
	length := int(duration/timeframe.Duration) + 1

	return data.NewEquity(timeframe, period.Start, length)
}

func setupPeriod(charts data.ChartContainer, maxDuration time.Duration) data.Period {
	start := charts.FirstStart()
	stop := start.Add(maxDuration)

	return data.NewPeriod(start, stop)
}

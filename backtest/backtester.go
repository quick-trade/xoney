package backtest

import (
	"fmt"
	"time"
	"xoney/common/data"
	"xoney/events"
	"xoney/exchange"
	st "xoney/strategy"
)

type Backtester struct {
	initialDepo float64
	equity      data.Equity
	simulator   exchange.Simulator
}

func NewBacktester(initialDepo float64, currency data.Currency) *Backtester {
	return &Backtester{
		initialDepo: initialDepo,
		equity:      data.Equity{},
		simulator:   exchange.NewSimulator(currency, initialDepo),
	}
}

func (b *Backtester) Backtest(
	charts data.ChartContainer,
	system st.Tradable,
) (data.Equity, error) {
	if vecTradable, ok := system.(st.VectorizedTradable); ok {
		return vecTradable.Backtest(b.initialDepo, charts)
	}

	err := b.setup(charts, system)
	if err != nil {
		return b.equity, fmt.Errorf("error during backtest setup: %w", err)
	}

	err = b.runTest(charts, system)
	if err != nil {
		return b.equity, fmt.Errorf("error during backtest: %w", err)
	}

	return b.equity, nil
}

func (b *Backtester) setup(
	charts data.ChartContainer,
	system st.Tradable,
) error {
	// TODO: add cleaning up an exchange
	durations := system.MinDurations()
	period := equityPeriod(charts, durations)

	b.equity = *generateEquity(charts, period, durations.Max())

	err := system.Start(charts.ChartsByPeriod(period))

	return err
}

func (b *Backtester) runTest(
	charts data.ChartContainer,
	system st.Tradable,
) error {
	start := b.equity.Start()
	timeframe := b.equity.Timeframe().Duration
	nextTime := start.Add(timeframe)

	for _, candle := range charts.Candles() {
		if err := b.updatePrices(candle); err != nil {
			return err
		}

		if candle.TimeClose.After(nextTime) {
			if err := b.processBalance(); err != nil {
				return err
			}

			nextTime = nextTime.Add(timeframe)
		}

		events, err := system.Next(candle)
		if err != nil {
			return err
		}

		b.processEvents(events)
	}

	return nil
}

func (b *Backtester) updatePrices(candle data.InstrumentCandle) error {
	return b.simulator.UpdatePrice(candle)
}

func (b *Backtester) processBalance() error {
	totalBalance, err := b.simulator.Total()
	if err != nil {
		return err
	}

	b.equity.AddValue(totalBalance)

	return nil
}

func (b *Backtester) processEvents(events []events.Event) {
	for _, e := range events {
		// TODO: handle errors
		e.Occur(&b.simulator)
	}
}

func equityPeriod(
	charts data.ChartContainer,
	durations st.Durations,
) data.Period {
	var firstStart time.Time
	var chartStart time.Time
	var instMinDuration time.Duration
	var instStart time.Time

	var latestEnd time.Time
	var chartEnd time.Time

	for inst, chart := range charts {
		chartStart = chart.Timestamp.Start()
		instMinDuration = durations[inst]
		instStart = chartStart.Add(instMinDuration)

		if firstStart.Before(instStart) {
			firstStart = chartStart
		}

		chartEnd = chart.Timestamp.End()

		if latestEnd.Before(chartEnd) {
			latestEnd = chartEnd
		}
	}

	return data.Period{firstStart, latestEnd}
}

func maxTimeFrame(charts data.ChartContainer) data.TimeFrame {
	var tf data.TimeFrame
	var chartTimeframe data.TimeFrame

	for _, chart := range charts {
		chartTimeframe = chart.Timestamp.Timeframe()
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
	duration := period[1].Sub(period[0])
	length := int(duration/timeframe.Duration) + 1

	return data.NewEquity(length, timeframe, period[0])
}

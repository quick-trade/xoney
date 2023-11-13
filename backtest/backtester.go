package backtest

import (
	"fmt"
	"time"

	"xoney/common/data"
	"xoney/events"
	"xoney/internal"
	st "xoney/strategy"
	"xoney/trade"
)

type Backtester struct {
	trades      trade.TradeHeap
	commission  float64
	initialDepo float64
	equity      data.Equity
}

func NewBacktester(commission float64, initialDepo float64) *Backtester {
	return &Backtester{
		trades:      trade.NewTradeHeap(internal.DefaultCapacity),
		commission:  commission,
		initialDepo: initialDepo,
		equity:      data.Equity{},
	}
}

func (b *Backtester) Backtest(
	charts data.ChartContainer,
	system st.Tradable,
) (data.Equity, error) {
	if vecTradable, ok := system.(st.VectorizedTradable); ok {
		return vecTradable.Backtest(b.commission, b.initialDepo, charts)
	}

	err := b.setup(charts, &system)
	if err != nil {
		return b.equity, fmt.Errorf("error during backtest setup: %w", err)
	}

	err = b.runTest(charts, &system)
	if err != nil {
		return b.equity, fmt.Errorf("error during backtest: %w", err)
	}

	return b.equity, nil
}

func (b *Backtester) setup(
	charts data.ChartContainer,
	system *st.Tradable,
) error {
	b.clearTrades()

	durations := (*system).MinDurations()
	period := equityPeriod(charts, durations)

	b.equity = *generateEquity(charts, period, durations, b.initialDepo)

	err := (*system).Start(charts.ChartsByPeriod(period))

	return err
}

func (b *Backtester) runTest(
	charts data.ChartContainer,
	system *st.Tradable,
) error {
	// equityTime := b.equity.Timestamp.Start()
	// nextTime := equityTime.Add(b.equity.Timeframe().Duration)
	for _, candle := range charts.Candles() {
		events, err := (*system).Next(candle)
		if err != nil {
			return err
		}

		b.processEvents(events)
	}
	// TODO: should write an equity
	return nil
}

func (b *Backtester) clearTrades() {
	b.trades = trade.NewTradeHeap(internal.DefaultCapacity)
}

func (b *Backtester) processEvents(events []events.Event) {
	for _, e := range events {
		e.HandleTrades(&b.trades)
	}
}

func equityPeriod(
	charts data.ChartContainer,
	durations st.Durations,
) data.Period {
	var start time.Time
	var chartStart time.Time
	var instMinDuration time.Duration
	var instStart time.Time

	var latestEnd time.Time
	var chartEnd time.Time

	for inst, chart := range charts {
		chartStart = chart.Timestamp.Start()
		instMinDuration = durations[inst]
		instStart = chartStart.Add(instMinDuration)

		if start.Before(instStart) {
			start = chartStart
		}

		chartEnd = chart.Timestamp.End()

		if latestEnd.Before(chartEnd) {
			latestEnd = chartEnd
		}
	}

	return data.Period{start, latestEnd}
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
	durations st.Durations,
	initial float64,
) *data.Equity {
	timeframe := maxTimeFrame(charts)
	duration := period[1].Sub(period[0])
	length := int(duration/timeframe.Duration) + 2

	return data.NewEquity(length, timeframe, period[0], initial)
}

package backtest

import (
	"fmt"
	"time"
	"xoney/common"
	"xoney/common/data"
	"xoney/events"
	"xoney/internal"
	"xoney/trade"

	st "xoney/strategy"
)

type Backtester struct {
	trades      trade.TradeHeap
	commission  float64
	initialDepo float64
	equity      data.Equity
	portfolio   common.Portfolio
	prices      map[data.Currency]float64
}

func NewBacktester(commission float64, initialDepo float64) *Backtester {
	return &Backtester{
		trades:      trade.NewTradeHeap(internal.DefaultCapacity),
		commission:  commission,
		initialDepo: initialDepo,
		equity:      data.Equity{},
		portfolio:   common.NewPortfolio(internal.DefaultCapacity),
		prices:      make(map[data.Currency]float64, internal.DefaultCapacity),
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

	b.equity = *generateEquity(charts, period, durations.Max(), b.initialDepo)

	err := (*system).Start(charts.ChartsByPeriod(period))

	return err
}

func (b *Backtester) runTest(
	charts data.ChartContainer,
	system *st.Tradable,
) error {
	start := b.equity.Timestamp.Start()
	timeframe := b.equity.Timeframe().Duration
	nextTime := start.Add(timeframe)

	clear(b.prices)

	for _, candle := range charts.Candles() {
		b.updatePrices(candle)
		if candle.TimeClose.After(nextTime) {
			if err := b.processBalance(); err != nil {
				return err
			}

			nextTime = nextTime.Add(timeframe)
		}

		events, err := (*system).Next(candle)
		if err != nil {
			return err
		}

		b.processEvents(events)
	}

	return nil
}

func (b *Backtester) updatePrices(candle data.InstrumentCandle) {
	b.prices[candle.Instrument.Symbol().Base()] = candle.Close
}

func (b *Backtester) clearTrades() {
	b.trades = trade.NewTradeHeap(internal.DefaultCapacity)
}

func (b *Backtester) processBalance() error {
	totalBalance, err := b.portfolio.Total(b.prices)
	if err != nil {
		return err
	}

	b.equity.AddValue(totalBalance)

	return nil
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
	initial float64,
) *data.Equity {
	period = period.ShiftedStart(-maxDuration)
	
	timeframe := maxTimeFrame(charts)
	duration := period[1].Sub(period[0])
	length := int(duration/timeframe.Duration) + 1

	return data.NewEquity(length, timeframe, period[0], initial)
}

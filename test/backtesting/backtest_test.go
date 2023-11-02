package backtesting_test

import (
	"os"
	"testing"
	"time"

	bt "xoney/backtest"
	"xoney/common/data"
	st "xoney/strategy"
	testdata "xoney/testdata/backtesting"
	dtr "xoney/testdata/dataread"
)

var (
	charts     data.ChartContainer
	instrument data.Instrument
)

func TestMain(m *testing.M) {
	// Uploading chart data once
	chart, err := dtr.LoadChartFromCSV("../../testdata/BTCUSDT15m.csv")
	if err != nil {
		panic(err)
	}
	charts = make(data.ChartContainer, 1)

	tf, err := data.NewTimeFrame(time.Minute*15, "15m")
	if err != nil {
		panic(err)
	}
	sym, err := data.NewSymbol("BTC", "USD", "BINANCE")
	if err != nil {
		panic(err)
	}
	instrument = data.NewInstrument(*sym, *tf)
	charts[instrument] = chart

	// Running all the test
	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestBacktestReturnsEquity(t *testing.T) {
	bt := bt.NewBacktester(0, 1)
	var system st.Tradable = testdata.NewBBStrategy(1000, 1.9, instrument)
	equity, err := bt.Backtest(charts, &system)
	if err != nil {
		t.Error(err.Error())
	}
	err = dtr.WriteSlice(equity.Deposit(), "Equity", "../../testdata/BBEquity.csv")
	if err != nil {
		t.Error(err.Error())

	}
}

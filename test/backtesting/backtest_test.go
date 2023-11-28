package backtesting_test

import (
	"os"
	"testing"
	"time"

	bt "xoney/backtest"
	"xoney/common"
	"xoney/common/data"
	"xoney/exchange"
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
	timeframe, err := data.NewTimeFrame(time.Minute*15, "15m")
	if err != nil {
		panic(err)
	}

	chart, err := dtr.LoadChartFromCSV("../../testdata/BTCUSDT15m.csv", *timeframe)
	if err != nil {
		panic(err)
	}

	charts = make(data.ChartContainer, 1)

	sym, err := data.NewSymbol("BTC", "USD", "BINANCE")
	if err != nil {
		panic(err)
	}

	instrument = data.NewInstrument(*sym, *timeframe)
	charts[instrument] = chart

	// Running all the tests
	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestBacktestReturnsEquity(t *testing.T) {
	currency := data.NewCurrency("USD", "BINANCE")
	portfolio := common.NewPortfolio(currency)
	portfolio.Set(currency, 20000)

	simulator := exchange.NewSimulator(portfolio)
	tester := bt.NewBacktester(simulator)

	var system st.Tradable = testdata.NewBBStrategy(300, 1.5, instrument)

	equity, err := tester.Backtest(charts, system)
	if err != nil {
		t.Error(err.Error())
	}

	history := equity.Deposit()
	balanceHistory := equity.PortfolioHistory()
	balanceHistory[data.NewCurrency("Total", "")] = history


	err = dtr.WriteMap(balanceHistory, "../../testdata/BBEquity.csv")
	if err != nil {
		t.Error(err.Error())
	}
}

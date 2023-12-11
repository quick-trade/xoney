package backtesting_test

import (
	"os"
	"testing"
	"time"
	"xoney/common"
	"xoney/common/data"
	"xoney/exchange"

	bt "xoney/backtest"

	testdata "xoney/testdata/backtesting"
	dtr "xoney/testdata/dataread"
)

var (
	btc15m data.Instrument
	charts data.ChartContainer
)

func btc15min() data.Instrument {
	btcUsd := data.NewSymbol("BTC", "USD", "BINANCE")
	m15, _ := data.NewTimeFrame(time.Minute*15, "15m")

	return data.NewInstrument(*btcUsd, *m15)
}

func getCharts() data.ChartContainer {
	m15 := btc15m.Timeframe()
	btc, err := dtr.LoadChartFromCSV("../../testdata/BTCUSDT15m.csv", m15)
	if err != nil {
		panic(err)
	}

	charts := make(data.ChartContainer, 1)

	charts[btc15m] = btc

	return charts
}

func btcStrategy() testdata.BBBStrategy {
	return *testdata.NewBBStrategy(300, 2, btc15m)
}

func TestMain(m *testing.M) {
	// Uploading chart data once
	btc15m = btc15min()
	charts = getCharts()

	// Running all the tests
	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestBacktestReturnsEquity(t *testing.T) {
	currency := data.NewCurrency("USD", "BINANCE")
	portfolio := common.NewPortfolio(currency)
	portfolio.Set(currency, 17100)

	simulator := exchange.NewMarginSimulator(portfolio)
	tester := bt.NewBacktester(&simulator)

	system := btcStrategy()

	equity, err := tester.Backtest(charts, &system)
	if err != nil {
		t.Error(err.Error())
	}

	history := equity.Deposit()
	balanceHistory := equity.PortfolioHistory()
	balanceHistory[data.NewCurrency("Total", "")] = history
	balanceHistory[data.NewCurrency("mean", "")] = system.Mean
	balanceHistory[data.NewCurrency("UB", "")] = system.UB
	balanceHistory[data.NewCurrency("LB", "")] = system.LB

	err = dtr.WriteMap(balanceHistory, "../../testdata/BBEquity.csv")
	if err != nil {
		t.Error(err.Error())
	}
}

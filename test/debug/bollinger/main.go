package main

import (
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
	btc15m data.Instrument
	charts data.ChartContainer
)

func btc15min() data.Instrument {
	btcUsd := data.NewSymbol("BTC", "USD", "BINANCE")
	m1, _ := data.NewTimeFrame(time.Minute, "1m")

	return data.NewInstrument(*btcUsd, *m1)
}

func getCharts() data.ChartContainer {
	m1 := btc15m.Timeframe()
	btc, err := dtr.LoadChartFromCSV("testdata/BTCUSDT1m.csv", m1, 1)
	if err != nil {
		panic(err)
	}

	charts := make(data.ChartContainer, 1)

	charts[btc15m] = btc

	return charts
}

func btcBBStrategy() testdata.BBBStrategy {
	return *testdata.NewBBStrategy(60*24*30, 2, btc15m)
}

func portfolio() common.Portfolio {
	currency := data.NewCurrency("USD", "BINANCE")
	portfolio := common.NewPortfolio(currency)
	portfolio.Set(currency, 4300)

	return portfolio
}

func backtester() bt.Backtester {
	portfolio := portfolio()

	simulator := exchange.NewMarginSimulator(portfolio)
	tester := bt.NewBacktester(&simulator)

	return *tester
}

func backtest(system st.Tradable) data.Equity {
	tester := backtester()

	equity, err := tester.Backtest(charts, system)
	if err != nil {
		panic(err)
	}
	return equity
}

func debugBollinger() {
	system := btcBBStrategy()
	equity := backtest(&system)

	history := equity.Deposit()
	balanceHistory := equity.PortfolioHistory()
	balanceHistory[data.NewCurrency("Total", "")] = history
	balanceHistory[data.NewCurrency("mean", "")] = system.Mean
	balanceHistory[data.NewCurrency("UB", "")] = system.UB
	balanceHistory[data.NewCurrency("LB", "")] = system.LB

	err := dtr.WriteMap(balanceHistory, "testdata/BBEquity.csv")
	if err != nil {
		panic(err)
	}
}

func main() {
	// Uploading chart data once
	btc15m = btc15min()
	charts = getCharts()

	debugBollinger()
}

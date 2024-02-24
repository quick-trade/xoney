package main

import (
	"time"

	bt "github.com/quick-trade/xoney/backtest"
	"github.com/quick-trade/xoney/common"
	"github.com/quick-trade/xoney/common/data"
	"github.com/quick-trade/xoney/exchange"
	st "github.com/quick-trade/xoney/strategy"
	testdata "github.com/quick-trade/xoney/testdata/backtesting"
	dtr "github.com/quick-trade/xoney/testdata/dataread"
	tk "github.com/quick-trade/xoney/toolkit"
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
	btc, err := dtr.LoadChartFromCSV("testdata/BTCUSDT15m.csv", m15, 0)
	if err != nil {
		panic(err)
	}

	charts := make(data.ChartContainer, 1)

	charts[btc15m] = btc

	return charts
}

func portfolio() common.Portfolio {
	currency := data.NewCurrency("USD", "BINANCE")
	portfolio := common.NewPortfolio(currency)
	portfolio.Set(currency, 20000)

	return portfolio
}

func backtester() bt.Backtester {
	portfolio := portfolio()

	simulator := exchange.NewMarginSimulator(portfolio, 0.001)
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

func gridBot() *tk.GridBot {
	generator := testdata.NewAutoGrid(100, 2, 1.5, 0.5)

	return tk.NewGridBot(generator, btc15m)
}

func debugGrid() {
	bot := gridBot()
	equity := backtest(bot)

	err := dtr.WriteMap(equity.PortfolioHistory(), "testdata/BBEquity.csv")
	if err != nil {
		panic(err)
	}
}

func main() {
	// Uploading chart data once
	btc15m = btc15min()
	charts = getCharts()

	debugGrid()
}

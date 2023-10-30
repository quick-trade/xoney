package backtesting_test

import (
	"testing"
	bt "xoney/backtest"
)

func TestBacktestReturnsEquity(t *testing.T) {
	bt := bt.NewBacktester(0, 1)
	bt.Backtest()
}

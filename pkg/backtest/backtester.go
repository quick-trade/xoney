package backtest

import (
	"xoney/internal"
	"xoney/pkg/common"
)

type Backtester struct {
	trades      internal.TradeHeap
	commission  float64
	initialDepo float64
}

func (b *Backtester) Backtest(
	charts common.ChartContainer,
	system common.Tradable,
	independent_testing bool,
) common.Equity {

}

func NewBacktester(commission float64, initial_depo float64) *Backtester {
	return &Backtester{
		trades:      internal.TradeHeap{},
		commission:  commission,
		initialDepo: initial_depo,
	}
}

package backtest

import (
	"xoney/common/data"
	st "xoney/strategy"
	"xoney/trade"
)

type Backtester struct {
	trades      trade.TradeHeap
	commission  float64
	initialDepo float64
}

func (b *Backtester) Backtest(
	charts data.ChartContainer,
	system *st.Tradable,
) (data.Equity, error) {
	if vecTradable, ok := (*system).(st.VectorizedTradable); ok {
		return vecTradable.Backtest(b.commission, b.initialDepo, charts)
	}

	panic("TODO: Implement")
}

func NewBacktester(commission float64, initialDepo float64) *Backtester {
	return &Backtester{
		trades:      *trade.NewTradeHeap(),
		commission:  commission,
		initialDepo: initialDepo,
	}
}

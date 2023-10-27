package backtest

import "xoney/common/data"

type Metric interface {
	IsPositive() bool
	Evaluate(equity data.Equity)
}

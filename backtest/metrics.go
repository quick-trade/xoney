package backtest

import (
	"math"
	"xoney/common/data"
	"xoney/internal"
)

type Metric interface {
	IsPositive() bool
	Evaluate(equity data.Equity) float64
}

type SharpeRatio struct {
	Rf float64
}

func (s SharpeRatio) IsPositive() bool { return true }
func (s SharpeRatio) Evaluate(equity data.Equity) float64 {
	deposit := equity.Deposit()
	mean, err := internal.RawMoment(deposit, 1)
	if err != nil {
		return 0
	}

	variance := internal.CentralMoment(deposit, mean, 2)
	std := math.Sqrt(variance)

	return (mean - s.Rf) / std
}

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

type CARA struct {
	Theta float64
}

func (c CARA) IsPositive() bool { return true }
func (c CARA) Evaluate(equity data.Equity) float64 {
	deposit := equity.Deposit()
	mean, err := internal.RawMoment(deposit, 1)
	if err != nil {
		return 0
	}

	variance := internal.CentralMoment(deposit, mean, 2)
	CentralMoment3 := internal.CentralMoment(deposit, mean, 3)
	CentralMoment4 := internal.CentralMoment(deposit, mean, 4)

	return mean -
		c.Theta*variance/2 +
		math.Pow(c.Theta, 2)*CentralMoment3/6 -
		c.Theta*CentralMoment4/720
}

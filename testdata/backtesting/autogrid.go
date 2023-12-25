package backtesting

import (
	"math"
	"time"
	"xoney/common/data"
	"xoney/internal"
	tk "xoney/toolkit"
)

type bounds struct {
	minimum float64
	maximum float64
}
func (b bounds) Contains(price float64) bool {
	return b.minimum <= price && b.maximum >= price
}
type AutoGrid struct {
	candles int
	nLevels int
	deviations float64
	allAmount float64
	priceBounds bounds
	prices []float64
	levels []tk.GridLevel
}

func (a *AutoGrid) MinDuration(timeframe data.TimeFrame) time.Duration {
	return timeframe.Duration * time.Duration(a.candles)
}

func (a *AutoGrid) Next(candle data.Candle) ([]tk.GridLevel, error) {
	a.prices = append(a.prices[1:], candle.Close)

	if a.priceBounds.Contains(candle.Close) {
		return nil, nil
	}

	a.updateBounds()

	return a.levels, nil
}
func (a *AutoGrid) updateBounds() {
	mean, _ := internal.RawMoment(a.prices, 1)
	variance := internal.CentralMoment(a.prices, mean, 2)
	std := math.Sqrt(variance)

	minPrice := mean - std * a.deviations
	maxPrice := mean + std * a.deviations

	a.priceBounds = bounds{minimum: minPrice, maximum: maxPrice}

	a.generateGrid()
}
func (a *AutoGrid) generateGrid() {
	clear(a.levels)

	prices := linspace(a.priceBounds, a.nLevels)

	levelAmount := a.allAmount / float64(a.nLevels)

	for _, price := range prices {
		level := tk.NewGridLevel(price, levelAmount)
		// amount=1 means that amount is weighted equally
		a.levels = internal.Append(a.levels, *level)
	}
}

func (a *AutoGrid) Start(chart data.Chart) error {
	a.prices = chart.Close

	a.updateBounds()

	return nil
}

func linspace(b bounds, n int) []float64 {
	step := (b.maximum - b.minimum) / float64(n-1)
	values := make([]float64, n)

	for i := 0; i < n; i++ {
		values[i] = b.minimum + float64(i)*step
	}

	return values
}

func NewAutoGrid(candles, nLevels int, deviations, sumAmount float64) *AutoGrid {
	return &AutoGrid{
		candles: candles,
		nLevels: nLevels,
		deviations: deviations,
		allAmount: sumAmount,
		priceBounds: bounds{0, 0},
		prices: nil,
		levels: nil,
	}
}

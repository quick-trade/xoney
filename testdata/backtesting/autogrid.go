package backtesting

import (
	"fmt"
	"math"
	"time"

	"github.com/quick-trade/xoney/common/data"
	"github.com/quick-trade/xoney/internal"
	tk "github.com/quick-trade/xoney/toolkit"
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

	if err := a.updateBounds(); err != nil {
		return nil, fmt.Errorf("failed to update bounds: %w", err)
	}

	return a.levels, nil
}
func (a *AutoGrid) updateBounds() error {
	mean, _ := internal.RawMoment(a.prices, 1)
	variance := internal.CentralMoment(a.prices, mean, 2)
	std := math.Sqrt(variance)

	minPrice := mean - std * a.deviations
	maxPrice := mean + std * a.deviations

	a.priceBounds = bounds{minimum: minPrice, maximum: maxPrice}

	if err := a.generateGrid(); err != nil {
		return fmt.Errorf("failed to generate grid: %w", err)
	}

	return nil
}
func (a *AutoGrid) generateGrid() error {
	a.levels = make([]tk.GridLevel, 0, len(a.levels))

	prices := linspace(a.priceBounds, a.nLevels)

	levelAmount := a.allAmount / float64(a.nLevels)

	var firstErr error
	for _, price := range prices {
		level, err := tk.NewGridLevel(price, levelAmount)
		if firstErr == nil && err != nil {
			firstErr = err
		}

		a.levels = internal.Append(a.levels, *level)
	}

	if firstErr != nil {
		return fmt.Errorf("failed to generate one or more levels: %w", firstErr)
	}

	return nil
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

// Package data provides types and functions for managing financial equity data.

package data

import (
	"time"

	"github.com/quick-trade/xoney/internal"
)

// Equity represents a financial entity with history of its value and portfolio.
// It tracks the historical and current value of a portfolio along with a timestamp
// for each recorded value.
type Equity struct {
	portfolioHistory []map[Currency]float64 // history of portfolio values by currency
	mainHistory      []float64              // history of main value changes
	Timestamp        TimeStamp              // timestamps corresponding to mainHistory records
	timeframe        TimeFrame              // timeframe for the historical data
}

// Timeframe returns the timeframe of the equity's historical data.
func (e *Equity) Timeframe() TimeFrame {
	return e.timeframe
}

// Deposit returns the history of main value changes.
func (e *Equity) Deposit() []float64 { return e.mainHistory }

// PortfolioHistory constructs and returns a history of portfolio values by currency.
// It maps each currency to a slice of its historical values.
func (e *Equity) PortfolioHistory() map[Currency][]float64 {
	if len(e.portfolioHistory) == 0 {
		return make(map[Currency][]float64)
	}

	last := e.portfolioHistory[len(e.portfolioHistory)-1]
	result := make(map[Currency][]float64, len(last))

	for currency := range last {
		result[currency] = make([]float64, 0, len(e.portfolioHistory))
	}

	for i := range e.portfolioHistory {
		for currency := range last {
			value := e.portfolioHistory[i][currency]
			result[currency] = internal.Append(result[currency], value)
		}
	}

	return result
}

// AddPortfolio appends a new set of portfolio values to the history.
func (e *Equity) AddPortfolio(portfolio map[Currency]float64) {
	element := internal.MapCopy(portfolio)
	e.portfolioHistory = internal.Append(e.portfolioHistory, element)
}

// AddValue appends a new value to the main history and associates it with a timestamp.
func (e *Equity) AddValue(value float64, timestamp time.Time) {
	e.mainHistory = internal.Append(e.mainHistory, value)
	e.Timestamp.Append(timestamp)
}

// Now returns the most recent value from the main history.
func (e *Equity) Now() float64 {
	return e.mainHistory[len(e.mainHistory)-1]
}

// Start returns the timestamp of the first recorded value in main history.
func (e *Equity) Start() time.Time { return e.Timestamp.At(0) }

// NewEquity creates and returns a new Equity instance with specified timeframe
// and capacity.
func NewEquity(
	timeframe TimeFrame,
	capacity int,
) *Equity {
	history := make([]float64, 0, capacity)
	timestamp := NewTimeStamp(timeframe, capacity)

	return &Equity{
		portfolioHistory: make([]map[Currency]float64, 0, internal.DefaultCapacity),
		mainHistory:      history,
		Timestamp:        timestamp,
		timeframe:        timeframe,
	}
}

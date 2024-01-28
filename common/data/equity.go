package data

import (
	"time"
	"xoney/internal"
)

type Equity struct {
	portfolioHistory []map[Currency]float64
	mainHistory      []float64
	Timestamp        TimeStamp
	timeframe        TimeFrame
}

func (e *Equity) Timeframe() TimeFrame {
	return e.timeframe
}

func (e *Equity) Deposit() []float64 { return e.mainHistory }
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

func (e *Equity) AddPortfolio(portfolio map[Currency]float64) {
	element := internal.MapCopy(portfolio)
	e.portfolioHistory = internal.Append(e.portfolioHistory, element)
}

func (e *Equity) AddValue(value float64, timestamp time.Time) {
	e.mainHistory = internal.Append(e.mainHistory, value)
	e.Timestamp.Append(timestamp)
}

func (e *Equity) Now() float64 {
	return e.mainHistory[len(e.mainHistory)-1]
}
func (e *Equity) Start() time.Time { return e.Timestamp.At(0) }

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

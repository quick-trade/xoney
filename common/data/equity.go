package data

import (
	"time"

	"xoney/internal"
)

type Equity struct {
	startTime time.Time
	portfolioHistory []map[Currency]float64
	mainHistory   []float64
	Timestamp TimeStamp
	timeframe TimeFrame
}

func (e *Equity) Timeframe() TimeFrame {
	return e.timeframe
}

func (e *Equity) Deposit() []float64 { return e.mainHistory }
func (e *Equity) PortfolioHistory() map[Currency][]float64 {
	last := e.portfolioHistory[len(e.portfolioHistory)-1]
	result := make(map[Currency][]float64, len(last))

	for currency := range last {
		result[currency] = make([]float64, 0, len(e.mainHistory))
	}

	for i := range e.mainHistory {
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
func (e *Equity) AddValue(value float64) {
	e.mainHistory = internal.Append(e.mainHistory, value)
	if e.Timestamp.Len() == 0 {
		e.Timestamp.Append(e.startTime)
	} else {
		e.Timestamp.Extend(1)
	}
}

func (e *Equity) Now() float64 {
	return e.mainHistory[len(e.mainHistory)-1]
}
func (e *Equity) Start() time.Time { return e.startTime }

func NewEquity(
	capacity int,
	timeframe TimeFrame,
	start time.Time,
) *Equity {
	history := make([]float64, 0, capacity)
	timestamp := NewTimeStamp(timeframe, capacity)

	return &Equity{
		startTime: start,
		portfolioHistory: make([]map[Currency]float64, 0, internal.DefaultCapacity),
		mainHistory:   history,
		Timestamp: timestamp,
		timeframe: timeframe,
	}
}

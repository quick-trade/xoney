package data

import (
	"time"
	"xoney/internal"
)

type Equity struct {
	history   []float64
	Timestamp TimeStamp
	timeframe TimeFrame
}

func (e *Equity) Timeframe() TimeFrame {
	return e.timeframe
}

func (e *Equity) Deposit() []float64 { return e.history }
func (e *Equity) AddValue(value float64) {
	e.history = internal.Append(e.history, value)
	e.Timestamp.Extend(1)
}

func (e *Equity) Now() float64 {
	return e.history[len(e.history)-1]
}

func NewEquity(
	capacity int,
	timeframe TimeFrame,
	start time.Time,
	initialDepo float64,
) *Equity {
	history := make([]float64, 0, capacity)
	history = internal.Append(history, initialDepo)
	timestamp := NewTimeStamp(timeframe, capacity)
	timestamp.Append(start)

	return &Equity{
		history:   history,
		Timestamp: timestamp,
		timeframe: timeframe,
	}
}

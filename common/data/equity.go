package data

import (
	"time"
	"xoney/internal"
)

type Equity struct {
	startTime time.Time
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
	if e.Timestamp.Len() == 0 {
		e.Timestamp.Append(e.startTime)
	} else {
		e.Timestamp.Extend(1)
	}
}

func (e *Equity) Now() float64 {
	return e.history[len(e.history)-1]
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
		history:   history,
		Timestamp: timestamp,
		timeframe: timeframe,
	}
}

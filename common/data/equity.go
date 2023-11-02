package data

import (
	"time"
	"xoney/common"
	"xoney/internal"
)

type Equity struct {
	history   []float64
	timestamp common.TimeStamp
	timeframe TimeFrame
}

func (e *Equity) Deposit() []float64 { return e.history }
func (e *Equity) AddValue(value float64) {
	e.history = internal.Append(e.history, value)
	last_time := e.timestamp[len(e.timestamp)-1]
	new_time := last_time.Add(e.timeframe.Duration)
	e.timestamp = internal.Append(e.timestamp, new_time)
}
func (e *Equity) Now() float64 {
	return e.history[len(e.timestamp)-1] }

func NewEquity(
	capacity int,
	timeframe TimeFrame,
	start time.Time,
	initialDepo float64,
	) *Equity {
		history  := make([]float64, 0, capacity)
		history = internal.Append(history, initialDepo)
		timestamp := common.NewTimeStamp(capacity)
		timestamp = internal.Append(timestamp, start)
	return &Equity{
		history:  history ,
		timestamp: timestamp,
		timeframe: timeframe,
	}
}

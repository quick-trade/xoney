package data

import (
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
}

func NewEquity(capacity int, timeframe TimeFrame) *Equity {
	return &Equity{
		history:   make([]float64, 0, capacity),
		timestamp: common.NewTimeStamp(capacity),
		timeframe: timeframe,
	}
}

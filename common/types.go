package common

import (
	"time"
)

type Period [2]time.Time

func (p Period) ShiftedStart(shift time.Duration) Period {
	p[0] = p[0].Add(shift)

	return p
}

type Result[T any] struct {
	Data  T
	Error error
}

func NewResult[T any](data T, err error) Result[T] {
	return Result[T]{Data: data, Error: err}
}

type TimeStamp []time.Time

func NewTimeStamp(capacity int) TimeStamp {
	return make(TimeStamp, 0, capacity)
}

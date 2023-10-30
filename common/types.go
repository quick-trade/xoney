package common

import (
	"time"
)

type Period [2]time.Time

func (p Period) ShiftedStart(shift time.Duration) Period {
	p[1] = p[1].Add(shift)

	return p
}
func NewPeriod()

type Result[T any] struct {
	Data  T
	Error error
}

func NewResult[T any](data T, err error) Result[T] {
	return Result[T]{Data: data, Error: err}
}

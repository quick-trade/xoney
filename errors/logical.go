package errors

import (
	"strings"
	"time"
)

type ZeroLengthChartError struct {
	Msg string
}

func (e ZeroLengthChartError) Error() string {
	return e.Msg
}

type IncorrectSymbolError struct{}

func (e IncorrectSymbolError) Error() string {
	return "incorrect symbol initialization"
}

type IncorrectDurationError struct {
	Duration time.Duration
}

func (e IncorrectDurationError) Error() string {
	var msg strings.Builder

	msg.WriteString("invalid duration: ")
	msg.WriteString(e.Duration.String())

	return msg.String()
}

func NewIncorrectDurationError(duration time.Duration) IncorrectDurationError {
	return IncorrectDurationError{Duration: duration}
}

type UnsuccessfulChartSlicingError struct{}

func (e UnsuccessfulChartSlicingError) Error() string {
	return "cannot slice a chart"
}

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

type IncorrectDurationError struct {
	Duration time.Duration
}

func (e IncorrectDurationError) Error() string {
	var msg strings.Builder
	msg.WriteString("Invalid duration: ")
	msg.WriteString(e.Duration.String())
	return msg.String()
}
func NewIncorrectDurationError(duration time.Duration) IncorrectDurationError {
	return IncorrectDurationError{Duration: duration}
}

package errors

import (
	"strings"
	"time"
)

type ZeroLengthError struct {
	Character string
}

func NewZeroLengthError(characterName string) ZeroLengthError {
	return ZeroLengthError{Character: characterName}
}

func (e ZeroLengthError) Error() string {
	var msg strings.Builder

	msg.WriteString("character ")
	msg.WriteString(e.Character)
	msg.WriteString(" has 0 length.")

	return e.Character
}

type IncorrectSymbolError struct{}

func (e IncorrectSymbolError) Error() string {
	return "incorrect symbol initialization."
}

type IncorrectDurationError struct {
	Duration time.Duration
}

func (e IncorrectDurationError) Error() string {
	var msg strings.Builder

	msg.WriteString("invalid duration: ")
	msg.WriteString(e.Duration.String())
	msg.WriteString(".")

	return msg.String()
}

func NewIncorrectDurationError(duration time.Duration) IncorrectDurationError {
	return IncorrectDurationError{Duration: duration}
}

type UnsuccessfulChartSlicingError struct{}

func (e UnsuccessfulChartSlicingError) Error() string {
	return "cannot slice a chart."
}

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

type MissingCurrencyError struct {
	currencies []string
}

func NewMissingCurrencyError(capacity int) MissingCurrencyError {
	return MissingCurrencyError{currencies: make([]string, 0, capacity)}
}

func (m *MissingCurrencyError) Add(currency string) {
	m.currencies = append(m.currencies, currency)
}
func (e MissingCurrencyError) Error() string {
	var msg strings.Builder

	msg.WriteString("missed currency(ies): ")
	for _, currency := range e.currencies {
		msg.WriteString(currency)
	}

	return msg.String()
}

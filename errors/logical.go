package errors

import (
	"strconv"
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

func (m MissingCurrencyError) Error() string {
	var msg strings.Builder

	msg.WriteString("missed currencies: ")

	for _, currency := range m.currencies {
		msg.WriteString(currency)
		msg.WriteString(", ")
	}

	msg.WriteRune('.')

	return msg.String()
}

type NotEnoughFundsError struct {
	Currency string
	Quantity float64
}

func (e NotEnoughFundsError) Error() string {
	var msg strings.Builder

	msg.WriteString("Not enough funds in portfolio: ")
	msg.WriteString(strconv.FormatFloat(e.Quantity, 'f', -1, 64))
	msg.WriteString(", ")
	msg.WriteString(e.Currency)

	msg.WriteRune('.')

	return msg.String()
}

func NewNotEnoughFundsError(currency string, quantity float64) NotEnoughFundsError {
	return NotEnoughFundsError{Currency: currency, Quantity: quantity}
}

type NoLimitOrderError struct {
	id uint64
}

func (e NoLimitOrderError) Error() string {
	var msg strings.Builder

	msg.WriteString("there is no such limit order with ID: ")
	msg.WriteString(strconv.FormatUint(e.id, 10))
	msg.WriteRune('.')

	return msg.String()
}

func NewNoLimitOrderError(id uint64) NoLimitOrderError {
	return NoLimitOrderError{id: id}
}

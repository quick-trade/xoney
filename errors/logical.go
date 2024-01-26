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

	return msg.String()
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

	msg.WriteString(strings.Join(m.currencies, ", "))

	msg.WriteRune('.')

	return msg.String()
}

type NotEnoughFundsError struct {
	Currency string
	Quantity float64
}

func (e NotEnoughFundsError) Error() string {
	var msg strings.Builder

	msg.WriteString("not enough funds in portfolio: ")
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

type InvalidOrderAmountError struct {
	Amount float64
}

func (e InvalidOrderAmountError) Error() string {
	var msg strings.Builder

	msg.WriteString("invalid order amount: ")
	msg.WriteString(strconv.FormatFloat(e.Amount, 'f', -1, 64))

	return msg.String()
}

func NewInvalidOrderAmountError(amount float64) InvalidOrderAmountError {
	return InvalidOrderAmountError{Amount: amount}
}

type InvalidSymbolError struct {
	Base  string
	Quote string
}

func (e InvalidSymbolError) Error() string {
	var msg strings.Builder
	msg.WriteString("invalid symbol: ")
	msg.WriteString(e.Base)
	msg.WriteString("/")
	msg.WriteString(e.Quote)
	return msg.String()
}

func NewInvalidSymbolError(base, quote string) InvalidSymbolError {
	return InvalidSymbolError{Base: base, Quote: quote}
}

type NoPriceError struct {
	currency string
}

func (e NoPriceError) Error() string {
	var msg strings.Builder
	msg.WriteString("there is no price for ")
	msg.WriteString(e.currency)
	return msg.String()
}

func NewNoPriceError(currency string) NoPriceError {
	return NoPriceError{currency: currency}
}

type InvalidWeightsError struct {
	sum float64
}

func (e InvalidWeightsError) Error() string {
	var msg strings.Builder

	msg.WriteString("invalid portfolio weights: sum of abs(weights): ")
	msg.WriteString(strconv.FormatFloat(e.sum, 'f', -1, 64))

	return msg.String()
}

func NewInvalidWeightsError(sum float64) InvalidWeightsError {
	return InvalidWeightsError{sum: sum}
}

type InvalidGridLevelAmountError struct {
	Amount float64
}

func (e InvalidGridLevelAmountError) Error() string {
	var msg strings.Builder

	msg.WriteString("invalid amount of grid level: ")
	msg.WriteString(strconv.FormatFloat(e.Amount, 'f', -1, 64))
	msg.WriteString(" (expected > 0)")

	return msg.String()
}

func NewInvalidGridLevelAmountError(amount float64) InvalidGridLevelAmountError {
	return InvalidGridLevelAmountError{Amount: amount}
}

type ParallelExecutionError struct {
	ErrorsList []string
}

func (e *ParallelExecutionError) Error() string {
	var msg strings.Builder

	msg.WriteString("errors occurred in parallel execution: ")
	msg.WriteString(strings.Join(e.ErrorsList, "; "))

	return msg.String()
}

func NewParallelExecutionError(errorsList []string) *ParallelExecutionError {
	return &ParallelExecutionError{ErrorsList: errorsList}
}

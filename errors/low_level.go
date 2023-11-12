package errors

import (
	"strconv"
	"strings"
)

type OutOfIndexError struct {
	Index int
}

func (e OutOfIndexError) Error() string {
	var builder strings.Builder

	builder.WriteString("index ")
	builder.WriteString(strconv.Itoa(e.Index))
	builder.WriteString(" is out of range.")

	return builder.String()
}

func NewOutOfIndexError(index int) OutOfIndexError {
	return OutOfIndexError{Index: index}
}

type ValueNotFoundError struct{}

func (e ValueNotFoundError) Error() string {
	return "value not found."
}

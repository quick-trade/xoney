package errors

import (
	"strconv"
	"strings"
)

type OutOfIndexError struct {
	index int
}

func (e OutOfIndexError) Error() string {
	var builder strings.Builder
	builder.WriteString("Index ")
	builder.WriteString(strconv.Itoa(e.index))
	builder.WriteString(" is out of range")
	return builder.String()
}

func NewOutOfIndexError(index int) OutOfIndexError {
	return OutOfIndexError{index: index}
}

type ValueNotFoundError struct{}

func (e ValueNotFoundError) Error() string {
	return "Value not found"
}

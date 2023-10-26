package errors

type ZeroLengthChartError struct {
	Msg string
}

func (e ZeroLengthChartError) Error() string {
	return e.Msg
}

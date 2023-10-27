package data

type Equity struct {
	history []float64
}

func (e *Equity) Deposit() []float64 { return e.history }
func (e *Equity) AddValue(value float64) {
	e.history = append(e.history, value)
}

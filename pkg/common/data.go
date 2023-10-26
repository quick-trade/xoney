package common

import (
	"time"
)

type Candle [4]float64

func (c *Candle) Open() float64 {
	return c[0]
}

func (c *Candle) High() float64 {
	return c[1]
}
func (c *Candle) Low() float64 {
	return c[2]
}
func (c *Candle) Close() float64 {
	return c[3]
}

func NewCandle(open, high, low, close float64) *Candle {
	return &Candle{open, high, low, close}
}

type Symbol struct {
	base     string
	quote    string
	exchange string // TODO: could be nil
	full     string
}

func (s *Symbol) String() string   { return s.full }
func (s *Symbol) Base() string     { return s.base }
func (s *Symbol) Quote() string    { return s.quote }
func (s *Symbol) Exchange() string { return s.exchange }
func NewSymbol(param string, rest ...string) (*Symbol, error) {
	switch len(rest) {
	case 0:
		return NewSymbolByAll()
	case 1:
		// TODO: implement
	}
}

type TimeFrame struct {
	TimeDelta     time.Duration
	Seconds       float32
	CandlesInYear int
	Name          string
}

func NewTimeFrame(seconds float32)

type Instrument interface {
	Symbol() Symbol
	Timeframe() TimeFrame
	ProcessEvent(event Event) []Event
}
type LinearInstrument struct {
	symbol    Symbol
	timeframe TimeFrame
}

func (l *LinearInstrument) ProcessEvent(event Event) []Event {
	return []Event{event}
}

func (l *LinearInstrument) Symbol() Symbol {
	return l.symbol
}

func (l *LinearInstrument) Timeframe() TimeFrame {
	return l.timeframe
}

type Chart struct {
	Open      []float64
	High      []float64
	Low       []float64
	Close     []float64
	Volume    []float64
	Timestamp []time.Time
} // TODO: add {get; set}

type ChartContainer struct {
	charts map[Instrument]*Chart
}

func (c *ChartContainer) ChartByInstrument(instrument Instrument) *Chart {
	return c.charts[instrument]
}

type Equity struct {
	history []float64
}

func (e *Equity) Deposit() []float64 { return e.history }
func (e *Equity) AddValue(value float64) {
	e.history = append(e.history, value)
}

package data

import (
	"time"

	"xoney/errors"
	"xoney/internal"
)

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
	panic("TODO: implement")
}

type TimeFrame struct {
	Duration       time.Duration
	Seconds        float64
	CandlesPerYear float64
	Name           string
}

func NewTimeFrame(duration time.Duration, name string) (*TimeFrame, error) {
	if duration <= 0 {
		return nil, errors.NewIncorrectDurationError(duration)
	}
	
	candles := internal.TimesInYear(duration)

	return &TimeFrame{
		Duration:       duration,
		Seconds:        duration.Seconds(),
		CandlesPerYear: candles,
		Name:           name,
	}, nil
}

type Instrument struct {
	symbol    Symbol
	timeframe TimeFrame
}

func (l *Instrument) Symbol() Symbol {
	return l.symbol
}

func (l *Instrument) Timeframe() TimeFrame {
	return l.timeframe
}

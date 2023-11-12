package data

import (
	"strings"
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
	var symbol Symbol
	switch len(rest) {
	case 2:
		symbol = symbolByBaseQuoteExchange(param, rest...)
	}

	return &symbol, nil
}

func symbolByBaseQuoteExchange(param string, rest ...string) Symbol {
	base := param
	quote := rest[0]
	exchange := rest[1]

	var full strings.Builder

	full.WriteString(exchange)
	full.WriteRune(':')
	full.WriteString(base)
	full.WriteRune('/')
	full.WriteString(quote)

	return Symbol{
		base:     base,
		quote:    quote,
		exchange: exchange,
		full:     full.String(),
	}
}

type TimeFrame struct {
	Duration       time.Duration
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
		CandlesPerYear: candles,
		Name:           name,
	}, nil
}

type Instrument struct {
	symbol    Symbol
	timeframe TimeFrame
}

func NewInstrument(symbol Symbol, timeframe TimeFrame) Instrument {
	return Instrument{
		symbol:    symbol,
		timeframe: timeframe,
	}
}

func (l *Instrument) Symbol() Symbol {
	return l.symbol
}

func (l *Instrument) Timeframe() TimeFrame {
	return l.timeframe
}

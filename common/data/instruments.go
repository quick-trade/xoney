package data

import (
	"strings"
	"time"

	"xoney/errors"
	"xoney/internal"
)

type Exchange string

type Currency struct {
	Asset    string
	Exchange Exchange
}

func (c Currency) String() string {
	var str strings.Builder

	str.WriteString(string(c.Exchange))
	str.WriteRune(':')
	str.WriteString(c.Asset)

	return str.String()
}

func NewCurrency[E Exchange | string](asset string, exchange E) Currency {
	return Currency{
		Asset:    asset,
		Exchange: Exchange(exchange),
	}
}

type Symbol struct {
	base  Currency
	quote Currency
}

func (s Symbol) String() string {
	var full strings.Builder

	full.WriteString(string(s.base.Exchange))
	full.WriteRune(':')
	full.WriteString(s.base.Asset)
	full.WriteRune('/')
	full.WriteString(s.quote.Asset)

	return full.String()
}

func (s Symbol) Base() Currency     { return s.base }
func (s Symbol) Quote() Currency    { return s.quote }
func (s Symbol) Exchange() Exchange { return s.base.Exchange }
func NewSymbol[E Exchange | string](base, quote string, exchange E) *Symbol {

	return &Symbol{
		base:  NewCurrency(base, exchange),
		quote: NewCurrency(quote, exchange),
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

func (i *Instrument) Symbol() Symbol {
	return i.symbol
}

func (i *Instrument) Timeframe() TimeFrame {
	return i.timeframe
}

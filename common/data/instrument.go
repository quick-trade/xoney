package data

import (
	"strings"
	"time"

	"github.com/quick-trade/xoney/errors"
	"github.com/quick-trade/xoney/internal"
)

type Exchange string

// Currency represents a financial asset on a specific exchange.
// It contains the asset identifier (e.g., "BTC" for Bitcoin) and
// the exchange where it is traded.
type Currency struct {
	Asset    string   // Identifier of the financial asset.
	Exchange Exchange // Exchange on which the asset is traded.
}

// String returns the string representation of the Currency in the format
// "Exchange:Asset". For example, if the Exchange is "NASDAQ" and the Asset is "AAPL",
// the String method would return "NASDAQ:AAPL".
func (c Currency) String() string {
	var str strings.Builder

	str.WriteString(string(c.Exchange))
	str.WriteRune(':')
	str.WriteString(c.Asset)

	return str.String()
}

// NewCurrency is a Currency constructor.
func NewCurrency[E Exchange | string](asset string, exchange E) Currency {
	return Currency{
		Asset:    asset,
		Exchange: Exchange(exchange),
	}
}

// Symbol represents a trading pair in the format 'Base/Quote'.
// The 'base' is the currency being bought or sold, and 'quote' is the currency
// that the 'base' is priced in.
type Symbol struct {
	base  Currency // Currency being traded.
	quote Currency // Currency used to price the 'base'.
}

// String returns the string representation of the Symbol in the format
// "Exchange:Base/Quote". For example, if the 'base' is "BTC" on "NASDAQ" and
// the 'quote' is "USD", the String method would return "NASDAQ:BTC/USD".
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

// NewSymbol is a Symbol constructor.
func NewSymbol[E Exchange | string](base, quote string, exchange E) *Symbol {
	return &Symbol{
		base:  NewCurrency(base, exchange),
		quote: NewCurrency(quote, exchange),
	}
}

func NewSymbolFromCurrencies(base, quote Currency) *Symbol {
	return &Symbol{
		base:  base,
		quote: quote,
	}
}

// TimeFrame represents a discretization period of a trading instrument.
type TimeFrame struct {
	Duration       time.Duration
	CandlesPerYear float64
	Name           string
}

// NewTimeFrame creates a new TimeFrame with the specified duration and name.
// It returns an error if the duration is not positive.
// Inputs:
// - duration: The length of the time frame in time.Duration format.
// - name: The display name of the time frame.
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

// Instrument represents the full description of a trading instrument in financial markets.
// It encapsulates all concrete details required for trading with no further abstractions for tradable currencies.
// The 'symbol' uniquely identifies the financial asset, and the 'timeframe' defines the trading interval.
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

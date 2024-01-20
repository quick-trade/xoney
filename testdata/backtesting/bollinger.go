package backtesting

import (
	"fmt"
	"math"
	"time"

	"errors"
	"xoney/common/data"
	"xoney/events"
	"xoney/exchange"
	"xoney/internal"
	"xoney/strategy"
)

type Decision string

const (
	BUY  Decision   = "buy"
	SELL  Decision  = "sell"
	NEUTRAL Decision= "neutral"
)

type BBBStrategy struct {
	Period     int
	Deviation  float64
	instrument data.Instrument
	prices []float64
	side Decision
	prevSide Decision
	Mean []float64
	UB []float64
	LB []float64
}

func NewBBStrategy(period int, deviation float64, instrument data.Instrument) *BBBStrategy {
	return &BBBStrategy{
		Period:     period,
		Deviation:  deviation,
		instrument: instrument,
		prices: make([]float64, 0, internal.DefaultCapacity),
		side: NEUTRAL,
		prevSide: NEUTRAL,
		Mean: make([]float64, 0, internal.DefaultCapacity),
		UB: make([]float64, 0, internal.DefaultCapacity),
		LB: make([]float64, 0, internal.DefaultCapacity),
	}
}

func (b *BBBStrategy) Next(candle data.InstrumentCandle) (events.Event, error) {
	b.computeBollinger(candle.Close)
	return b.fetchEvents(candle.Close)
}

func (b *BBBStrategy) Start(charts data.ChartContainer) error {
	b.prices = charts[b.instrument].Close
	return nil
}

func (b BBBStrategy) MinDurations() strategy.Durations {
	return strategy.Durations{
		b.instrument: b.instrument.Timeframe().Duration * time.Duration(b.Period),
	}
}

func (b *BBBStrategy) computeBollinger(price float64) {
	b.prices = append(b.prices[1:], price)

    mean, _ := internal.RawMoment(b.prices, 1)
    variance := internal.CentralMoment(b.prices, mean, 2)
    std := math.Sqrt(variance)

    if price > mean+b.Deviation*std {
        b.side = BUY
    } else if price < mean-b.Deviation*std {
        b.side = SELL
    }
    b.Mean = append(b.Mean, mean)
    b.UB = append(b.UB, mean+b.Deviation*std)
    b.LB = append(b.LB, mean-b.Deviation*std)
}
func (b *BBBStrategy) fetchEvents(price float64) (events.Event, error) {
	resultEvents := events.NewSequential()

	var event events.Event
	var err error = nil

	if b.side != b.prevSide {
		if b.side == BUY {
			event, err = NewEntryAllDeposit(b.instrument.Symbol(), "market", "buy", price)
		} else if b.side == SELL {
			event, err = NewEntryAllDeposit(b.instrument.Symbol(), "market", "sell", price)
		}
		resultEvents.Add(event)
	}
	b.prevSide = b.side


	return resultEvents, err
}

type VectorizedBollinger struct {
	BBBStrategy
}

func (b *VectorizedBollinger) Backtest(
	simulator exchange.Simulator,
	charts data.ChartContainer,
) (data.Equity, error) {
	initialDepo := simulator.Portfolio().Balance(b.instrument.Symbol().Quote())

	chart := charts[b.instrument]

	startIndex := b.Period - 1
	start := chart.Timestamp.At(startIndex)
	equity := *data.NewEquity(b.instrument.Timeframe(), start, len(chart.Close))
	equity.AddValue(initialDepo, start)

	price := chart.Close
	flag := NEUTRAL

	var diff float64

	for i, p := range price {
		price[i] = math.Log(p)
	}

	average, _ := internal.RawMoment(price[:b.Period+1], 1)

	for i := startIndex; i < len(price); i++ {
		diff = price[i] - price[i-1]
		moment := chart.Timestamp.At(i)

		if flag == BUY {
			equity.AddValue(equity.Now() + diff, moment)
		} else if flag == SELL {
			equity.AddValue(equity.Now() - diff, moment)
		}

		average += (price[i] - price[i-b.Period+1]) / float64(b.Period)

		stdDev := 0.0
		for j := 0; j < b.Period; j++ {
			deviation := price[i-j] - average
			stdDev += deviation * deviation
		}
		stdDev = math.Sqrt(stdDev / float64(b.Period))

		upperBand := average + b.Deviation*stdDev
		lowerBand := average - b.Deviation*stdDev

		if price[i] > upperBand {
			flag = BUY
		}

		if price[i] < lowerBand {
			flag = SELL
		}
	}

	return equity, nil
}

type EntryAllDeposit struct {
	symbol data.Symbol
	orderType exchange.OrderType
	side exchange.OrderSide
	price float64
}

func (b *EntryAllDeposit) Occur(connector exchange.Connector) error {
	var neededCurrency data.Currency

	if b.side == exchange.Buy {
		neededCurrency = b.symbol.Quote()
	} else if b.side == exchange.Sell {
		neededCurrency = b.symbol.Base()
	} else {
		return errors.New("unknown side: " + string(b.side))
	}

	amount := connector.Portfolio().Balance(neededCurrency)

	if b.side == exchange.Buy {
		amount /= b.price
	}

	order, err := exchange.NewOrder(b.symbol, b.orderType, b.side, b.price, amount)
	if err != nil {
		return fmt.Errorf("error creating order: %w", err)
	}

	return events.NewOpenOrder(*order).Occur(connector)
}
func NewEntryAllDeposit(symbol data.Symbol, orderType, side string, price float64) (*EntryAllDeposit, error) {
	var Type exchange.OrderType
	if orderType == "market" {
		Type = exchange.Market
	} else if orderType == "limit" {
		Type = exchange.Limit
	} else {
		return nil, errors.New("no such order type")
	}

	var Side exchange.OrderSide
	if side == "buy" {
		Side = exchange.Buy
	} else if side == "sell" {
		Side = exchange.Sell
	} else {
		return nil, errors.New("no such order type")
	}

	return &EntryAllDeposit{
		symbol: symbol,
		orderType: Type,
		side: Side,
		price: price,
	}, nil
}

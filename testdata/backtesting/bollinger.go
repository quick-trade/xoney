package backtesting

import (
	"math"
	"time"

	"xoney/common/data"
	"xoney/events"
	"xoney/exchange"
	"xoney/internal"
	"xoney/strategy"
)

const (
	BUY     = iota
	SELL    = iota
	NEUTRAL = iota
)

type BBBStrategy struct {
	Period     int
	Deviation  float64
	instrument data.Instrument
	chart      data.Chart
	first bool
}

func (b *BBBStrategy) Bacdktest(
	initialDepo float64,
	charts data.ChartContainer,
) (data.Equity, error) {
	b.chart = charts[b.instrument]
	equity := *data.NewEquity(len(b.chart.Close), b.instrument.Timeframe(), time.Now())
	equity.AddValue(initialDepo)

	price := b.chart.Close
	flag := NEUTRAL

	var diff float64

	for i, p := range price {
		price[i] = math.Log(p)
	}

	average, _ := internal.RawMoment(price[:b.Period+1], 1)

	for i := b.Period - 1; i < len(price); i++ {
		diff = price[i] - price[i-1]

		if flag == BUY {
			equity.AddValue(equity.Now() + diff)
		} else if flag == SELL {
			equity.AddValue(equity.Now() - diff)
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
func NewBBStrategy(period int, deviation float64, instrument data.Instrument) *BBBStrategy {
	return &BBBStrategy{
		Period:     period,
		Deviation:  deviation,
		instrument: instrument,
		chart:      data.RawChart(instrument.Timeframe(), 0),
		first: true,
	}
}

func (b *BBBStrategy) Next(candle data.InstrumentCandle) ([]events.Event, error) {
	if !b.first { return []events.Event{}, nil}
	if candle.Instrument == b.instrument {
		b.first = false
		return []events.Event{
			events.NewOpenOrder(
				*exchange.NewOrder(candle.Symbol(), exchange.Limit, exchange.Buy, 16510, 17099.96/candle.Close),
			),
			events.NewOpenOrder(
				*exchange.NewOrder(candle.Symbol(), exchange.Limit, exchange.Sell, 24000, 17099.96/candle.Close),
			),
			}, nil
		}

	return []events.Event{}, nil
}

func (b *BBBStrategy) Start(charts data.ChartContainer) error {
	return nil
}

func (b BBBStrategy) MinDurations() strategy.Durations {
	return strategy.Durations{
		b.instrument: b.instrument.Timeframe().Duration * time.Duration(b.Period),
	}
}

func (b *BBBStrategy) Parameters() []strategy.Parameter {
	return []strategy.Parameter{
		strategy.NewIntParameter("Period", 1, 600),
		strategy.NewFloatParameter("Deviation", 0, 4),
	}
}

package backtesting

import (
	"math"
	"time"
	"xoney/common/data"
	"xoney/events"
)

const (
	BUY = iota
	SELL = iota
	NEUTRAL = iota
)

type BBBStrategy struct {
	Period int
	Deviation float64
	instrument data.Instrument
	chart data.Chart
}

func (b *BBBStrategy) Backtest(commission float64, initialDepo float64, charts data.ChartContainer) (data.Equity, error) {
	b.chart = charts[b.instrument]
	equity := *data.NewEquity(len(b.chart.Close), b.instrument.Timeframe(), time.Now(), 1)
	equity.AddValue(initialDepo)

	price := b.chart.Close
	flag := NEUTRAL

	var diff float64

	// Расчет Bollinger Bands на основе цен
	for i := b.Period-1; i < len(price); i++ {
		// Рассчитываем среднее значение цены за период
		sum := 0.0
		for j := 0; j < b.Period; j++ {
			sum += price[i-j]
		}
		average := sum / float64(b.Period)

		// Рассчитываем стандартное отклонение цены за период
		stdDev := 0.0
		for j := 0; j < b.Period; j++ {
			deviation := price[i-j] - average
			stdDev += deviation * deviation
		}
		stdDev = math.Sqrt(stdDev / float64(b.Period))

		// Рассчитываем верхнюю и нижнюю полосы Bollinger Bands
		upperBand := average + b.Deviation*stdDev
		lowerBand := average - b.Deviation*stdDev

		// Ваша логика для торговли на основе Bollinger Bands
		// Здесь можно определить условия покупки и продажи

		// Пример условия: если цена закрытия выше верхней полосы, покупаем
		if price[i] > upperBand {
			flag = BUY
		}

		// Пример условия: если цена закрытия ниже нижней полосы, продаем
		if price[i] < lowerBand {
			flag = SELL
		}
		diff = price[i] - price[i-1]

		if flag == BUY {
			equity.AddValue(equity.Now()+diff)
		} else if flag == SELL {
			equity.AddValue(equity.Now()-diff)
		}
	}

	return equity, nil
}
func NewBBStrategy(period int, deviation float64, instrument data.Instrument) *BBBStrategy {
	return &BBBStrategy{
		Period: period,
		Deviation: deviation,
		instrument: instrument,
		chart: data.RawChart(0),
	}
}

func (b *BBBStrategy) MinDuration() time.Duration {
	return b.instrument.Timeframe().Duration * time.Duration(b.Period)
}

func (b *BBBStrategy) Next(candle data.Candle) []events.Event {
	panic("not implemented")
}

func (b *BBBStrategy) Start(charts data.ChartContainer) {
	panic("not implemented")
}

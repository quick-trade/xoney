package data

import (
	"time"

	"xoney/common"
	"xoney/internal"
	"xoney/internal/search"
)

type Candle struct {
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
	Timestamp time.Time
}

func NewCandle(open, high, low, c, volume float64, timestamp time.Time) *Candle {
	return &Candle{
		Open:      open,
		High:      high,
		Low:       low,
		Close:     c,
		Volume:    volume,
		Timestamp: timestamp,
	}
}

type Chart struct {
	Open      []float64
	High      []float64
	Low       []float64
	Close     []float64
	Volume    []float64
	Timestamp common.TimeStamp
}

func RawChart(capacity int) Chart {
	return Chart{
		Open:      make([]float64, 0, capacity),
		High:      make([]float64, 0, capacity),
		Low:       make([]float64, 0, capacity),
		Close:     make([]float64, 0, capacity),
		Volume:    make([]float64, 0, capacity),
		Timestamp: make(common.TimeStamp, 0, capacity),
	}
}

func (c *Chart) Add(candle Candle) {
	c.Open = internal.Append(c.Open, candle.Open)
	c.High = internal.Append(c.High, candle.High)
	c.Low = internal.Append(c.Low, candle.Low)
	c.Close = internal.Append(c.Close, candle.Close)
	c.Volume = internal.Append(c.Volume, candle.Volume)
	c.Timestamp = internal.Append(c.Timestamp, candle.Timestamp)
}

func (c *Chart) Slice(period common.Period) Chart {
	start, err := search.LastBeforeIdx(c.Timestamp, period[0])
	if err != nil {
		return RawChart(0)
	}

	stop, err := search.LastBeforeIdx(c.Timestamp, period[1])
	if err != nil {
		return RawChart(0)
	}

	return Chart{
		Open:      c.Open[start:stop],
		High:      c.High[start:stop],
		Low:       c.Low[start:stop],
		Close:     c.Close[start:stop],
		Volume:    c.Volume[start:stop],
		Timestamp: c.Timestamp[start:stop],
	}
}

type ChartContainer map[Instrument]Chart

func (c *ChartContainer) ChartsByPeriod(period common.Period) ChartContainer {
	result := make(ChartContainer, len(*c))
	for instrument, chart := range *c {
		result[instrument] = chart.Slice(period)
	}

	return result
}

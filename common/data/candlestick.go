package data

import (
	"time"

	"xoney/common"
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
	Timestamp []time.Time
}

func (c *Chart) Slice(period common.Period) Chart {
	panic("TODO: implement")
}

type ChartContainer map[Instrument]Chart

func (c *ChartContainer) ChartsByPeriod(period common.Period) ChartContainer {
	result := make(ChartContainer, len(*c))
	for instrument, chart := range *c {
		result[instrument] = chart.Slice(period)
	}

	return result
}

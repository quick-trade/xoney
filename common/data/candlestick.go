package data

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

package data

import (
	"time"

	"xoney/common"
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

func NewCandle(open, high, low, c float64) *Candle {
	return &Candle{open, high, low, c}
}

type Chart struct {
	Open      []float64
	High      []float64
	Low       []float64
	Close     []float64
	Volume    []float64
	Timestamp []time.Time
}

func (c *Chart) Slice(period common.Period) (Chart, error) {
	panic("not implemented")
}

type ChartContainer map[Instrument]Chart

func (c *ChartContainer) ChartsByPeriod(period common.Period) (ChartContainer, error) {
	var err error
	success := true
	result := make(ChartContainer, len(*c))
	for instrument, chart := range *c {
		result[instrument], err = chart.Slice(period)
		if err != nil {
			success = false
		}
	}
	if !success {
		err = 
	}

	return result, err
}

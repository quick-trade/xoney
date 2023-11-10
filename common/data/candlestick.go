package data

import (
	"time"
	"xoney/errors"
	"xoney/internal"
)

type Period [2]time.Time

func (p Period) ShiftedStart(shift time.Duration) Period {
	p[0] = p[0].Add(shift)

	return p
}

type TimeStamp struct {
	timeframe TimeFrame
	Timestamp []time.Time
}

func NewTimeStamp(timeframe TimeFrame, capacity int) TimeStamp {
	return TimeStamp{
		timeframe: timeframe,
		Timestamp: make([]time.Time, 0, capacity),
	}
}

func (t TimeStamp) Timeframe() TimeFrame {
	return t.timeframe
}

func (t *TimeStamp) Extend(n int) {
	last := t.Timestamp[len(t.Timestamp)-1]
	new := last.Add(t.timeframe.Duration)
	t.Timestamp = internal.Append(t.Timestamp, new)
}

func (t *TimeStamp) Append(moments ...time.Time) {
	t.Timestamp = internal.Append(t.Timestamp, moments...)
}

func (t TimeStamp) sliceIdx(start, stop int) TimeStamp {
	return TimeStamp{
		timeframe: t.timeframe,
		Timestamp: t.Timestamp[start:stop],
	}
}

func (t TimeStamp) End() time.Time {
	return t.Timestamp[len(t.Timestamp)-1]
}

func (t TimeStamp) Start() time.Time {
	return t.Timestamp[0]
}

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

type InstrumentCandle struct {
	Candle
	Instrument
}

type Chart struct {
	Open      []float64
	High      []float64
	Low       []float64
	Close     []float64
	Volume    []float64
	Timestamp TimeStamp
}

func RawChart(timeframe TimeFrame, capacity int) Chart {
	return Chart{
		Open:      make([]float64, 0, capacity),
		High:      make([]float64, 0, capacity),
		Low:       make([]float64, 0, capacity),
		Close:     make([]float64, 0, capacity),
		Volume:    make([]float64, 0, capacity),
		Timestamp: NewTimeStamp(timeframe, capacity),
	}
}

func (c *Chart) Add(candle Candle) {
	c.Open = internal.Append(c.Open, candle.Open)
	c.High = internal.Append(c.High, candle.High)
	c.Low = internal.Append(c.Low, candle.Low)
	c.Close = internal.Append(c.Close, candle.Close)
	c.Volume = internal.Append(c.Volume, candle.Volume)
	c.Timestamp.Extend(1)
}

func (c *Chart) Slice(period Period) Chart {
	start, err := findIndexBeforeOrAtTime(c.Timestamp, period[0])
	if err != nil {
		return RawChart(c.Timestamp.Timeframe(), 0)
	}

	stop, err := findIndexBeforeOrAtTime(c.Timestamp, period[1])
	if err != nil {
		return RawChart(c.Timestamp.Timeframe(), 0)
	}

	return Chart{
		Open:      c.Open[start:stop],
		High:      c.High[start:stop],
		Low:       c.Low[start:stop],
		Close:     c.Close[start:stop],
		Volume:    c.Volume[start:stop],
		Timestamp: c.Timestamp.sliceIdx(start, stop),
	}
}

type ChartContainer map[Instrument]Chart

func (c *ChartContainer) ChartsByPeriod(period Period) ChartContainer {
	result := make(ChartContainer, len(*c))
	for instrument, chart := range *c {
		result[instrument] = chart.Slice(period)
	}

	return result
}

func (c *ChartContainer) Candles() []InstrumentCandle {
	panic("TODO: Implement")
}

func findIndexBeforeOrAtTime(
	series TimeStamp,
	moment time.Time,
) (int, error) {
	lastIndex := len(series.Timestamp) - 1

	if len(series.Timestamp) == 0 {
		return -1, errors.NewZeroLengthError("timestamp series")
	}

	begin := series.Timestamp[0]
	if moment.Before(begin) {
		return -1, errors.ValueNotFoundError{}
	}

	idx := int(moment.Sub(begin) / series.timeframe.Duration)

	if idx > lastIndex {
		idx = lastIndex
	}

	return idx, nil
}

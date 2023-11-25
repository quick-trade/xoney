package data

import (
	"sort"
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

func (t *TimeStamp) Timeframe() TimeFrame {
	return t.timeframe
}

func (t TimeStamp) At(index int) time.Time {
	return t.Timestamp[index]
}

func (t *TimeStamp) Extend(n int) {
	last := t.At(len(t.Timestamp) - 1)
	for i := 0; i < n; i++ {
		last = last.Add(t.timeframe.Duration)
		t.Timestamp = internal.Append(t.Timestamp, last)
	}
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
	return t.At(len(t.Timestamp) - 1)
}

func (t TimeStamp) Start() time.Time {
	return t.At(0)
}

type Candle struct {
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
	TimeClose time.Time
}

func NewCandle(open, high, low, c, volume float64, timeClose time.Time) *Candle {
	return &Candle{
		Open:      open,
		High:      high,
		Low:       low,
		Close:     c,
		Volume:    volume,
		TimeClose: timeClose,
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
		return RawChart(c.Timestamp.timeframe, 0)
	}

	stop, err := findIndexBeforeOrAtTime(c.Timestamp, period[1])
	if err != nil {
		return RawChart(c.Timestamp.timeframe, 0)
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

func (c *Chart) Len() int {
	return len(c.Timestamp.Timestamp)
}

func (c *Chart) CandleByIndex(index int) (*Candle, error) {
	if index >= c.Len() {
		return nil, errors.OutOfIndexError{Index: index}
	}

	return NewCandle(
		c.Open[index],
		c.High[index],
		c.Low[index],
		c.Close[index],
		c.Volume[index],
		c.Timestamp.At(index),
	), nil
}

type ChartContainer map[Instrument]Chart

func (c *ChartContainer) ChartsByPeriod(period Period) ChartContainer {
	result := make(ChartContainer, len(*c))
	for instrument, chart := range *c {
		result[instrument] = chart.Slice(period)
	}

	return result
}

func (c *ChartContainer) sortedInstruments() []Instrument {
	keys := make([]Instrument, 0, len(*c))

	for instrument := range *c {
		keys = internal.Append(keys, instrument)
	}

	sort.Slice(keys, func(i, j int) bool {
		durationI := keys[i].timeframe.Duration
		durationJ := keys[j].timeframe.Duration

		return durationI < durationJ
	})

	return keys
}

func (c ChartContainer) Candles() []InstrumentCandle {
	sumLength := 0
	for _, chart := range c {
		sumLength += chart.Len()
	}

	result := make([]InstrumentCandle, 0, sumLength)
	pointers := make([]int, len(c))
	instruments := c.sortedInstruments()

	for {
		var minChart Chart
		var minInstrument Instrument
		var minKey int

		minIndex := -1
		minTime := time.Time{}

		for instIdx, inst := range instruments {
			minKey = instIdx
			idx := pointers[instIdx]
			chart := c[inst]
			moment := chart.Timestamp.At(idx)

			if idx < chart.Len() && (minIndex == -1 || moment.Before(minTime)) {
				minTime = moment
				minChart = chart
				minInstrument = inst
				minIndex = idx
			}
		}

		if minIndex == -1 {
			break
		}

		candle, _ := minChart.CandleByIndex(minIndex)
		instCandle := InstrumentCandle{Candle: *candle, Instrument: minInstrument}
		result = internal.Append(result, instCandle)

		pointers[minKey]++
	}

	return result
}

func findIndexBeforeOrAtTime(
	series TimeStamp,
	moment time.Time,
) (int, error) {
	lastIndex := len(series.Timestamp) - 1

	if len(series.Timestamp) == 0 {
		return -1, errors.NewZeroLengthError("timestamp series")
	}

	begin := series.At(0)
	if moment.Before(begin) {
		return -1, errors.ValueNotFoundError{}
	}

	idx := int(moment.Sub(begin) / series.timeframe.Duration)

	if idx > lastIndex {
		idx = lastIndex
	}

	return idx, nil
}

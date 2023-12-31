package data

import (
	"sort"
	"time"

	"xoney/errors"
	"xoney/internal"
)

type Period struct {
	Start time.Time
	End   time.Time
}

func NewPeriod(start, end time.Time) Period {
	return Period{Start: start, End: end}
}

func (p Period) ShiftedStart(shift time.Duration) Period {
	p.Start = p.Start.Add(shift)

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

func (t TimeStamp) Slice(start, stop int) TimeStamp {
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
func (t TimeStamp) Len() int { return len(t.Timestamp) }

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
	c.Timestamp.Append(candle.TimeClose)
}

func (c *Chart) Slice(period Period) Chart {
	start, err := findIndexBeforeOrAtTime(c.Timestamp, period.Start)
	if err != nil {
		return RawChart(c.Timestamp.timeframe, 0)
	}

	stop, _ := findIndexBeforeOrAtTime(c.Timestamp, period.End)
	// Any errors that might occur would be related
	// to the processing of the period start.
	stop++

	return Chart{
		Open:      c.Open[start:stop],
		High:      c.High[start:stop],
		Low:       c.Low[start:stop],
		Close:     c.Close[start:stop],
		Volume:    c.Volume[start:stop],
		Timestamp: c.Timestamp.Slice(start, stop),
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

func (c *ChartContainer) FirstStart() time.Time {
	var first time.Time

	for _, chart := range *c {
		start := chart.Timestamp.Start()
		if first.IsZero() || start.Before(first) {
			first = start
		}
	}

	return first
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
	// It is just merging from merge-sort algorithm
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

			if idx >= chart.Len() {
				break
			}

			moment := chart.Timestamp.At(idx)

			if minIndex == -1 || moment.Before(minTime) {
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
	if len(series.Timestamp) == 0 {
		return -1, errors.NewZeroLengthError("timestamp series")
	}

	begin := series.At(0)
	if moment.Before(begin) {
		return -1, errors.ValueNotFoundError{}
	}

	idx := binarySearch(series, moment)

	return idx, nil
}

func binarySearch(series TimeStamp, target time.Time) int {
	low, high := 0, len(series.Timestamp)-1

	var result int

	for low <= high {
		mid := (low + high) / 2
		midTime := series.At(mid)

		switch {
		case midTime.Before(target):
			low = mid + 1
			result = mid
		case midTime.After(target):
			high = mid - 1
		default:
			return mid
		}
	}

	return result
}

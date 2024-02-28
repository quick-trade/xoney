package data

import (
	"sort"
	"time"

	"github.com/quick-trade/xoney/errors"
	"github.com/quick-trade/xoney/internal"
)

// Period is a utility structure used for denoting intervals in time.
// It allows for convenient slicing of various time-series data.
type Period struct {
	Start time.Time
	End   time.Time
}

func NewPeriod(start, end time.Time) Period {
	return Period{Start: start, End: end}
}

// ShiftedStart returns a new Period with the start time shifted by the given duration.
// The end time of the Period is not affected.
func (p Period) ShiftedStart(shift time.Duration) Period {
	p.Start = p.Start.Add(shift)

	return p
}

// TimeStamp represents a sequence of time moments.
// These moments are not required to have a constant step between them,
// but the data should be sequential.
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
// At returns the time at the specified index within the TimeStamp.
func (t TimeStamp) At(index int) time.Time {
	return t.Timestamp[index]
}

// Extend increases the length of the TimeStamp by n timeframes,
// adding new moments sequentially based on the TimeStamp's timeframe.
func (t *TimeStamp) Extend(n int) {
	last := t.At(len(t.Timestamp) - 1)
	for i := 0; i < n; i++ {
		last = last.Add(t.timeframe.Duration)
		t.Timestamp = internal.Append(t.Timestamp, last)
	}
}

// Append adds the provided moments to the end of the TimeStamp.
// If you need to extend the TimeStamp by N timeframes, consider using Extend instead.
func (t *TimeStamp) Append(moments ...time.Time) {
	t.Timestamp = internal.Append(t.Timestamp, moments...)
}

// Slice returns a new TimeStamp consisting of the time moments within the range [start, stop).
func (t TimeStamp) Slice(start, stop int) TimeStamp {
	return TimeStamp{
		timeframe: t.timeframe,
		Timestamp: t.Timestamp[start:stop],
	}
}

// End returns the last time moment within the TimeStamp.
func (t TimeStamp) End() time.Time {
	return t.At(len(t.Timestamp) - 1)
}

// Start returns the first time moment within the TimeStamp.
func (t TimeStamp) Start() time.Time {
	return t.At(0)
}

// Len returns the number of time moments within the TimeStamp.
func (t TimeStamp) Len() int { return len(t.Timestamp) }

// Candle represents a single candlestick data point in a financial chart,
// encapsulating the open, high, low, close values and the volume of trading
// over a particular time period, with TimeClose marking the end of that period.
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
// InstrumentCandle represents a candlestick data point with an associated financial instrument.
// It combines the detailed candlestick information such as OHLCV with the specific instrument
// for which this data is relevant. This data structure is commonly used in trading strategies
// as the primary source of information from which trading signals are generated.
type InstrumentCandle struct {
	Candle
	Instrument
}

func NewInstrumentCandle(candle Candle, instrument Instrument) *InstrumentCandle {
	return &InstrumentCandle{
		Candle:     candle,
		Instrument: instrument,
	}
}

// Chart represents a sequence of candlestick data points for a specific instrument.
// It contains slices of open, high, low, close values and the trading volume,
// along with the corresponding timestamps for each data point.
type Chart struct {
	Open      []float64
	High      []float64
	Low       []float64
	Close     []float64
	Volume    []float64
	Timestamp TimeStamp
}

// RawChart creates a new Chart with slices preallocated to the specified capacity.
// This function is used to initialize a Chart for a certain timeframe and with a capacity
// hint to optimize memory allocations for the underlying slices.
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

// Add appends a single candle to the end of the chart.
// It appends the open, high, low, close, and volume data from the candle
// to their respective slices in the Chart and updates the timestamp.
func (c *Chart) Add(candle Candle) {
	c.Open = internal.Append(c.Open, candle.Open)
	c.High = internal.Append(c.High, candle.High)
	c.Low = internal.Append(c.Low, candle.Low)
	c.Close = internal.Append(c.Close, candle.Close)
	c.Volume = internal.Append(c.Volume, candle.Volume)
	c.Timestamp.Append(candle.TimeClose)
}

// Slice returns a new Chart consisting of the candlestick data within the specified period.
// It performs a binary search to find the start and end indices, ensuring O(log N) complexity.
// The resulting slice includes the boundaries of the period. If an exact match for the boundaries
// is not found, the nearest preceding element is included in the slice.
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
// Len returns the number of timestamps (and hence the number of candles) in the Chart.
func (c *Chart) Len() int {
	return len(c.Timestamp.Timestamp)
}

// CandleByIndex retrieves the candle at the specified index from the Chart.
// If the index is out of range, an error is returned.
// Parameters:
//   index - The index of the candle to retrieve.
// Returns:
//   pointer to a Candle and nil error if successful, or
//   nil pointer and an OutOfIndexError if the index is invalid.
func (c *Chart) CandleByIndex(index int) (*Candle, error) {
	if index >= c.Len() {
		return nil, errors.NewOutOfIndexError(index)
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
// ChartContainer represents a collection of instruments and their corresponding charts.
// It can be used to inform your trading system about your investment universe during testing and training.
type ChartContainer map[Instrument]Chart

// ChartsByPeriod slices each chart in the ChartContainer to a new chart based on the provided period.
// It's analogous to the Slice method of the Chart, but it operates on the entire container, slicing each chart.
func (c *ChartContainer) ChartsByPeriod(period Period) ChartContainer {
	result := make(ChartContainer, len(*c))
	for instrument, chart := range *c {
		result[instrument] = chart.Slice(period)
	}

	return result
}

// FirstStart finds the earliest start time among all charts in the ChartContainer.
// It iterates through all charts and returns the earliest timestamp found.
// If the ChartContainer is empty, it returns the zero value of time.Time.
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

// LastEnd finds the latest end time among all charts in the ChartContainer.
// It iterates through all charts and returns the latest timestamp found.
// If the ChartContainer is empty, it returns the zero value of time.Time.
func (c *ChartContainer) LastEnd() time.Time {
	var last time.Time

	for _, chart := range *c {
		end := chart.Timestamp.End()
		if last.IsZero() || end.After(last) {
			last = end
		}
	}

	return last
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

// Candles returns all the candles in the ChartContainer.
// It implements a merging stage of the merge-sort algorithm.
// Complexity: O(NK), where N is the number of candles and K is the number of instruments.
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

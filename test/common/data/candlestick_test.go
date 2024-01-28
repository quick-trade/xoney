package data_test

import (
	"testing"
	"time"
	"xoney/common/data"
	"xoney/errors"
)

func newTestTimeFrame() data.TimeFrame {
	tf, _ := data.NewTimeFrame(time.Minute*10, "10m")

	return *tf
}

func newChart() data.Chart {
	rawChart := data.RawChart(newTestTimeFrame(), 10)

	// Adding candles to the chart
	for i := 0; i < 5; i++ {
		candle := data.NewCandle(
			float64(i*10),
			float64(i*10+10),
			float64(i*10-5),
			float64(i*10+5),
			float64(i*100),
			time.Now().Add(time.Duration(i)*time.Minute),
		)
		rawChart.Add(*candle)
	}

	return rawChart
}

func TestTimeStampTimeframe(t *testing.T) {
	timeframe := newTestTimeFrame()
	ts := data.NewTimeStamp(timeframe, 10)

	result := ts.Timeframe()

	if result != timeframe {
		t.Errorf("Expected: %v, got: %v", timeframe, result)
	}
}

func TestTimeStampAt(t *testing.T) {
	timeframe := newTestTimeFrame()
	ts := data.NewTimeStamp(timeframe, 10)

	// Adding some timestamps
	ts.Append(time.Now(), time.Now().Add(time.Minute), time.Now().Add(2*time.Minute))

	// Testing At method
	expected := ts.Timestamp[1]
	result := ts.At(1)

	if result != expected {
		t.Errorf("Expected: %v, got: %v", expected, result)
	}
}

func TestTimeStampExtend(t *testing.T) {
	timeframe := newTestTimeFrame()
	ts := data.NewTimeStamp(timeframe, 10)

	// Adding some timestamps
	ts.Append(time.Now(), time.Now().Add(time.Minute), time.Now().Add(2*time.Minute))

	// Extending timestamps
	ts.Extend(3)

	// Testing the length
	expectedLength := 6
	resultLength := len(ts.Timestamp)

	if resultLength != expectedLength {
		t.Errorf("Expected length: %v, got: %v", expectedLength, resultLength)
	}

	// Testing the extended values
	expected := ts.Timestamp[5]
	result := ts.Timestamp[4].Add(time.Minute * 10)

	if result != expected {
		t.Errorf("Expected: %v, got: %v", expected, result)
	}
}

func TestTimeStampAppend(t *testing.T) {
	timeframe := newTestTimeFrame()
	ts := data.NewTimeStamp(timeframe, 10)

	// Adding some timestamps
	ts.Append(time.Now(), time.Now().Add(time.Minute), time.Now().Add(2*time.Minute))

	// Testing the length
	expectedLength := 3
	resultLength := len(ts.Timestamp)

	if resultLength != expectedLength {
		t.Errorf("Expected length: %v, got: %v", expectedLength, resultLength)
	}

	// Testing the appended values
	expected := ts.Timestamp[2]
	result := ts.Timestamp[2]

	if result != expected {
		t.Errorf("Expected: %v, got: %v", expected, result)
	}
}

func TestTimeStampSlice(t *testing.T) {
	timeframe := newTestTimeFrame()
	ts := data.NewTimeStamp(timeframe, 10)

	// Adding some timestamps
	ts.Append(time.Now(), time.Now().Add(time.Minute), time.Now().Add(2*time.Minute), time.Now().Add(3*time.Minute))

	// Slicing the timestamps
	sliced := ts.Slice(1, 3)

	// Testing the length
	expectedLength := 2
	resultLength := len(sliced.Timestamp)

	if resultLength != expectedLength {
		t.Errorf("Expected length: %v, got: %v", expectedLength, resultLength)
	}

	// Testing the sliced values
	expected := ts.Timestamp[1]
	result := sliced.Timestamp[0]

	if result != expected {
		t.Errorf("Expected: %v, got: %v", expected, result)
	}
}

func TestTimeStampEnd(t *testing.T) {
	timeframe := newTestTimeFrame()
	ts := data.NewTimeStamp(timeframe, 10)

	// Adding some timestamps
	ts.Append(time.Now(), time.Now().Add(time.Minute), time.Now().Add(2*time.Minute), time.Now().Add(3*time.Minute))

	// Testing the End method
	expected := ts.Timestamp[3]
	result := ts.End()

	if result != expected {
		t.Errorf("Expected: %v, got: %v", expected, result)
	}
}

func TestTimeStampStart(t *testing.T) {
	timeframe := newTestTimeFrame()
	ts := data.NewTimeStamp(timeframe, 10)

	// Adding some timestamps
	ts.Append(time.Now(), time.Now().Add(time.Minute), time.Now().Add(2*time.Minute), time.Now().Add(3*time.Minute))

	// Testing the Start method
	expected := ts.Timestamp[0]
	result := ts.Start()

	if result != expected {
		t.Errorf("Expected: %v, got: %v", expected, result)
	}
}

func TestTimeStampLen(t *testing.T) {
	timeframe := newTestTimeFrame()
	ts := data.NewTimeStamp(timeframe, 10)

	// Adding some timestamps
	ts.Append(time.Now(), time.Now().Add(time.Minute), time.Now().Add(2*time.Minute), time.Now().Add(3*time.Minute))

	// Testing the Len method
	expected := 4
	result := ts.Len()

	if result != expected {
		t.Errorf("Expected: %v, got: %v", expected, result)
	}
}

func TestChartAdd(t *testing.T) {
	// Creating a raw chart with a time frame of 1 minute and capacity 10
	rawChart := data.RawChart(newTestTimeFrame(), 10)

	// Adding a candle to the chart
	candle := data.NewCandle(50.0, 60.0, 45.0, 55.0, 100.0, time.Now())
	rawChart.Add(*candle)

	// Checking the length of the chart
	expectedLength := 1
	resultLength := len(rawChart.Timestamp.Timestamp)

	if resultLength != expectedLength {
		t.Errorf("Expected length: %v, got: %v", expectedLength, resultLength)
	}

	// Checking the values of the added candle
	expectedCandle := *candle
	resultCandle, err := rawChart.CandleByIndex(0)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if *resultCandle != expectedCandle {
		t.Errorf("Expected: %v, got: %v", expectedCandle, *resultCandle)
	}
}

func TestChartAddEmptyCandle(t *testing.T) {
	// Creating a raw chart with a time frame of 1 minute and capacity 10
	rawChart := data.RawChart(newTestTimeFrame(), 10)

	// Adding an empty candle to the chart
	emptyCandle := &data.Candle{}
	rawChart.Add(*emptyCandle)

	// Checking the length of the chart
	expectedLength := 1
	resultLength := len(rawChart.Timestamp.Timestamp)

	if resultLength != expectedLength {
		t.Errorf("Expected length: %v, got: %v", expectedLength, resultLength)
	}

	// Checking the values of the added candle
	expectedCandle := *emptyCandle
	resultCandle, err := rawChart.CandleByIndex(0)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if *resultCandle != expectedCandle {
		t.Errorf("Expected: %v, got: %v", expectedCandle, *resultCandle)
	}
}

func TestChartAddError(t *testing.T) {
	// Creating a raw chart with a time frame of 1 minute and capacity 10
	rawChart := data.RawChart(newTestTimeFrame(), 1)

	// Adding a candle to the chart
	candle := data.NewCandle(50.0, 60.0, 45.0, 55.0, 100.0, time.Now())
	rawChart.Add(*candle)

	// Adding another candle to exceed the capacity
	secondCandle := data.NewCandle(60.0, 70.0, 55.0, 65.0, 120.0, time.Now().Add(time.Minute))
	rawChart.Add(*secondCandle)

	// Checking the error for exceeding capacity
	_, err := rawChart.CandleByIndex(2)

	if _, ok := err.(errors.OutOfIndexError); !ok {
		t.Errorf("Expected OutOfIndexError, got: %v", err)
	}
}

func TestChartSlice(t *testing.T) {
	// Creating a raw chart with a time frame of 1 minute and capacity 10
	chart := data.RawChart(newTestTimeFrame(), 10)

	timeStart := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

	// Checking the result for raw chart
	sliced := chart.Slice(data.NewPeriod(timeStart, timeStart.Add(time.Hour*999)))
	if sliced.Len() != 0 {
		t.Error("Expected slice length to be 0 with raw source")
	}

	// Adding candles to the chart
	for i := 0; i < 5; i++ {
		candle := data.NewCandle(
			float64(i*10),
			float64(i*10+10),
			float64(i*10-5),
			float64(i*10+5),
			float64(100500),
			timeStart.Add(time.Minute*time.Duration(10*i)),
		)

		chart.Add(*candle)
	}

	// Slicing the chart
	sliced = chart.Slice(data.NewPeriod(
		timeStart.Add(11*time.Minute),
		timeStart.Add(30*time.Minute),
	))

	// Checking the length of the sliced chart
	expectedLength := 3
	resultLength := sliced.Timestamp.Len()

	if resultLength != expectedLength {
		t.Errorf("Expected length: %d, got: %d", expectedLength, resultLength)
	}

	// Checking the values of the sliced candles
	for i := 0; i < 3; i++ {
		expectedCandle, err := chart.CandleByIndex(i + 1)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		resultCandle, err := sliced.CandleByIndex(i)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if *resultCandle != *expectedCandle {
			t.Errorf("Expected: %v, got: %v", expectedCandle, resultCandle)
		}
	}

	sliced = chart.Slice(data.NewPeriod(
		timeStart,
		timeStart.Add(time.Hour*24), // outside of range
	))
	if sliced.Len() != chart.Len() {
		t.Errorf("Expected length: %d, got: %d", chart.Len(), sliced.Len())
	}
}

func TestChartSliceError(t *testing.T) {
	// Creating a raw chart with a time frame of 1 minute and capacity 10
	rawChart := newChart()

	// Slicing the chart with an incorrect period
	period := data.NewPeriod(
		time.Now().Add(3*time.Minute),
		time.Now().Add(2*time.Minute),
	)
	sliced := rawChart.Slice(period)

	// Checking the error for an incorrect period
	_, err := sliced.CandleByIndex(0)

	if _, ok := err.(errors.OutOfIndexError); !ok {
		t.Errorf("Expected OutOfIndexError, got: %v", err)
	}
}

func TestChartLen(t *testing.T) {
	// Creating a raw chart with a time frame of 1 minute and capacity 10
	rawChart := newChart()

	// Checking the length of the chart
	expectedLength := 5
	resultLength := rawChart.Len()

	if resultLength != expectedLength {
		t.Errorf("Expected length: %v, got: %v", expectedLength, resultLength)
	}
}

func TestChartSlicePeriodEndsBeforeStart(t *testing.T) {
	// Создание Chart с таймфреймом 1 минута и емкостью 10
	chart := newChart()

	// Создание периода, который заканчивается раньше, чем начинается Chart
	period := data.NewPeriod(
		time.Now().Add(-5*time.Minute),
		time.Now().Add(-2*time.Minute),
	)

	// Вызов метода Slice
	sliced := chart.Slice(period)

	// Проверка, что возвращается новый Chart с нулевой емкостью
	expectedCapacity := 0
	resultCapacity := cap(sliced.Close)

	if resultCapacity != expectedCapacity {
		t.Errorf("Expected capacity: %v, got: %v", expectedCapacity, resultCapacity)
	}
}

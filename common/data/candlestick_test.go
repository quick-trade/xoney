package data

import (
	"reflect"
	"testing"
	"time"
	"xoney/errors"
)

func timeFrameMinute() TimeFrame {
	tf, _ := NewTimeFrame(time.Minute, "1m")

	return *tf
}

func TestFindIndexBeforeOrAtTime(t *testing.T) {
	// Test case: Empty timestamp series
	emptySeries := NewTimeStamp(timeFrameMinute(), 0)
	timeStart := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

	_, err := findIndexBeforeOrAtTime(emptySeries, timeStart)
	if _, ok := err.(errors.ZeroLengthError); !ok {
		t.Errorf("Expected ZeroLengthError, got: %T", err)
	}

	// Test case: Moment is before the beginning of the series
	nonEmptySeries := NewTimeStamp(timeFrameMinute(), 3)
	nonEmptySeries.Append(timeStart, timeStart.Add(2*time.Minute), timeStart.Add(3*time.Minute))
	momentBefore := timeStart.Add(-time.Minute)
	_, err = findIndexBeforeOrAtTime(nonEmptySeries, momentBefore)

	if _, ok := err.(errors.ValueNotFoundError); !ok || err.Error() != "value not found." {
		t.Errorf("Expected ValueNotFoundError, got: %T", err)
	}

	// Test case: Valid scenario
	momentValid := timeStart.Add(time.Minute)
	index, err := findIndexBeforeOrAtTime(nonEmptySeries, momentValid)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedIndex := 0
	if index != expectedIndex {
		t.Errorf("Expected index: %d, got: %d", expectedIndex, index)
	}
}

func TestSortedInstruments(t *testing.T) {
	// Создание инструментов и добавление их в контейнер
	m1, _ := NewTimeFrame(time.Minute, "1m")
	h1, _ := NewTimeFrame(time.Hour, "1h")
	instrument1 := NewInstrument(*NewSymbol("BTC", "USD", "BINANCE"), *m1)
	instrument2 := NewInstrument(*NewSymbol("ETH", "USD", "BINANCE"), *h1)

	container := ChartContainer{
		instrument1: RawChart(instrument1.timeframe, 10),
		instrument2: RawChart(instrument2.timeframe, 10),
	}

	// Вызов приватного метода
	result := container.sortedInstruments()

	// Ожидаемый порядок инструментов по duration
	expectedOrder := []Instrument{instrument1, instrument2}

	// Сравнение результата с ожидаемым порядком
	if !reflect.DeepEqual(result, expectedOrder) {
		t.Errorf("Expected order: %v, got: %v", expectedOrder, result)
	}
}

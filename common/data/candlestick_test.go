package data

import (
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
    if _, ok := err.(errors.ValueNotFoundError); !ok {
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

package internal_test

import (
	"testing"

	"xoney/internal"
)

func array() []float64 {
	return []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
}

func TestMean(t *testing.T) {
	arr := array()
	mean, err := internal.RawMoment(arr, 1)
	if err != nil {
		t.Error(err.Error())
	}

	if mean != 5 {
		t.Errorf("expected mean to be 5, got %v", mean)
	}
}

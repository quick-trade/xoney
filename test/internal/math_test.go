package internal_test

import (
	"testing"

	"github.com/quick-trade/xoney/errors"
	"github.com/quick-trade/xoney/internal"
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

func TestZeroLen(t *testing.T) {
	arr := []float64{}
	_, err := internal.RawMoment(arr, 1)

	expected := errors.NewZeroLengthError("series")
	if err.Error() != expected.Error() {
		t.Errorf("there is no correct error, got %v", err)
	}
}

func TestDiff(t *testing.T) {
	arr := []float64{1, 2, 3, 13}

	diff := internal.Diff(arr)
	expected := []float64{1, 1, 10}

	if len(diff) != len(expected) {
		t.Error("incorrect diff length")
	}

	for i := range diff {
		if diff[i] != expected[i] {
			t.Errorf("incorrect diff: expected %v, got %v", expected, diff)

			break
		}
	}
}

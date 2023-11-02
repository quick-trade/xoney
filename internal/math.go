package internal

import (
	"math"

	"xoney/errors"
)

type data []float64

func RawMoment(sample data, degree float64) (float64, error) {
	n := len(sample)

	if n == 0 {
		return 0, errors.NewZeroLengthError("series")
	}

	var sum float64

	for _, v := range sample {
		sum += math.Pow(v, degree)
	}

	return sum / float64(n), nil
}

func CentralMoment(sample data, mean float64, degree float64) float64 {
	var moment float64

	n := len(sample)

	for _, v := range sample {
		moment += math.Pow(v-mean, degree)
	}

	return moment / float64(n)
}

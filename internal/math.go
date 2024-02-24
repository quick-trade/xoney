package internal

import (
	"math"

	"github.com/quick-trade/xoney/errors"
)

type Data []float64

func RawMoment(sample Data, degree float64) (float64, error) {
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

func CentralMoment(sample Data, mean float64, degree float64) float64 {
	var moment float64

	n := len(sample)

	for _, v := range sample {
		moment += math.Pow(v-mean, degree)
	}

	return moment / float64(n)
}

func Diff(sample Data) Data {
	diff := make(Data, 0, len(sample))
	for i := 1; i < len(sample); i++ {
		diff = Append(diff, sample[i]-sample[i-1])
	}

	return diff
}

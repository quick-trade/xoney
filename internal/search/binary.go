package search

import (
	"xoney/errors"
	math "xoney/internal"
)

func Index(series math.Data, moment float64) (int, error) {
	left, right := 0, len(series)-1

	if len(series) == 0 {
		return -1, errors.NewZeroLengthError("series")
	}

	if moment < series[left] {
		return -1, errors.ValueNotFoundError{}
	}

	var result int

	for left <= right {
		mid := (left + right) / 2

		if series[mid] < moment {
			result = mid
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return result, nil
}

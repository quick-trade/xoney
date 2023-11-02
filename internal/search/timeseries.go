package search

import (
	"time"

	"xoney/common"
	"xoney/errors"
)

func LastBeforeIdx(series common.TimeStamp, moment time.Time) (int, error) {
	left, right := 0, len(series)-1

	if len(series) == 0 {
		return -1, errors.NewZeroLengthError("timestamp series")
	}

	if moment.Before(series[left]) {
		return -1, errors.ValueNotFoundError{}
	}

	var result int

	for left <= right {
		mid := (left + right) / 2

		if series[mid].Before(moment) {
			result = mid
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return result, nil
}

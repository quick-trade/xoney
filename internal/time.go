package internal

import "time"

const Year = time.Hour * 24 * 365

func TimesInYear(duration time.Duration) float64 {
	return float64(Year) / float64(duration)
}

package internal

import "time"

const Year = time.Hour * 24 * 365

// TimesInYear calculates the number of intervals of the specified duration
// that fit within a standard year. For example, passing a duration of one day
// would return the number of days in a year.
//
// Parameters:
//
//	duration - A time.Duration value representing the length of the interval.
//
// Returns:
//
//	The number of times the interval fits into a year, as a float64.
func TimesInYear(duration time.Duration) float64 {
	return float64(Year) / float64(duration)
}

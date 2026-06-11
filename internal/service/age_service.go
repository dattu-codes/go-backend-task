package service

import (
	"time"
)

// CalculateAge computes the exact age in years based on a birthdate (dob) and a reference time (usually current time).
// It implements boundary checks for months and days, including leap-year scenarios (e.g. Feb 29).
func CalculateAge(dob time.Time, refTime time.Time) int {
	age := refTime.Year() - dob.Year()

	// Adjust age if the birth date has not been reached yet in the current year.
	// This compares the month and day of reference time against the birth date.
	if refTime.Month() < dob.Month() || (refTime.Month() == dob.Month() && refTime.Day() < dob.Day()) {
		age--
	}
	return age
}

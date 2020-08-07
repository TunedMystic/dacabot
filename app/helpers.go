package app

import (
	"log"
	"time"
)

// MustParseDate accepts a date string and returns a time.Time value.
func MustParseDate(dateString string) time.Time {
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		log.Fatalf("Could not convert %v to a date.\n", dateString)
	}
	return date
}

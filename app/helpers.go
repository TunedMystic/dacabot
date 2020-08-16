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

func TrimText(text string, truncLength int) string {
	if len(text) > truncLength {
		// Split string by rune.
		// Ref: https://stackoverflow.com/a/46416046
		// TODO: Not the best solution. Consider Adrian's approach in the SO answer.
		return string([]rune(text)[:truncLength-3]) + "..."
	}
	return text
}

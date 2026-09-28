package admin

import (
	"fmt"
	"time"
)

// parseProductOfferDates parses start and end date strings in YYYY-MM-DD format
// and ensures start <= end. If both strings are empty, returns zero-times and nil error.
func parseProductOfferDates(startStr, endStr string) (time.Time, time.Time, error) {
	const layout = "2006-01-02"

	if startStr == "" && endStr == "" {
		return time.Time{}, time.Time{}, nil
	}

	if startStr == "" || endStr == "" {
		return time.Time{}, time.Time{}, fmt.Errorf("both start_date and end_date must be provided")
	}

	start, err := time.Parse(layout, startStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid start_date format: %w", err)
	}

	end, err := time.Parse(layout, endStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid end_date format: %w", err)
	}

	if end.Before(start) {
		return time.Time{}, time.Time{}, fmt.Errorf("end_date must be same or after start_date")
	}

	return start, end, nil
}

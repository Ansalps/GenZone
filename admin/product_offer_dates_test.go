package admin

import (
	"testing"
	"time"
)

func TestParseProductOfferDates(t *testing.T) {
	start, end, err := parseProductOfferDates("2026-10-01", "2026-10-15")
	if err != nil {
		t.Fatalf("expected valid date range: %v", err)
	}

	if start.Year() != 2026 || start.Month() != time.October || start.Day() != 1 {
		t.Fatalf("unexpected start date: %v", start)
	}

	if end.Year() != 2026 || end.Month() != time.October || end.Day() != 15 {
		t.Fatalf("unexpected end date: %v", end)
	}
}

func TestParseProductOfferDatesRejectsInvalidRange(t *testing.T) {
	_, _, err := parseProductOfferDates("2026-10-15", "2026-10-01")
	if err == nil {
		t.Fatal("expected invalid date range error")
	}
}

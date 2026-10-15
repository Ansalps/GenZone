package admin

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func resolveReportDateRange(period, day, month, year string, now time.Time) (time.Time, time.Time, error) {
	if day != "" {
		parsed, err := time.Parse("2006-01-02", day)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid date format")
		}
		start := time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, parsed.Location())
		end := time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 23, 59, 59, 999999999, parsed.Location())
		return start, end, nil
	}

	if month != "" {
		monthNum, err := strconv.Atoi(month)
		if err != nil || monthNum < 1 || monthNum > 12 {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid month format")
		}
		parsedYear := now.Year()
		if year != "" {
			parsedYear, err = strconv.Atoi(year)
			if err != nil || parsedYear < 1900 || parsedYear > now.Year() {
				return time.Time{}, time.Time{}, fmt.Errorf("invalid year format")
			}
		}
		start := time.Date(parsedYear, time.Month(monthNum), 1, 0, 0, 0, 0, now.Location())
		end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)
		return start, end, nil
	}

	if year != "" {
		parsedYear, err := strconv.Atoi(year)
		if err != nil || parsedYear < 1900 || parsedYear > now.Year() {
			return time.Time{}, time.Time{}, fmt.Errorf("year out of range, must be between 1900 and current year")
		}
		start := time.Date(parsedYear, time.January, 1, 0, 0, 0, 0, now.Location())
		end := time.Date(parsedYear, time.December, 31, 23, 59, 59, 999999999, now.Location())
		return start, end, nil
	}

	normalizedPeriod := normalizeReportPeriod(period)
	if normalizedPeriod == "" {
		return time.Time{}, time.Time{}, nil
	}

	start, end := currentRangeForPeriod(normalizedPeriod, now)
	return start, end, nil
}

func normalizeReportPeriod(period string) string {
	switch strings.ToLower(strings.TrimSpace(period)) {
	case "", "day", "daily":
		return "daily"
	case "week", "weekly":
		return "weekly"
	case "2_weekly", "2 weekly", "2-weekly", "biweekly":
		return "2_weekly"
	case "month", "monthly":
		return "monthly"
	case "3_monthly", "3 monthly", "3-monthly", "quarterly":
		return "3_monthly"
	case "4_monthly", "4 monthly", "4-monthly":
		return "4_monthly"
	case "6_monthly", "6 monthly", "6-monthly":
		return "6_monthly"
	case "year", "yearly":
		return "yearly"
	default:
		return ""
	}
}

func currentRangeForPeriod(period string, now time.Time) (time.Time, time.Time) {
	loc := now.Location()
	switch period {
	case "daily":
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
		end := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, loc)
		return start, end
	case "weekly":
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc).AddDate(0, 0, -6)
		end := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, loc)
		return start, end
	case "2_weekly":
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc).AddDate(0, 0, -13)
		end := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, loc)
		return start, end
	case "monthly":
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
		end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)
		return start, end
	case "3_monthly":
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc).AddDate(0, -2, 0)
		end := time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 999999999, loc)
		return start, end
	case "4_monthly":
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc).AddDate(0, -3, 0)
		end := time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 999999999, loc)
		return start, end
	case "6_monthly":
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc).AddDate(0, -5, 0)
		end := time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 999999999, loc)
		return start, end
	case "yearly":
		start := time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, loc)
		end := time.Date(now.Year(), time.December, 31, 23, 59, 59, 999999999, loc)
		return start, end
	default:
		return time.Time{}, time.Time{}
	}
}

func buildReportDateFilter(columnName, period, day, month, year string, now time.Time) (string, []interface{}, error) {
	start, end, err := resolveReportDateRange(period, day, month, year, now)
	if err != nil {
		return "", nil, err
	}
	if start.IsZero() && end.IsZero() {
		return "", nil, nil
	}
	return " WHERE " + columnName + " >= ? AND " + columnName + " <= ? ", []interface{}{start, end}, nil
}

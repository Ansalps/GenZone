package admin

import (
	"testing"
	"time"
)

func TestResolveReportDateRange(t *testing.T) {
	now := time.Date(2026, time.June, 15, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		period   string
		day      string
		month    string
		year     string
		wantFrom time.Time
		wantTo   time.Time
	}{
		{name: "daily period", period: "daily", wantFrom: time.Date(2026, time.June, 15, 0, 0, 0, 0, time.UTC), wantTo: time.Date(2026, time.June, 15, 23, 59, 59, 999999999, time.UTC)},
		{name: "weekly period", period: "weekly", wantFrom: time.Date(2026, time.June, 9, 0, 0, 0, 0, time.UTC), wantTo: time.Date(2026, time.June, 15, 23, 59, 59, 999999999, time.UTC)},
		{name: "2 weekly period", period: "2_weekly", wantFrom: time.Date(2026, time.June, 2, 0, 0, 0, 0, time.UTC), wantTo: time.Date(2026, time.June, 15, 23, 59, 59, 999999999, time.UTC)},
		{name: "monthly period", period: "monthly", wantFrom: time.Date(2026, time.June, 1, 0, 0, 0, 0, time.UTC), wantTo: time.Date(2026, time.June, 30, 23, 59, 59, 999999999, time.UTC)},
		{name: "3 monthly period", period: "3_monthly", wantFrom: time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC), wantTo: time.Date(2026, time.June, 30, 23, 59, 59, 999999999, time.UTC)},
		{name: "4 monthly period", period: "4_monthly", wantFrom: time.Date(2026, time.March, 1, 0, 0, 0, 0, time.UTC), wantTo: time.Date(2026, time.June, 30, 23, 59, 59, 999999999, time.UTC)},
		{name: "6 monthly period", period: "6_monthly", wantFrom: time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC), wantTo: time.Date(2026, time.June, 30, 23, 59, 59, 999999999, time.UTC)},
		{name: "yearly period", period: "yearly", wantFrom: time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC), wantTo: time.Date(2026, time.December, 31, 23, 59, 59, 999999999, time.UTC)},
		{name: "explicit month override", month: "02", year: "2026", wantFrom: time.Date(2026, time.February, 1, 0, 0, 0, 0, time.UTC), wantTo: time.Date(2026, time.February, 28, 23, 59, 59, 999999999, time.UTC)},
	}

	for _, tt := range tests {
		gotFrom, gotTo, err := resolveReportDateRange(tt.period, tt.day, tt.month, tt.year, now)
		if err != nil {
			t.Fatalf("period=%s err=%v", tt.period, err)
		}
		if !gotFrom.Equal(tt.wantFrom) || !gotTo.Equal(tt.wantTo) {
			t.Fatalf("period=%s got=(%s,%s) want=(%s,%s)", tt.period, gotFrom.Format(time.RFC3339Nano), gotTo.Format(time.RFC3339Nano), tt.wantFrom.Format(time.RFC3339Nano), tt.wantTo.Format(time.RFC3339Nano))
		}
	}
}

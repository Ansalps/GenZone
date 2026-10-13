package user

import "testing"

func TestShouldCancelEntireOrder(t *testing.T) {
	tests := []struct {
		name     string
		statuses []string
		want     bool
	}{
		{name: "no remaining items", statuses: nil, want: true},
		{name: "all remaining items cancelled", statuses: []string{"cancelled", "cancelled", "returned"}, want: true},
		{name: "one active item remains", statuses: []string{"pending", "cancelled"}, want: false},
		{name: "mixed active and returned", statuses: []string{"pending", "returned"}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldCancelEntireOrder(tt.statuses); got != tt.want {
				t.Fatalf("shouldCancelEntireOrder(%v) = %v, want %v", tt.statuses, got, tt.want)
			}
		})
	}
}

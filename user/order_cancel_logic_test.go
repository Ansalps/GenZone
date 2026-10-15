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

func TestValidateCancelQuantity(t *testing.T) {
	tests := []struct {
		name            string
		itemQuantity    uint
		requestedAmount uint
		want            uint
		wantErr         bool
	}{
		{name: "no quantity provided uses full item", itemQuantity: 3, requestedAmount: 0, want: 3},
		{name: "partial cancellation allowed", itemQuantity: 5, requestedAmount: 2, want: 2},
		{name: "quantity over limit rejected", itemQuantity: 3, requestedAmount: 4, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateCancelQuantity(tt.itemQuantity, tt.requestedAmount)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("validateCancelQuantity(%d, %d) expected error, got nil", tt.itemQuantity, tt.requestedAmount)
				}
				return
			}
			if err != nil {
				t.Fatalf("validateCancelQuantity(%d, %d) unexpected error: %v", tt.itemQuantity, tt.requestedAmount, err)
			}
			if got != tt.want {
				t.Fatalf("validateCancelQuantity(%d, %d) = %d, want %d", tt.itemQuantity, tt.requestedAmount, got, tt.want)
			}
		})
	}
}

func TestCanCancelByDeliveryStatus(t *testing.T) {
	tests := []struct {
		name          string
		deliveryState string
		want          bool
	}{
		{name: "pending is cancellable", deliveryState: "pending", want: true},
		{name: "uppercase pending is cancellable", deliveryState: "PENDING", want: true},
		{name: "shipped is not cancellable", deliveryState: "shipped", want: false},
		{name: "delivered is not cancellable", deliveryState: "delivered", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := canCancelByDeliveryStatus(tt.deliveryState); got != tt.want {
				t.Fatalf("canCancelByDeliveryStatus(%q) = %v, want %v", tt.deliveryState, got, tt.want)
			}
		})
	}
}

func TestProportionalAmount(t *testing.T) {
	got := proportionalAmount(1000, 2, 5)
	want := 400.0
	if got != want {
		t.Fatalf("proportionalAmount(1000, 2, 5) = %v, want %v", got, want)
	}

	if got2 := proportionalAmount(500, 0, 5); got2 != 0 {
		t.Fatalf("proportionalAmount(500, 0, 5) = %v, want 0", got2)
	}
}

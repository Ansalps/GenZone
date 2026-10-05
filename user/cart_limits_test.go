package user

import "testing"

func TestValidateCartAddition(t *testing.T) {
	tests := []struct {
		name             string
		currentItemQty   uint
		requestedQty     uint
		currentCartTotal uint
		wantOK           bool
		wantMsg          string
	}{
		{
			name:             "within limits",
			currentItemQty:   2,
			requestedQty:     3,
			currentCartTotal: 10,
			wantOK:           true,
		},
		{
			name:             "item quantity cap reached",
			currentItemQty:   5,
			requestedQty:     1,
			currentCartTotal: 10,
			wantOK:           false,
			wantMsg:          "Maximum quantity allowed for this product is 5",
		},
		{
			name:             "requested quantity above per product cap",
			currentItemQty:   0,
			requestedQty:     6,
			currentCartTotal: 0,
			wantOK:           false,
			wantMsg:          "Maximum quantity allowed for each product is 5",
		},
		{
			name:             "cart total cap reached",
			currentItemQty:   0,
			requestedQty:     5,
			currentCartTotal: 21,
			wantOK:           false,
			wantMsg:          "Cart total quantity cannot exceed 25",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, msg := validateCartAddition(tt.currentItemQty, tt.requestedQty, tt.currentCartTotal)
			if ok != tt.wantOK {
				t.Fatalf("validateCartAddition() ok = %v, want %v", ok, tt.wantOK)
			}
			if msg != tt.wantMsg {
				t.Fatalf("validateCartAddition() msg = %q, want %q", msg, tt.wantMsg)
			}
		})
	}
}

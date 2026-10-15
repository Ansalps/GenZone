package admin

import "testing"

func TestCanUpdateDeliveryStatus(t *testing.T) {
	tests := []struct {
		name     string
		current  string
		next     string
		expected bool
	}{
		{name: "pending to shipped", current: "pending", next: "shipped", expected: true},
		{name: "pending to failed", current: "pending", next: "failed", expected: true},
		{name: "shipped to delivered", current: "shipped", next: "delivered", expected: true},
		{name: "shipped to failed", current: "shipped", next: "failed", expected: true},
		{name: "delivered cannot fail", current: "delivered", next: "failed", expected: false},
		{name: "failed cannot resume", current: "failed", next: "shipped", expected: false},
	}

	for _, tc := range tests {
		if got := canUpdateDeliveryStatus(tc.current, tc.next); got != tc.expected {
			t.Fatalf("canUpdateDeliveryStatus(%q, %q) = %v, want %v", tc.current, tc.next, got, tc.expected)
		}
	}
}

func TestCanUpdatePaymentStatus(t *testing.T) {
	tests := []struct {
		name     string
		current  string
		next     string
		expected bool
	}{
		{name: "pending to paid", current: "pending", next: "paid", expected: true},
		{name: "pending to failed", current: "pending", next: "failed", expected: true},
		{name: "paid cannot revert", current: "paid", next: "pending", expected: false},
		{name: "failed cannot resume", current: "failed", next: "paid", expected: false},
	}

	for _, tc := range tests {
		if got := canUpdatePaymentStatus(tc.current, tc.next); got != tc.expected {
			t.Fatalf("canUpdatePaymentStatus(%q, %q) = %v, want %v", tc.current, tc.next, got, tc.expected)
		}
	}
}

func TestShouldCompleteOrderOnDelivery(t *testing.T) {
	tests := []struct {
		name          string
		paymentMethod string
		expected      bool
	}{
		{name: "cod does not auto-complete", paymentMethod: "COD", expected: false},
		{name: "razorpay auto-completes", paymentMethod: "RazorPay", expected: true},
		{name: "wallet auto-completes", paymentMethod: "Wallet", expected: true},
	}

	for _, tc := range tests {
		if got := shouldCompleteOrderOnDelivery(tc.paymentMethod); got != tc.expected {
			t.Fatalf("shouldCompleteOrderOnDelivery(%q) = %v, want %v", tc.paymentMethod, got, tc.expected)
		}
	}
}

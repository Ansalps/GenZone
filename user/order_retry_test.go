package user

import "testing"

func TestNormalizeRetryPaymentMethod(t *testing.T) {
	if got := normalizeRetryPaymentMethod("razorpay"); got != "RazorPay" {
		t.Fatalf("expected RazorPay, got %q", got)
	}
	if got := normalizeRetryPaymentMethod("wallet"); got != "Wallet" {
		t.Fatalf("expected Wallet, got %q", got)
	}
	if got := normalizeRetryPaymentMethod("cod"); got != "COD" {
		t.Fatalf("expected COD, got %q", got)
	}
}

func TestBuildRetryPaymentUpdates(t *testing.T) {
	updates := buildRetryOrderUpdates("RazorPay", "rzp_123")
	if updates["payment_method"] != "RazorPay" {
		t.Fatalf("expected RazorPay payment_method, got %#v", updates["payment_method"])
	}
	if updates["razorpay_order_id"] != "rzp_123" {
		t.Fatalf("expected razorpay_order_id to be set, got %#v", updates["razorpay_order_id"])
	}

	walletUpdates := buildRetryOrderUpdates("Wallet", "")
	if walletUpdates["payment_method"] != "Wallet" {
		t.Fatalf("expected Wallet payment_method, got %#v", walletUpdates["payment_method"])
	}
}

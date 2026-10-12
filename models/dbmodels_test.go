package models

import "testing"

func TestOrderHasOptionalCouponID(t *testing.T) {
	order := Order{
		UserID:      1,
		AddressID:   1,
		TotalAmount: 100,
		OrderStatus: "pending",
		FinalAmount: 100,
	}

	if order.CouponID != nil {
		t.Fatal("expected coupon id to be nil when no coupon is applied")
	}
}

func TestOrderTracksCouponCodeWhenApplied(t *testing.T) {
	order := Order{
		UserID:      1,
		AddressID:   1,
		TotalAmount: 100,
		OrderStatus: "pending",
		FinalAmount: 100,
		CouponCode:  "SAVE10",
	}

	if order.CouponCode != "SAVE10" {
		t.Fatal("expected coupon code to be retained on the order")
	}
}

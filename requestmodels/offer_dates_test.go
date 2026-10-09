package requestmodels

import "testing"

func TestOfferIncludesDateRange(t *testing.T) {
	req := Offer{
		ProductID:          2,
		DiscountPercentage: 12,
		StartAt:            "2026-01-10",
		EndAt:              "2026-01-20",
	}

	if req.StartAt == "" {
		t.Fatal("expected start_at to be present")
	}

	if req.EndAt == "" {
		t.Fatal("expected end_at to be present")
	}
}

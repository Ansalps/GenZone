package responsemodels

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestOrderItemsQuantityIsIncludedInJSON(t *testing.T) {
	item := OrderItems{
		ID:          1,
		OrderID:     "10",
		ProductID:   "20",
		ProductName: "Widget",
		Quantity:    3,
		Price:       49.99,
	}

	payload, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	if !strings.Contains(string(payload), "\"quantity\":3") {
		t.Fatalf("expected JSON to include quantity field, got %s", payload)
	}
}

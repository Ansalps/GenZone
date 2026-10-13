package responsemodels

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestOrderItemsStatusJSONKey(t *testing.T) {
	item := OrderItems{OrderStatus: "cancelled", DeliveryStatus: "shipped"}
	data, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	jsonText := string(data)
	if jsonText == "" || jsonText == "{}" {
		t.Fatalf("expected JSON payload, got %s", jsonText)
	}
	if strings.Contains(jsonText, "order_staus") {
		t.Fatalf("status key should be order_status, not order_staus: %s", jsonText)
	}
	if !strings.Contains(jsonText, "\"order_status\":\"cancelled\"") {
		t.Fatalf("expected order_status field in JSON, got %s", jsonText)
	}
	if !strings.Contains(jsonText, "\"delivery_status\":\"shipped\"") {
		t.Fatalf("expected delivery_status field in JSON, got %s", jsonText)
	}
}

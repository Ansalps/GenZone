package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/middleware"
	"github.com/gin-gonic/gin"
)

func TestTotalQuantity_EmptyCartReturnsZero(t *testing.T) {
	gin.SetMode(gin.TestMode)

	if database.DB == nil {
		database.Initialize()
		if err := database.AutoMigrate(); err != nil {
			t.Fatalf("AutoMigrate failed: %v", err)
		}
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/cart-total-quantity", nil)
	c.Set("claims", &middleware.CustomClaims{ID: 999999})

	TotalQuantity(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d with body: %s", w.Code, w.Body.String())
	}

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}

	if body["status"] != true {
		t.Fatalf("expected status true, got %#v", body["status"])
	}

	if total, ok := body["total_quantity"]; !ok || total != float64(0) {
		t.Fatalf("expected total_quantity to be 0, got %#v", body["total_quantity"])
	}
}

package admin

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/models"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestReadProductsIncludesDeletedProductsAndStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}
	database.DB = db

	if err := db.AutoMigrate(&models.Category{}, &models.Product{}, &models.ProductVariant{}, &models.Offer{}); err != nil {
		t.Fatalf("failed to migrate models: %v", err)
	}

	category := models.Category{CategoryName: "Test Category", Description: "Desc", ImageURL: "https://example.com/cat.jpg"}
	if err := db.Create(&category).Error; err != nil {
		t.Fatalf("failed to create category: %v", err)
	}

	active := models.Product{
		CategoryID:  category.ID,
		ProductName: "Active Product",
		Description: "Should remain visible",
		ImageURL:    "https://example.com/active.jpg",
		Price:       10.99,
	}
	if err := db.Create(&active).Error; err != nil {
		t.Fatalf("failed to create active product: %v", err)
	}

	deleted := models.Product{
		CategoryID:  category.ID,
		ProductName: "Deleted Product",
		Description: "Should appear as discontinued",
		ImageURL:    "https://example.com/deleted.jpg",
		Price:       19.99,
	}
	if err := db.Create(&deleted).Error; err != nil {
		t.Fatalf("failed to create deleted product: %v", err)
	}
	if err := db.Delete(&deleted).Error; err != nil {
		t.Fatalf("failed to soft delete product: %v", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodGet, "/admin/product", nil)
	c.Request = req

	ReadProducts(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d with body %s", w.Code, w.Body.String())
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to parse response JSON: %v", err)
	}

	data, ok := payload["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("missing data field: %v", payload)
	}
	products, ok := data["products"].([]interface{})
	if !ok {
		t.Fatalf("expected products array, got %#v", data["products"])
	}
	if len(products) != 2 {
		t.Fatalf("expected 2 products (active + discontinued), got %d: %#v", len(products), data["products"])
	}

	statuses := make(map[string]bool)
	for _, item := range products {
		p, ok := item.(map[string]interface{})
		if !ok {
			t.Fatalf("product entry was not an object: %#v", item)
		}
		if status, ok := p["product_status"].(string); ok {
			statuses[status] = true
		}
	}

	if !statuses["active"] || !statuses["discontinued"] {
		t.Fatalf("expected both active and discontinued statuses, got %#v", statuses)
	}
}

package public

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

func TestReadProductsIncludesDiscontinuedProducts(t *testing.T) {
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

	product := models.Product{
		CategoryID:  category.ID,
		ProductName: "Discontinued Item",
		Description: "Soft deleted",
		ImageURL:    "https://example.com/product.jpg",
		Price:       99.99,
	}
	if err := db.Create(&product).Error; err != nil {
		t.Fatalf("failed to create product: %v", err)
	}
	if err := db.Delete(&product).Error; err != nil {
		t.Fatalf("failed to soft delete product: %v", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodGet, "/public/product", nil)
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
	if !ok || len(products) != 1 {
		t.Fatalf("expected 1 discontinued product in response, got %#v", data["products"])
	}
	productData, ok := products[0].(map[string]interface{})
	if !ok {
		t.Fatalf("product entry was not an object: %#v", products[0])
	}
	if productData["product_status"] != "discontinued" {
		t.Fatalf("expected discontinued product_status, got %#v", productData["product_status"])
	}
}

func TestReadProductsMarksOnlyDeletedProductsAsDiscontinued(t *testing.T) {
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

	product := models.Product{
		CategoryID:  category.ID,
		ProductName: "Active Product",
		Description: "Still active",
		ImageURL:    "https://example.com/active.jpg",
		Price:       79.99,
	}
	if err := db.Create(&product).Error; err != nil {
		t.Fatalf("failed to create active product: %v", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodGet, "/public/product", nil)
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
	if !ok || len(products) != 1 {
		t.Fatalf("expected 1 active product in response, got %#v", data["products"])
	}
	productData, ok := products[0].(map[string]interface{})
	if !ok {
		t.Fatalf("product entry was not an object: %#v", products[0])
	}
	if _, exists := productData["deleted_at"]; !exists {
		t.Fatalf("expected deleted_at field to be present for null-value response, got %#v", productData)
	}
	if productData["deleted_at"] != nil {
		t.Fatalf("expected deleted_at to be null for active product, got %#v", productData["deleted_at"])
	}
	if productData["product_status"] != "active" {
		t.Fatalf("expected active product_status, got %#v", productData["product_status"])
	}
}

func TestGetProductIncludesDiscontinuedProducts(t *testing.T) {
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

	product := models.Product{
		CategoryID:  category.ID,
		ProductName: "Discontinued Search Product",
		Description: "Soft deleted by search endpoint",
		ImageURL:    "https://example.com/product-search.jpg",
		Price:       149.99,
	}
	if err := db.Create(&product).Error; err != nil {
		t.Fatalf("failed to create product: %v", err)
	}
	if err := db.Delete(&product).Error; err != nil {
		t.Fatalf("failed to soft delete product for search endpoint: %v", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodGet, "/public/product", nil)
	c.Request = req

	GetProduct(c)

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
	if !ok || len(products) != 1 {
		t.Fatalf("expected 1 discontinued product in search response, got %#v", data["products"])
	}
	productData, ok := products[0].(map[string]interface{})
	if !ok {
		t.Fatalf("product entry was not an object: %#v", products[0])
	}
	if productData["product_status"] != "discontinued" {
		t.Fatalf("expected discontinued product_status from search endpoint, got %#v", productData["product_status"])
	}
}

package user

import (
	"testing"
	"time"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestIsRetryWindowOpen(t *testing.T) {
	now := time.Now()
	if !isRetryWindowOpen(now.Add(-30*time.Second), now) {
		t.Fatal("expected retry window to remain open within 1 minute")
	}
	if isRetryWindowOpen(now.Add(-2*time.Minute), now) {
		t.Fatal("expected retry window to close after 1 minute")
	}
}

func TestShouldExpirePayment(t *testing.T) {
	if !shouldExpirePayment(time.Now().Add(-2*time.Minute), time.Now()) {
		t.Fatal("expected payment to expire after 1 minute")
	}
	if shouldExpirePayment(time.Now().Add(-30*time.Second), time.Now()) {
		t.Fatal("expected payment to stay active before 1 minute")
	}
}

func TestReduceAndRestoreVariantStock(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}
	database.DB = db
	if err := db.AutoMigrate(&models.Product{}, &models.ProductVariant{}); err != nil {
		t.Fatalf("failed to migrate models: %v", err)
	}

	product := models.Product{
		ProductName: "Test Product",
		Description: "test description",
		ImageURL:    "https://example.com/test.jpg",
		Price:       100,
	}
	if err := db.Create(&product).Error; err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	variants := []models.ProductVariant{
		{ProductID: product.ID, Size: "S", Stock: 5},
		{ProductID: product.ID, Size: "M", Stock: 4},
	}
	if err := db.Create(&variants).Error; err != nil {
		t.Fatalf("failed to create variants: %v", err)
	}

	if err := reduceVariantStock(product.ID, 7); err != nil {
		t.Fatalf("reduceVariantStock returned error: %v", err)
	}

	var updated []models.ProductVariant
	if err := db.Where("product_id = ?", product.ID).Order("id ASC").Find(&updated).Error; err != nil {
		t.Fatalf("failed to fetch updated variants: %v", err)
	}
	if len(updated) != 2 {
		t.Fatalf("expected 2 variants after reduction, got %d", len(updated))
	}
	if updated[0].Stock != 0 || updated[1].Stock != 2 {
		t.Fatalf("expected reduced stock to be [0, 2], got [%d, %d]", updated[0].Stock, updated[1].Stock)
	}

	if err := restoreVariantStock(product.ID, 7); err != nil {
		t.Fatalf("restoreVariantStock returned error: %v", err)
	}
	if err := db.Where("product_id = ?", product.ID).Order("id ASC").Find(&updated).Error; err != nil {
		t.Fatalf("failed to fetch updated variants after restore: %v", err)
	}
	if updated[0].Stock != 7 || updated[1].Stock != 2 {
		t.Fatalf("expected restored stock to be [7, 2] while preserving total inventory, got [%d, %d]", updated[0].Stock, updated[1].Stock)
	}
}

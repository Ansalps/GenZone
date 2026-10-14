package user

import (
	"testing"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestOrderItemVariantTracking(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}
	database.DB = db

	if err := db.AutoMigrate(&models.Product{}, &models.ProductVariant{}, &models.Order{}, &models.OrderItems{}); err != nil {
		t.Fatalf("failed to migrate models: %v", err)
	}

	product := models.Product{
		ProductName: "Variant Product",
		Description: "For tracking test",
		ImageURL:    "https://example.com/product.jpg",
		Price:       200,
	}
	if err := db.Create(&product).Error; err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	var variant models.ProductVariant
	if err := db.Create(&models.ProductVariant{ProductID: product.ID, Size: "M", Stock: 5}).Error; err != nil {
		t.Fatalf("failed to create product variant: %v", err)
	}
	variant = models.ProductVariant{ProductID: product.ID, Size: "M", Stock: 5}
	if err := db.Where("product_id = ? AND size = ?", product.ID, "M").First(&variant).Error; err != nil {
		t.Fatalf("failed to fetch created variant: %v", err)
	}

	order := models.Order{UserID: 1, AddressID: 1, TotalAmount: 200, PaymentMethod: "COD", FinalAmount: 200}
	if err := db.Create(&order).Error; err != nil {
		t.Fatalf("failed to create order: %v", err)
	}

	item := models.OrderItems{
		OrderID:          order.ID,
		ProductID:        product.ID,
		ProductVariantID: variant.ID,
		Size:             "M",
		Quantity:         2,
		UnitPrice:        100,
		TotalAmount:      200,
		OrderStatus:      "pending",
		DeliveryStatus:   "pending",
		PaymentMethod:    "COD",
	}
	if err := db.Create(&item).Error; err != nil {
		t.Fatalf("failed to create order item: %v", err)
	}

	if err := reduceStockForOrder(order.ID); err != nil {
		t.Fatalf("reduceStockForOrder returned error: %v", err)
	}

	var updated models.ProductVariant
	if err := db.First(&updated, variant.ID).Error; err != nil {
		t.Fatalf("failed to fetch variant after stock deduction: %v", err)
	}
	if updated.Stock != 3 {
		t.Fatalf("expected reduced stock for M variant to be 3, got %d", updated.Stock)
	}
}

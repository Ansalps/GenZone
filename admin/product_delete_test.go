package admin

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/models"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestProductDeleteRemovesRelatedDatabaseRecords(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}
	database.DB = db

	if err := db.AutoMigrate(
		&models.Category{},
		&models.Product{},
		&models.ProductVariant{},
		&models.Offer{},
		&models.Cart{},
		&models.CartItem{},
		&models.Wishlist{},
		&models.Order{},
		&models.OrderItems{},
	); err != nil {
		t.Fatalf("failed to migrate models: %v", err)
	}

	category := models.Category{
		CategoryName: "Test Category",
		Description:  "Test description",
		ImageURL:     "https://example.com/category.jpg",
	}
	if err := db.Create(&category).Error; err != nil {
		t.Fatalf("failed to create category: %v", err)
	}

	product := models.Product{
		CategoryID:  category.ID,
		ProductName: "Delete Me",
		Description: "Will be removed",
		ImageURL:    "https://example.com/product.jpg",
		Price:       99.99,
	}
	if err := db.Create(&product).Error; err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	if err := db.Create(&models.ProductVariant{ProductID: product.ID, Size: "M", Stock: 10}).Error; err != nil {
		t.Fatalf("failed to create variant: %v", err)
	}
	if err := db.Create(&models.Offer{ProductID: product.ID, DiscountPercentage: 10}).Error; err != nil {
		t.Fatalf("failed to create offer: %v", err)
	}

	cart := models.Cart{UserID: 10}
	if err := db.Create(&cart).Error; err != nil {
		t.Fatalf("failed to create cart: %v", err)
	}
	if err := db.Create(&models.CartItem{CartID: cart.ID, ProductID: product.ID, Size: "M", Quantity: 1, UnitPrice: 99.99}).Error; err != nil {
		t.Fatalf("failed to create cart item: %v", err)
	}
	if err := db.Create(&models.Wishlist{UserID: 10, ProductID: product.ID, ProductName: product.ProductName}).Error; err != nil {
		t.Fatalf("failed to create wishlist: %v", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodDelete, "/admin/product/"+fmt.Sprintf("%d", product.ID), nil)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", product.ID)}}

	ProductDelete(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 status on delete, got %d and body %s", w.Code, w.Body.String())
	}

	var productCount int64
	if err := db.Model(&models.Product{}).Where("id = ?", product.ID).Count(&productCount).Error; err != nil {
		t.Fatalf("failed to query product after delete: %v", err)
	}
	if productCount != 0 {
		t.Fatalf("expected product row to be deleted, found %d rows", productCount)
	}

	var variantCount int64
	if err := db.Model(&models.ProductVariant{}).Where("product_id = ?", product.ID).Count(&variantCount).Error; err != nil {
		t.Fatalf("failed to query variants after delete: %v", err)
	}
	if variantCount != 0 {
		t.Fatalf("expected product variants to be deleted, found %d rows", variantCount)
	}

	var offerCount int64
	if err := db.Model(&models.Offer{}).Where("product_id = ?", product.ID).Count(&offerCount).Error; err != nil {
		t.Fatalf("failed to query offers after delete: %v", err)
	}
	if offerCount != 0 {
		t.Fatalf("expected product offers to be deleted, found %d rows", offerCount)
	}

	var cartItemCount int64
	if err := db.Model(&models.CartItem{}).Where("product_id = ?", product.ID).Count(&cartItemCount).Error; err != nil {
		t.Fatalf("failed to query cart items after delete: %v", err)
	}
	if cartItemCount != 0 {
		t.Fatalf("expected cart items to be removed, found %d rows", cartItemCount)
	}

	var wishlistCount int64
	if err := db.Model(&models.Wishlist{}).Where("product_id = ?", product.ID).Count(&wishlistCount).Error; err != nil {
		t.Fatalf("failed to query wishlist after delete: %v", err)
	}
	if wishlistCount != 0 {
		t.Fatalf("expected wishlist rows to be removed, found %d rows", wishlistCount)
	}
}

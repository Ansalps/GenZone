package helper

import (
	"testing"
	"time"

	"github.com/Ansalps/GeZOne/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGetActiveOfferDiscount(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite database: %v", err)
	}

	if err := db.AutoMigrate(&models.Offer{}); err != nil {
		t.Fatalf("failed to migrate offer table: %v", err)
	}

	now := time.Now()
	offer := models.Offer{
		ProductID:          1,
		DiscountPercentage: 12,
		StartAt:            now.Add(-time.Hour),
		EndAt:              now.Add(time.Hour),
	}
	if err := db.Create(&offer).Error; err != nil {
		t.Fatalf("failed to create offer: %v", err)
	}

	discount, ok := GetActiveOfferDiscount(db, 1)
	if !ok {
		t.Fatal("expected an active offer to be found")
	}
	if discount != 12 {
		t.Fatalf("expected discount 12, got %.2f", discount)
	}
}

func TestGetActiveOfferForProduct(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite database: %v", err)
	}

	if err := db.AutoMigrate(&models.Offer{}); err != nil {
		t.Fatalf("failed to migrate offer table: %v", err)
	}

	now := time.Now()
	offer := models.Offer{
		ProductID:          42,
		DiscountPercentage: 18,
		StartAt:            now.Add(-time.Hour),
		EndAt:              now.Add(time.Hour),
	}
	if err := db.Create(&offer).Error; err != nil {
		t.Fatalf("failed to create offer: %v", err)
	}

	activeOffer, ok := GetActiveOfferForProduct(db, 42)
	if !ok {
		t.Fatal("expected an active offer record to be found")
	}
	if activeOffer == nil {
		t.Fatal("expected active offer pointer to be non-nil")
	}
	if activeOffer.ID != offer.ID {
		t.Fatalf("expected offer id %d, got %d", offer.ID, activeOffer.ID)
	}
}

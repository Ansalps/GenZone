package helper

import (
	"time"

	"github.com/Ansalps/GeZOne/models"
	"gorm.io/gorm"
)

func GetActiveOfferForProduct(db *gorm.DB, productID uint) (*models.Offer, bool) {
	var offer models.Offer
	now := time.Now()

	if err := db.Where("product_id = ? AND start_at <= ? AND (end_at IS NULL OR end_at >= ?)", productID, now, now).
		Order("start_at DESC").
		First(&offer).Error; err != nil {
		return nil, false
	}

	return &offer, true
}

func GetActiveOfferDiscount(db *gorm.DB, productID uint) (float64, bool) {
	offer, ok := GetActiveOfferForProduct(db, productID)
	if !ok || offer == nil {
		return 0, false
	}

	return offer.DiscountPercentage, true
}

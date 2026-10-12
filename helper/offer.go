package helper

import (
	"time"

	"github.com/Ansalps/GeZOne/models"
	"gorm.io/gorm"
)

func GetActiveOfferDiscount(db *gorm.DB, productID uint) (float64, bool) {
	var offer models.Offer
	now := time.Now()

	if err := db.Where("product_id = ? AND start_at <= ? AND (end_at IS NULL OR end_at >= ?)", productID, now, now).
		Order("start_at DESC").
		First(&offer).Error; err != nil {
		return 0, false
	}

	return offer.DiscountPercentage, true
}

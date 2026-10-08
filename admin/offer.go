package admin

import (
	"fmt"
	"net/http"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/helper"
	"github.com/Ansalps/GeZOne/models"
	requestmodemodels "github.com/Ansalps/GeZOne/requestmodels"
	"github.com/Ansalps/GeZOne/responsemodels"
	"github.com/gin-gonic/gin"
)

func OfferList(c *gin.Context) {
	var offers []responsemodels.Offer
	if err := database.DB.Raw(`
		SELECT
			offers.id,
			offers.created_at,
			offers.updated_at,
			offers.deleted_at,
			offers.product_id,
			offers.discount_percentage,
			products.product_name,
			categories.category_name,
			products.description,
			products.image_url,
			products.price,
			COALESCE(SUM(product_variants.stock), 0) AS stock,
			products.popular,
			COALESCE(string_agg(DISTINCT product_variants.size, ',' ORDER BY product_variants.size), '') AS size
		FROM offers
		JOIN products ON offers.product_id = products.id
		JOIN categories ON categories.id = products.category_id
		LEFT JOIN product_variants ON product_variants.product_id = products.id
		WHERE offers.deleted_at IS NULL
		GROUP BY offers.id, offers.created_at, offers.updated_at, offers.deleted_at,
			offers.product_id, offers.discount_percentage,
			products.product_name, categories.category_name, products.description,
			products.image_url, products.price, products.popular
	`).Scan(&offers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "failed to fetch offers",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "listing offers successfully",
		"data": gin.H{
			"offers": offers,
		},
	})
}

func OfferAdd(c *gin.Context) {
	var req requestmodemodels.Offer
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "failed to bind request",
		})
		return
	}

	if err := helper.Validate(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"message":    err.Error(),
			"error_code": http.StatusBadRequest,
		})
		return
	}

	if req.DiscountPercentage < 0 || req.DiscountPercentage > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "discount percentage must be between 0 and 100",
		})
		return
	}

	var product models.Product
	if err := database.DB.First(&product, req.ProductID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "product does not exist",
		})
		return
	}

	var count int64
	if err := database.DB.Model(&models.Offer{}).Where("product_id = ? AND deleted_at IS NULL", req.ProductID).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "failed to check existing offer",
		})
		return
	}
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "product offer already exists, a product cannot have more than 1 offer",
		})
		return
	}

	offer := models.Offer{
		ProductID:          req.ProductID,
		DiscountPercentage: req.DiscountPercentage,
	}
	if err := database.DB.Create(&offer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "failed to add offer",
		})
		return
	}

	if err := database.DB.Model(&models.Product{}).Where("id = ?", req.ProductID).Updates(map[string]interface{}{
		"has_offer":              true,
		"offer_discount_percent": req.DiscountPercentage,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "offer created but product discount update failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "offer added for the product",
	})
}

func OfferEdit(c *gin.Context) {
	offerID := c.Param("id")
	var existing models.Offer
	if err := database.DB.First(&existing, offerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "offer id does not exist",
		})
		return
	}

	oldProductID := existing.ProductID

	var req requestmodemodels.Offer
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "failed to bind request",
		})
		return
	}

	if err := helper.Validate(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"message":    err.Error(),
			"error_code": http.StatusBadRequest,
		})
		return
	}

	if req.DiscountPercentage < 0 || req.DiscountPercentage > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "discount percentage must be between 0 and 100",
		})
		return
	}

	productID := req.ProductID
	if productID == 0 {
		productID = existing.ProductID
	}

	var product models.Product
	if err := database.DB.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "product does not exist",
		})
		return
	}

	if oldProductID != productID {
		var another int64
		if err := database.DB.Model(&models.Offer{}).Where("product_id = ? AND id != ? AND deleted_at IS NULL", productID, existing.ID).Count(&another).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "failed to verify offer target product",
			})
			return
		}
		if another > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "another offer already exists for this product",
			})
			return
		}
	}

	existing.ProductID = productID
	existing.DiscountPercentage = req.DiscountPercentage
	if err := database.DB.Save(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "failed to update offer",
		})
		return
	}

	if err := database.DB.Model(&models.Product{}).Where("id = ?", productID).Updates(map[string]interface{}{
		"has_offer":              true,
		"offer_discount_percent": existing.DiscountPercentage,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "offer updated but product discount status could not be refreshed",
		})
		return
	}

	if oldProductID != productID {
		_ = database.DB.Model(&models.Product{}).Where("id = ?", oldProductID).Updates(map[string]interface{}{
			"has_offer":              false,
			"offer_discount_percent": 0,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "offer updated successfully",
	})
}

func OfferRemove(c *gin.Context) {
	offerID := c.Param("id")
	var offer models.Offer
	if err := database.DB.First(&offer, offerID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "offer id does not exist",
		})
		return
	}

	if err := database.DB.Delete(&offer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "failed to delete offer",
		})
		return
	}

	if err := database.DB.Model(&models.Product{}).Where("id = ?", offer.ProductID).Updates(map[string]interface{}{
		"has_offer":              false,
		"offer_discount_percent": 0,
	}).Error; err != nil {
		fmt.Println("failed to clear offer state for product:", err)
	}

	c.JSON(http.StatusOK, gin.H{"status": true, "message": "offer deleted successfully"})
}

package admin

import (
	"fmt"
	"net/http"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/helper"
	"github.com/Ansalps/GeZOne/models"
	"github.com/gin-gonic/gin"
)

func OfferList(c *gin.Context) {
	var offer []models.Offer
	database.DB.Raw(`SELECT * FROM offers WHERE deleted_at IS NULL`).Scan(&offer)
	c.JSON(http.StatusOK, gin.H{
		"data":    offer,
		"message": "listing offers successfully",
	})
}

func OfferAdd(c *gin.Context) {
	var offeradd models.OfferAdd
	err := c.BindJSON(&offeradd)
	response := gin.H{
		"status":  false,
		"message": "failed to bind request",
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// Validate the content of the JSON
	if err := helper.Validate(offeradd); err != nil {
		fmt.Println("", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"message":    err.Error(),
			"error_code": http.StatusBadRequest,
		})
		return
	}
	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM products WHERE id = ?`, offeradd.ProductID).Scan(&count)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "product does not exist",
		})
	}
	var count1 int64
	database.DB.Raw(`SELECT COUNT(*) FROM offers WHERE product_id = ?`, offeradd.ProductID).Scan(&count1)
	if count1 > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "productoffer already exists, a product cannot have more than 1 offer",
		})
	}
	offer := models.Offer{
		ProductID:          offeradd.ProductID,
		DiscountPercentage: offeradd.DiscountPercentage,
	}
	database.DB.Create(&offer)
	database.DB.Model(&models.Product{}).Where("id = ?", offeradd.ProductID).Update("has_offer", true)
	c.JSON(http.StatusOK, gin.H{
		"message": "offer added for the product",
	})
}

func OfferRemove(c *gin.Context) {
	OfferID := c.Param("id")
	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM offers WHERE id = ? AND deleted_at IS NULL`, OfferID).Scan(&count)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "offer id does not exist",
		})
		return
	}
	var productid uint
	database.DB.Model(&models.Offer{}).Where("id = ?", OfferID).Pluck("product_id", &productid)
	database.DB.Where("id = ?", OfferID).Delete(&models.Offer{})
	database.DB.Model(&models.Product{}).Where("id = ?", productid).Update("has_offer", false)
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "offer deleted succesfully"})
}

package admin

import (
	"fmt"
	"net/http"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/helper"
	"github.com/Ansalps/GeZOne/models"
	"github.com/Ansalps/GeZOne/responsemodels"
	"github.com/gin-gonic/gin"
)

func OfferList(c *gin.Context) {
	var offer []responsemodels.Offer
	database.DB.Raw(`SELECT offers.id,offers.created_at,offers.updated_at,offers.deleted_at,offers.product_id,offers.discount_percentage,products.product_name,categories.category_name,products.description,products.image_url,products.price,products.stock,products.popular,products.size,products.has_offer,products.offer_discount_percent FROM offers JOIN products ON offers.product_id = products.id JOIN categories ON categories.id = products.category_id WHERE offers.deleted_at IS NULL`).Scan(&offer)

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
	database.DB.Raw(`SELECT COUNT(*) FROM offers WHERE product_id = ? and deleted_at IS NULL`, offeradd.ProductID).Scan(&count1)
	if count1 > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "productoffer already exists, a product cannot have more than 1 offer",
		})
		return
	}
	offer := models.Offer{
		ProductID:          offeradd.ProductID,
		DiscountPercentage: offeradd.DiscountPercentage,
	}
	database.DB.Create(&offer)
	database.DB.Model(&models.Product{}).Where("id = ?", offeradd.ProductID).Update("has_offer", true)
	database.DB.Model(&models.Product{}).Where("id = ?", offeradd.ProductID).Update("offer_discount_percent", offeradd.DiscountPercentage)
	var qty uint
	database.DB.Model(&models.CartItems{}).Where("product_id = ?", offeradd.ProductID).Pluck("qty", &qty)
	var price float64
	database.DB.Model(&models.CartItems{}).Where("product_id = ?", offeradd.ProductID).Pluck("price", &price)
	offerdiscount := price * float64(offeradd.DiscountPercentage) / 100
	database.DB.Model(&models.CartItems{}).Where("product_id = ?", offeradd.ProductID).Update("discount", float64(qty)*offerdiscount)
	var totalamount float64
	database.DB.Model(&models.CartItems{}).Where("product_id = ?", offeradd.ProductID).Pluck("total_amount", &totalamount)
	database.DB.Model(&models.CartItems{}).Where("product_id = ?", offeradd.ProductID).Update("final_amount", totalamount-float64(qty)*offerdiscount)
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
	database.DB.Model(&models.Product{}).Where("id = ?", productid).Update("offer_discount_percent", 0)
	database.DB.Model(&models.CartItems{}).Where("product_id = ?", productid).Update("discount", 0.00)
	var totalamount float64
	database.DB.Model(&models.CartItems{}).Where("product_id = ?", productid).Pluck("total_amount", &totalamount)
	database.DB.Model(&models.CartItems{}).Where("product_id = ?", productid).Update("final_amount", totalamount)
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "offer deleted succesfully"})
}

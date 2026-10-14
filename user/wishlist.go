package user

import (
	"fmt"
	"net/http"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/helper"
	"github.com/Ansalps/GeZOne/middleware"
	"github.com/Ansalps/GeZOne/models"
	requestmodels "github.com/Ansalps/GeZOne/requestmodels"
	"github.com/Ansalps/GeZOne/responsemodels"
	"github.com/gin-gonic/gin"
)

func Wishlist(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Claims not found"})
		return
	}

	customClaims, ok := claims.(*middleware.CustomClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
		return
	}

	userID := customClaims.ID

	var wishlist []responsemodels.Wishlist
	database.DB.Raw(`
		SELECT
			wishlists.id,
			wishlists.created_at,
			wishlists.updated_at,
			wishlists.user_id,
			wishlists.product_id,
			products.product_name,
			categories.category_name,
			products.description AS product_description,
			products.description,
			products.image_url,
			products.image_url AS product_image_url,
			products.price,
			COALESCE(SUM(product_variants.stock), 0) AS stock,
			products.popular,
			COALESCE(string_agg(DISTINCT product_variants.size, ',' ORDER BY product_variants.size), '') AS size,
			COALESCE(MAX(CASE WHEN offers.start_at <= now() AND (offers.end_at IS NULL OR offers.end_at >= now()) THEN offers.discount_percentage ELSE 0 END), 0) AS discount_percentage,
			COALESCE(MAX(to_char(offers.start_at, 'YYYY-MM-DD')), '') AS start_date,
			COALESCE(MAX(to_char(offers.end_at, 'YYYY-MM-DD')), '') AS end_date
		FROM wishlists
		JOIN products ON wishlists.product_id = products.id
		JOIN categories ON categories.id = products.category_id
		LEFT JOIN product_variants ON product_variants.product_id = products.id
		LEFT JOIN offers ON offers.product_id = products.id
		WHERE wishlists.user_id = ?
		GROUP BY
			wishlists.id,
			wishlists.created_at,
			wishlists.updated_at,
			wishlists.user_id,
			wishlists.product_id,
			products.product_name,
			categories.category_name,
			products.description,
			products.image_url,
			products.price,
			products.popular
	`, userID).Scan(&wishlist)
	c.JSON(http.StatusOK, gin.H{
		"data":    wishlist,
		"message": "succesfully shown wishlist",
	})
}

func WishlistAdd(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Claims not found"})
		return
	}

	customClaims, ok := claims.(*middleware.CustomClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
		return
	}

	userID := customClaims.ID

	var wishlistadd requestmodels.WishlistAdd
	err := c.BindJSON(&wishlistadd)
	response := gin.H{
		"status":  false,
		"message": "failed to bind request",
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// Validate the content of the JSON
	if err := helper.Validate(wishlistadd); err != nil {
		fmt.Println("", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"message":    err.Error(),
			"error_code": http.StatusBadRequest,
		})
		return
	}
	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM wishlists WHERE product_id = ?`, wishlistadd.ProductID).Scan(&count)
	if count != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "product already added in wishlist",
		})
		return
	}
	wishlist := models.Wishlist{
		UserID:    userID,
		ProductID: wishlistadd.ProductID,
	}
	database.DB.Create(&wishlist)
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "Wishlist Added"})
}

func WishlistRemove(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Claims not found"})
		return
	}

	customClaims, ok := claims.(*middleware.CustomClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
		return
	}

	userID := customClaims.ID

	var wishlistadd requestmodels.WishlistAdd
	err := c.BindJSON(&wishlistadd)
	response := gin.H{
		"status":  false,
		"message": "failed to bind request",
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// Validate the content of the JSON
	if err := helper.Validate(wishlistadd); err != nil {
		fmt.Println("", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"message":    err.Error(),
			"error_code": http.StatusBadRequest,
		})
		return
	}
	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM wishlists WHERE product_id = ? AND user_id = ?`, wishlistadd.ProductID, userID).Scan(&count)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "product does not exist wishlist for this particular user",
		})
		return
	}
	database.DB.Where("product_id = ? AND user_id = ?", wishlistadd.ProductID, userID).Delete(&models.Wishlist{})
	c.JSON(http.StatusOK, gin.H{
		"message": "product removed from wishlist successfully",
	})
}

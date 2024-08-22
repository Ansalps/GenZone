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

func CouponList(c *gin.Context) {
	var coupon []responsemodels.Coupon
	database.DB.Raw(`SELECT * FROM coupons WHERE deleted_at IS NULL`).Scan(&coupon)
	c.JSON(http.StatusOK, gin.H{
		"data":    coupon,
		"message": "listing coupons successfully",
	})
}

func CouponAdd(c *gin.Context) {
	var couponadd models.CouponAdd
	err := c.BindJSON(&couponadd)
	response := gin.H{
		"status":  false,
		"message": "failed to bind request",
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// Validate the content of the JSON
	if err := helper.Validate(couponadd); err != nil {
		fmt.Println("", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"message":    err.Error(),
			"error_code": http.StatusBadRequest,
		})
		return
	}
	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM coupons WHERE code = ? AND deleted_at IS NULL`, couponadd.Code).Scan(&count)
	if count != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "coupon code already exists",
		})
		return
	}
	coupon := models.Coupon{
		Code:        couponadd.Code,
		Discount:    couponadd.Discount,
		MinPurchase: couponadd.MinPurchase,
	}
	database.DB.Create(&coupon)
	c.JSON(http.StatusOK, gin.H{
		"message": "coupon addded succeessfully",
	})
}
func CouponRemove(c *gin.Context) {
	CouponID := c.Param("id")
	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM coupons WHERE id = ? AND deleted_at IS NULL`, CouponID).Scan(&count)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Coupon id id does not exist",
		})
		return
	}
	database.DB.Where("id = ?", CouponID).Delete(&models.Coupon{})
	c.JSON(http.StatusOK, gin.H{
		"message": "Coupon deleted successfully",
	})
}

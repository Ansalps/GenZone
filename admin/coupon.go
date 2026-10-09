package admin

import (
	"net/http"
	"strings"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/helper"
	"github.com/Ansalps/GeZOne/models"
	"github.com/Ansalps/GeZOne/requestmodels"
	"github.com/Ansalps/GeZOne/responsemodels"
	"github.com/gin-gonic/gin"
)

func CouponList(c *gin.Context) {
	var coupons []responsemodels.Coupon
	if err := database.DB.Order("created_at DESC").Find(&coupons).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "failed to fetch coupons",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "listing coupons successfully",
		"data": gin.H{
			"coupons": coupons,
		},
	})
}

func CouponAdd(c *gin.Context) {
	var req requestmodels.CouponAdd
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

	code := strings.TrimSpace(req.Code)
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "coupon code is required",
		})
		return
	}

	if req.Discount < 0 || req.MinPurchase < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "discount and minimum purchase must be greater than or equal to 0",
		})
		return
	}

	startAt, endAt, err := parseProductOfferDates(req.StartAt, req.EndAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	var existing models.Coupon
	if err := database.DB.Where("code = ?", code).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "coupon code already exists",
		})
		return
	}

	coupon := models.Coupon{
		Code:        strings.ToUpper(code),
		Discount:    req.Discount,
		MinPurchase: req.MinPurchase,
		StartAt:     startAt,
		EndAt:       endAt,
		IsActive:    req.IsActive,
	}
	if err := database.DB.Create(&coupon).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "failed to add coupon",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "coupon added successfully",
	})
}

func CouponEdit(c *gin.Context) {
	couponID := c.Param("id")
	var existing models.Coupon
	if err := database.DB.First(&existing, couponID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "coupon id does not exist",
		})
		return
	}

	var req requestmodels.CouponAdd
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

	if req.Discount < 0 || req.MinPurchase < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "discount and minimum purchase must be greater than or equal to 0",
		})
		return
	}

	startAt, endAt, err := parseProductOfferDates(req.StartAt, req.EndAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	code := strings.TrimSpace(req.Code)
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "coupon code is required",
		})
		return
	}

	var duplicate models.Coupon
	if err := database.DB.Where("code = ? AND id != ?", code, existing.ID).First(&duplicate).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "coupon code already exists",
		})
		return
	}

	existing.Code = strings.ToUpper(code)
	existing.Discount = req.Discount
	existing.MinPurchase = req.MinPurchase
	existing.StartAt = startAt
	existing.EndAt = endAt
	existing.IsActive = req.IsActive
	if err := database.DB.Save(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "failed to update coupon",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "coupon updated successfully",
	})
}

func CouponActivate(c *gin.Context) {
	couponID := c.Param("id")
	var coupon models.Coupon
	if err := database.DB.First(&coupon, couponID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "coupon id does not exist",
		})
		return
	}

	coupon.IsActive = true
	if err := database.DB.Save(&coupon).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "failed to activate coupon",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "coupon activated successfully",
	})
}

func CouponInactivate(c *gin.Context) {
	couponID := c.Param("id")
	var coupon models.Coupon
	if err := database.DB.First(&coupon, couponID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "coupon id does not exist",
		})
		return
	}

	coupon.IsActive = false
	if err := database.DB.Save(&coupon).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "failed to deactivate coupon",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "coupon deactivated successfully",
	})
}

func CouponRemove(c *gin.Context) {
	couponID := c.Param("id")
	var coupon models.Coupon
	if err := database.DB.First(&coupon, couponID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "coupon id does not exist",
		})
		return
	}

	if err := database.DB.Delete(&coupon).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "failed to delete coupon",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "coupon deleted successfully",
	})
}

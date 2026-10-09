package user

import (
	"net/http"
	"strings"
	"time"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/middleware"
	"github.com/Ansalps/GeZOne/models"
	"github.com/Ansalps/GeZOne/requestmodels"
	"github.com/gin-gonic/gin"
)

func CouponCheckout(c *gin.Context) {
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

	var req requestmodels.CouponCheckout
	if err := c.ShouldBindJSON(&req); err != nil && c.Request.ContentLength > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "failed to bind request",
		})
		return
	}

	couponCode := strings.TrimSpace(req.CouponCode)
	if couponCode == "" {
		couponCode = strings.TrimSpace(c.Query("coupon_code"))
	}
	if couponCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "coupon code is required",
		})
		return
	}

	var coupon models.Coupon
	if err := database.DB.Where("code = ? AND deleted_at IS NULL", couponCode).First(&coupon).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "such coupon does not exist",
		})
		return
	}

	if !coupon.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "this coupon is currently inactive",
		})
		return
	}

	now := time.Now()
	if !coupon.StartAt.IsZero() && coupon.StartAt.After(now) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "this coupon is not active yet",
		})
		return
	}
	if !coupon.EndAt.IsZero() && coupon.EndAt.Before(now) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "this coupon has expired",
		})
		return
	}

	var cartTotal float64
	query := `
		WITH item_summary AS (
			SELECT
				ci.id,
				ci.product_id,
				ci.quantity,
				CASE
					WHEN COALESCE(MAX(CASE WHEN o.start_at <= now() AND (o.end_at IS NULL OR o.end_at >= now()) THEN o.discount_percentage ELSE 0 END), 0) > 0
					THEN p.price * (1 - (COALESCE(MAX(CASE WHEN o.start_at <= now() AND (o.end_at IS NULL OR o.end_at >= now()) THEN o.discount_percentage ELSE 0 END), 0) / 100.0))
					ELSE p.price
				END AS offer_unit_price,
				p.price * ci.quantity AS original_total,
				CASE
					WHEN COALESCE(MAX(CASE WHEN o.start_at <= now() AND (o.end_at IS NULL OR o.end_at >= now()) THEN o.discount_percentage ELSE 0 END), 0) > 0
					THEN (p.price * ci.quantity) * (1 - (COALESCE(MAX(CASE WHEN o.start_at <= now() AND (o.end_at IS NULL OR o.end_at >= now()) THEN o.discount_percentage ELSE 0 END), 0) / 100.0))
					ELSE p.price * ci.quantity
				END AS final_total
			FROM carts c
			JOIN cart_items ci ON c.id = ci.cart_id
			JOIN products p ON ci.product_id = p.id
			LEFT JOIN offers o ON o.product_id = p.id
			WHERE c.user_id = ? AND ci.quantity > 0
			GROUP BY ci.id, ci.product_id, ci.quantity, p.price
		)
		SELECT COALESCE(SUM(final_total), 0)
		FROM item_summary
	`
	if err := database.DB.Raw(query, userID).Scan(&cartTotal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "failed to calculate cart total",
		})
		return
	}

	if cartTotal < coupon.MinPurchase {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "coupon can be applied if you purchase with a minimum amount",
		})
		return
	}

	couponDiscount := coupon.Discount
	if couponDiscount > cartTotal {
		couponDiscount = cartTotal
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "coupon applied successfully",
		"data": gin.H{
			"coupon_code":           couponCode,
			"coupon_discount":       couponDiscount,
			"coupon_applied_amount": cartTotal - couponDiscount,
			"minimum_purchase":      coupon.MinPurchase,
		},
	})
}

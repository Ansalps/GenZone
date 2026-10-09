package user

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/helper"
	"github.com/Ansalps/GeZOne/middleware"
	"github.com/Ansalps/GeZOne/models"
	"github.com/Ansalps/GeZOne/requestmodels"
	"github.com/Ansalps/GeZOne/responsemodels"
	"github.com/gin-gonic/gin"
)

func Order(c *gin.Context) {
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

	var OrderAdd requestmodels.OrderAdd
	if err := c.BindJSON(&OrderAdd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "failed to bind request",
		})
		return
	}

	if err := helper.Validate(OrderAdd); err != nil {
		fmt.Println("", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"message":    err.Error(),
			"error_code": http.StatusBadRequest,
		})
		return
	}

	var count1 int64
	database.DB.Raw(`SELECT COUNT(*) FROM addresses WHERE id = ? AND user_id = ? AND deleted_at IS NULL`, OrderAdd.AddressID, userID).Scan(&count1)
	if count1 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "address_id does not exist for this particular user",
		})
		return
	}

	var cart models.Cart
	if err := database.DB.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "cart empty, order can't be placed",
		})
		return
	}

	var cartItems []models.CartItem
	if err := database.DB.Where("cart_id = ?", cart.ID).Find(&cartItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "failed to load cart items",
		})
		return
	}
	if len(cartItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "cart empty, order can't be placed",
		})
		return
	}

	totalAmount := 0.0
	totalQuantity := uint(0)
	for _, item := range cartItems {
		if item.Quantity == 0 {
			continue
		}
		totalAmount += item.UnitPrice * float64(item.Quantity)
		totalQuantity += item.Quantity
	}

	if totalQuantity == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "cart empty, order can't be placed",
		})
		return
	}

	if OrderAdd.PaymentMethod == "COD" && totalAmount > 100000 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "orders above Rs.1000 cannot be done through cash on delivery",
		})
		return
	}

	discountAmount := 0.0
	var couponID uint
	if OrderAdd.CouponCode != "" {
		var coupon models.Coupon
		if err := database.DB.Where("code = ? AND is_active = true AND start_at <= ? AND end_at >= ?", OrderAdd.CouponCode, time.Now(), time.Now()).First(&coupon).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "such coupon does not exists or is not active",
			})
			return
		}
		if totalAmount < coupon.MinPurchase {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "A minimum purchase amount is required to apply this coupon",
			})
			return
		}
		discountAmount = coupon.Discount
		couponID = coupon.ID
	}

	offerDiscountAmount := 0.0
	for _, item := range cartItems {
		if item.Quantity == 0 {
			continue
		}
		var product models.Product
		if err := database.DB.First(&product, item.ProductID).Error; err != nil {
			continue
		}

		var activeOffer models.Offer
		if err := database.DB.Where("product_id = ? AND start_at <= ? AND end_at >= ?", item.ProductID, time.Now(), time.Now()).Order("created_at DESC").First(&activeOffer).Error; err == nil {
			offerDiscountAmount += (item.UnitPrice * float64(item.Quantity)) * (activeOffer.DiscountPercentage / 100)
		}
	}

	paymentMethod := OrderAdd.PaymentMethod
	if paymentMethod == "" {
		paymentMethod = "COD"
	}

	finalAmount := totalAmount - discountAmount
	if finalAmount < 0 {
		finalAmount = 0
	}

	order := models.Order{
		UserID:              userID,
		AddressID:           OrderAdd.AddressID,
		TotalAmount:         totalAmount,
		PaymentMethod:       paymentMethod,
		OrderStatus:         "pending",
		CouponID:            couponID,
		TotalDiscountAmount: offerDiscountAmount + discountAmount,
		FinalAmount:         finalAmount,
	}
	if err := database.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "failed to create order",
		})
		return
	}

	for _, item := range cartItems {
		if item.Quantity == 0 {
			continue
		}

		unitPrice := item.UnitPrice
		if unitPrice == 0 {
			var product models.Product
			if err := database.DB.First(&product, item.ProductID).Error; err == nil {
				unitPrice = product.Price
			}
		}

		itemOfferDiscount := 0.0
		var activeOffer models.Offer
		if err := database.DB.Where("product_id = ? AND start_at <= ? AND end_at >= ?", item.ProductID, time.Now(), time.Now()).Order("created_at DESC").First(&activeOffer).Error; err == nil {
			itemOfferDiscount = (unitPrice * float64(item.Quantity)) * (activeOffer.DiscountPercentage / 100)
		}

		itemCouponDiscount := 0.0
		if discountAmount > 0 && totalAmount > 0 {
			itemCouponDiscount = (unitPrice * float64(item.Quantity) / totalAmount) * discountAmount
		}

		orderItem := models.OrderItems{
			OrderID:        order.ID,
			ProductID:      item.ProductID,
			Price:          unitPrice,
			OrderStatus:    "pending",
			PaymentMethod:  paymentMethod,
			CouponDiscount: itemCouponDiscount,
			OfferDiscount:  itemOfferDiscount,
			TotalDiscount:  itemOfferDiscount + itemCouponDiscount,
			PaidAmount:     unitPrice * float64(item.Quantity),
		}
		if err := database.DB.Create(&orderItem).Error; err != nil {
			continue
		}
	}

	payment := models.Payments{
		UserID:        userID,
		OrderID:       order.ID,
		TotalAmount:   finalAmount,
		PaymentType:   paymentMethod,
		PaymentStatus: "pending",
	}
	if err := database.DB.Create(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "failed to create payment record",
		})
		return
	}

	if err := database.DB.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "failed to clear cart after placing order",
		})
		return
	}

	var responseOrder responsemodels.Order
	var address responsemodels.Address
	var responseOrderItems []responsemodels.OrderItems

	database.DB.Raw(`SELECT orders.id,orders.created_at,orders.updated_at,orders.deleted_at,orders.user_id,orders.address_id,orders.total_amount,orders.payment_method,orders.order_status,orders.offer_applied,orders.coupon_code,orders.discount_amount,orders.final_amount FROM orders JOIN addresses ON orders.address_id = addresses.id WHERE orders.user_id = ? ORDER BY orders.created_at DESC LIMIT 1`, userID).Scan(&responseOrder)

	var latestOrderID uint
	database.DB.Raw(`SELECT id FROM orders WHERE user_id = ? ORDER BY created_at DESC LIMIT 1`, userID).Scan(&latestOrderID)
	var latestAddressID uint
	database.DB.Raw(`SELECT address_id FROM orders WHERE user_id = ? ORDER BY created_at DESC LIMIT 1`, userID).Scan(&latestAddressID)
	database.DB.Raw(`SELECT * FROM addresses WHERE id = ?`, latestAddressID).Scan(&address)
	responseOrder.Address = address
	database.DB.Raw(`SELECT order_items.id,order_items.created_at,order_items.updated_at,order_items.deleted_at,order_items.order_id,order_items.product_id,products.product_name,order_items.price,order_items.order_status,order_items.payment_method,order_items.coupon_discount,order_items.offer_discount,order_items.total_discount,order_items.paid_amount FROM order_items JOIN products ON order_items.product_id = products.id WHERE order_items.order_id = ? ORDER BY order_items.id`, latestOrderID).Scan(&responseOrderItems)

	c.JSON(http.StatusOK, gin.H{
		"message":     "Order added successfully",
		"order":       responseOrder,
		"order_items": responseOrderItems,
	})
}

func OrderList(c *gin.Context) {
	//userID := c.Param("user_id")
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

	listorder := c.Query("list_order")
	var orders []responsemodels.Order
	var address responsemodels.Address

	sql := `SELECT orders.id,orders.created_at,orders.updated_at,orders.user_id,orders.address_id,orders.total_amount,orders.offer_discount,orders.coupon_discount,orders.total_discount_amount,orders.coupon_id,orders.order_status,orders.final_amount,orders.payment_method,addresses.user_id,addresses.country,addresses.state,addresses.street_name,addresses.city,addresses.pin_code,addresses.phone,addresses.default
	FROM orders
	JOIN addresses ON orders.address_id = addresses.id where orders.user_id = ?`
	switch listorder {
	case "":
		sql += ` ORDER BY orders.id ASC`
	case "ASC":
		sql += ` ORDER BY orders.id ASC`
	case "DSC":
		sql += ` ORDER BY orders.id DESC`
	}

	database.DB.Raw(sql, userID).Scan(&orders)
	
	for i, v := range orders {
		database.DB.Raw(`SELECT *
	        FROM orders
	        JOIN addresses ON orders.address_id = addresses.id
	        WHERE orders.user_id = ? AND orders.id = ?`, userID, v.ID).Scan(&address)
		orders[i].Address = address
	}

	// query.Find(&Address)
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "successfully retrieved user informations",
		"data": gin.H{
			//"Address": Address,
			"Order": orders,
		},
	})
}

func OrderItemsList(c *gin.Context) {
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

	orderId := c.Param("order_id")
	listorder := c.Query("list_order")

	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM orders where id = ? AND user_id = ?`, orderId, userID).Scan(&count)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "order id does not exist for this particular user",
		})
		return
	}
	var orderitems []responsemodels.OrderItems

	sql := `SELECT order_items.id,order_items.created_at,order_items.updated_at,order_items.deleted_at,order_items.order_id,order_items.product_id,products.product_name,order_items.price,order_items.order_status,order_items.payment_method,order_items.coupon_discount,order_items.offer_discount,order_items.total_discount,order_items.paid_amount FROM order_items join products on order_items.product_id=products.id WHERE order_items.order_id = ?`

	switch listorder {
	case "":
		sql += ` ORDER BY order_items.id ASC`
	case "ASC":
		sql += ` ORDER BY order_items.id ASC`
	case "DSC":
		sql += ` ORDER BY order_items.id DESC`
	}

	database.DB.Raw(sql, orderId).Scan(&orderitems)
	c.JSON(http.StatusOK, gin.H{
		"order items": orderitems,
	})
}
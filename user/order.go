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

func shouldCancelEntireOrder(remainingStatuses []string) bool {
	if len(remainingStatuses) == 0 {
		return true
	}
	for _, status := range remainingStatuses {
		if status != "cancelled" && status != "return" && status != "" {
			return false
		}
	}
	return true
}

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

	couponDiscountAmount := 0.0
	var couponID *uint
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
		couponDiscountAmount = coupon.Discount
		couponID = &coupon.ID
	}

	if couponID != nil {
		if err := database.DB.Model(&models.Coupon{}).Where("id = ?", *couponID).Pluck("discount", &couponDiscountAmount).Error; err != nil {
			couponDiscountAmount = 0
		}
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

	finalAmount := totalAmount - offerDiscountAmount - couponDiscountAmount
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
		CouponCode:          OrderAdd.CouponCode,
		OfferDiscount:       offerDiscountAmount,
		CouponDiscount:      couponDiscountAmount,
		TotalDiscountAmount: offerDiscountAmount + couponDiscountAmount,
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
		if couponDiscountAmount > 0 && totalAmount > 0 {
			itemCouponDiscount = (unitPrice * float64(item.Quantity) / totalAmount) * couponDiscountAmount
		}

		fmt.Println("quantit", item.Quantity)

		lineTotalAmount := unitPrice * float64(item.Quantity)
		paidAmount := lineTotalAmount - itemOfferDiscount - itemCouponDiscount
		orderItem := models.OrderItems{
			OrderID:        order.ID,
			ProductID:      item.ProductID,
			Quantity:       item.Quantity,
			UnitPrice:      unitPrice,
			TotalAmount:    lineTotalAmount,
			OrderStatus:    "pending",
			PaymentMethod:  paymentMethod,
			CouponDiscount: itemCouponDiscount,
			OfferDiscount:  itemOfferDiscount,
			TotalDiscount:  itemOfferDiscount + itemCouponDiscount,
			PaidAmount:     paidAmount,
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

	database.DB.Raw(`SELECT orders.id,orders.created_at,orders.updated_at,orders.user_id,orders.address_id,orders.total_amount,orders.payment_method,orders.order_status,orders.offer_discount,orders.coupon_discount,orders.total_discount_amount,orders.final_amount FROM orders JOIN addresses ON orders.address_id = addresses.id WHERE orders.user_id = ? ORDER BY orders.created_at DESC LIMIT 1`, userID).Scan(&responseOrder)

	var latestOrderID uint
	database.DB.Raw(`SELECT id FROM orders WHERE user_id = ? ORDER BY created_at DESC LIMIT 1`, userID).Scan(&latestOrderID)
	var latestAddressID uint
	database.DB.Raw(`SELECT address_id FROM orders WHERE user_id = ? ORDER BY created_at DESC LIMIT 1`, userID).Scan(&latestAddressID)
	database.DB.Raw(`SELECT * FROM addresses WHERE id = ?`, latestAddressID).Scan(&address)
	responseOrder.Address = address
	database.DB.Raw(`SELECT order_items.id,order_items.created_at,order_items.updated_at,order_items.deleted_at,order_items.order_id,order_items.product_id,products.product_name,order_items.quantity,order_items.unit_price AS price,order_items.total_amount,order_items.order_status,order_items.payment_method,order_items.coupon_discount,order_items.offer_discount,order_items.total_discount,order_items.paid_amount FROM order_items JOIN products ON order_items.product_id = products.id WHERE order_items.order_id = ? ORDER BY order_items.id`, latestOrderID).Scan(&responseOrderItems)

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

	sql := `SELECT order_items.id,order_items.created_at,order_items.updated_at,order_items.deleted_at,order_items.order_id,order_items.product_id,products.product_name,order_items.quantity,order_items.unit_price AS price,order_items.total_amount,order_items.order_status,order_items.payment_method,order_items.coupon_discount,order_items.offer_discount,order_items.total_discount,order_items.paid_amount FROM order_items join products on order_items.product_id=products.id WHERE order_items.order_id = ?`

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

func CancelOrder(c *gin.Context) {
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

	orderID := c.Param("order_id")
	var orderid uint
	database.DB.Model(&models.Order{}).Where("id = ?", orderID).Pluck("id", &orderid)
	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM orders where id = ? AND user_id = ?`, orderID, userID).Scan(&count)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "order id does not exist for this particular user",
		})
		return
	}

	var orderstatus string
	database.DB.Model(&models.Order{}).Where("id = ?", orderID).Pluck("order_status", &orderstatus)
	if orderstatus == "cancelled" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cannot cancel the order item, item is already cancelled bu user or admin",
		})
		return
	}
	if orderstatus == "shipped" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cannot cancel the order, items are already shipped",
		})
		return
	}
	if orderstatus == "delivered" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cannot cancel the order, items are delivered",
		})
		return
	}
	if orderstatus == "failed" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "your order is already failed due to some inventory issues",
		})
		return
	}
	database.DB.Model(&models.Order{}).Where("id = ?", orderID).Update("order_status", "cancelled")
	database.DB.Model(&models.OrderItems{}).Where("order_id = ?", orderID).Update("order_status", "cancelled")
	var paymentmethod string
	database.DB.Model(&models.Order{}).Where("id = ?", orderID).Pluck("payment_method", &paymentmethod)
	var totalamount float64
	database.DB.Model(&models.Order{}).Where("id = ?", orderID).Pluck("final_amount", &totalamount)
	if paymentmethod == "RazorPay" || paymentmethod == "Wallet" {

		var count int64
		database.DB.Raw(`SELECT COUNT(*) FROM wallets WHERE user_id = ?`, userID).Scan(&count)
		if count == 0 {
			wallet := models.Wallet{
				UserID:  userID,
				Balance: totalamount,
			}
			database.DB.Create(&wallet)
		} else {
			var balance float64
			database.DB.Model(&models.Wallet{}).Where("user_id = ?", userID).Pluck("balance", &balance)
			balance = balance + totalamount
			wallet := models.Wallet{
				UserID:  userID,
				Balance: balance,
			}
			database.DB.Model(&models.Wallet{}).Where("user_id = ?", userID).Updates(&wallet)
		}
		transaction := models.WalletTransaction{
			UserID:          userID,
			Amount:          totalamount,
			TransactionType: "Credit",
			Description:     "Refund for cancelling order",
		}
		database.DB.Create(&transaction)
		now := time.Now()
		today := now.Format("2006-01-02")
		payment := models.Payments{
			UserID:        userID,
			OrderID:       orderid,
			TotalAmount:   totalamount,
			PaymentDate:   today,
			PaymentType:   paymentmethod,
			PaymentStatus: "refund",
			Description:   "refund for cancelling order",
		}
		database.DB.Create(&payment)
	}
	//database.DB.Raw(`UPDATE order_items SET order_status = 'cancelled' WHERE `)
	c.JSON(http.StatusOK, gin.H{"message": "Order cancelled successfully"})
}

func CancelSingleOrderItem(c *gin.Context) {
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

	//plucking order item id
	ItemID := c.Param("orderitem_id")
	var itemid uint
	database.DB.Model(&models.OrderItems{}).Where("id = ?", itemid).Pluck("id", &itemid)
	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM order_items where id = ?`, ItemID).Scan(&count)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "order item id does not exist for this particular user",
		})
		return
	}

	//plucking orderstatus
	var orderstatus string
	database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Pluck("order_status", &orderstatus)

	if orderstatus == "cancelled" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cannot cancel the order item, item is already cancelled",
		})
		return
	}
	if orderstatus == "return" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cannot cancel the order item, item is already returned",
		})
		return
	}
	if orderstatus == "shipped" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cannot cancel the order item, item is already shipped",
		})
		return
	}
	if orderstatus == "delivered" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cannot cancel the order item, item is delivered",
		})
		return
	}
	if orderstatus == "failed" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "your order is already failed due to some inventory issues",
		})
		return
	}

	var unitPrice float64
	database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Pluck("unit_price", &unitPrice)
	var orderid uint
	database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Pluck("order_id", &orderid)
	var offerdiscount float64
	database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Pluck("offer_discount", &offerdiscount)
	var itemPaidAmount float64
	database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Pluck("paid_amount", &itemPaidAmount)
	if itemPaidAmount == 0 {
		itemPaidAmount = (unitPrice - offerdiscount)
	}
	var remainingStatuses []string
	database.DB.Raw(`SELECT order_status FROM order_items WHERE order_id = ? AND id != ?`, orderid, ItemID).Scan(&remainingStatuses)
	finalizeOrderCancellation := shouldCancelEntireOrder(remainingStatuses)
	var paymentmethod string
	database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Pluck("payment_method", &paymentmethod)
	var finalamount float64
	database.DB.Model(&models.Order{}).Where("id = ?", orderid).Pluck("final_amount", &finalamount)
	if finalamount > itemPaidAmount {
		database.DB.Model(&models.Order{}).Where("id = ?", orderid).Update("final_amount", finalamount-itemPaidAmount)
	}

	database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Update("order_status", "cancelled")
	if finalizeOrderCancellation {
		database.DB.Model(&models.Order{}).Where("id = ?", orderid).Updates(map[string]interface{}{
			"order_status":          "cancelled",
			"final_amount":          0.0,
			"coupon_code":           "",
			"coupon_discount":       0.0,
			"offer_discount":        0.0,
			"total_discount_amount": 0.0,
		})
	}

	if paymentmethod == "RazorPay" || paymentmethod == "Wallet" {
		var count int64
		database.DB.Raw(`SELECT COUNT(*) FROM wallets WHERE user_id = ?`, userID).Scan(&count)
		if count == 0 {
			wallet := models.Wallet{UserID: userID, Balance: itemPaidAmount}
			database.DB.Create(&wallet)
		} else {
			var balance float64
			database.DB.Model(&models.Wallet{}).Where("user_id = ?", userID).Pluck("balance", &balance)
			wallet := models.Wallet{UserID: userID, Balance: balance + itemPaidAmount}
			database.DB.Model(&models.Wallet{}).Where("user_id = ?", userID).Updates(&wallet)
		}
		transaction := models.WalletTransaction{
			UserID:          userID,
			Amount:          itemPaidAmount,
			TransactionType: "Credit",
			Description:     "Refund for cancelling single order item",
		}
		database.DB.Create(&transaction)
		now := time.Now()
		today := now.Format("2006-01-02")
		payment := models.Payments{
			UserID:        userID,
			OrderID:       orderid,
			OrderItemID:   itemid,
			TotalAmount:   itemPaidAmount,
			PaymentDate:   today,
			PaymentType:   paymentmethod,
			PaymentStatus: "refund",
			Description:   "refund for cancelling single order item",
		}
		database.DB.Create(&payment)
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order item cancelled successfully"})
}

func ReturnSingleOrderItem(c *gin.Context) {
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

	ItemID := c.Param("orderitem_id")
	var itemid uint
	database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Pluck("id", &itemid)
	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM order_items where id = ?`, ItemID).Scan(&count)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "order item id does not exist for this particular user",
		})
		return
	}
	var orderstatus string
	database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Pluck("order_status", &orderstatus)

	if orderstatus == "delivered" {

		// var paymentmethod string
		// database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Pluck("payment_method", paymentmethod)

		var paidamount float64
		database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Pluck("paid_amount", &paidamount)
		var orderid uint
		database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Pluck("order_id", &orderid)
		var orderTotalAmount float64
		database.DB.Model(&models.Order{}).Where("id = ?", orderid).Pluck("total_amount", &orderTotalAmount)
		var unitPrice float64
		database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Pluck("unit_price", &unitPrice)
		var offerdiscount float64
		database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Pluck("offer_discount", &offerdiscount)
		offerdiscounted := unitPrice - offerdiscount

		orderTotalAmount = orderTotalAmount - offerdiscounted
		database.DB.Model(&models.Order{}).Where("id = ?", orderid).Update("total_amount", orderTotalAmount)
		var finalamount float64
		database.DB.Model(&models.Order{}).Where("id = ?", orderid).Pluck("final_amount", &finalamount)
		finalamount = finalamount - paidamount
		database.DB.Model(&models.Order{}).Where("id = ?", orderid).Update("final_amount", finalamount)
		var coupondiscount float64
		database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Pluck("coupon_discount", &coupondiscount)
		var totalDiscountAmount float64
		database.DB.Model(&models.Order{}).Where("id = ?", orderid).Pluck("total_discount_amount", &totalDiscountAmount)
		totalDiscountAmount = totalDiscountAmount - coupondiscount
		database.DB.Model(&models.Order{}).Where("id = ?", orderid).Update("total_discount_amount", totalDiscountAmount)
		var balance float64
		database.DB.Model(&models.Wallet{}).Where("user_id = ?", userID).Pluck("balance", &balance)
		balance = balance + paidamount
		database.DB.Model(&models.Wallet{}).Where("user_id = ?", userID).Update("balance", balance)
		//database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Update("paid_amount", 0.00)
		var count1 int
		database.DB.Raw(`SELECT COUNT(*) FROM order_items WHERE order_id = ? AND order_status != 'cancelled' AND order_status != 'return'`, orderid).Scan(&count1)

		if count1 == 0 {
			database.DB.Model(&models.Order{}).Where("id = ?", orderid).Updates(map[string]interface{}{
				"order_status":          "cancelled",
				"final_amount":          0.00,
				"coupon_discount":       0.00,
				"offer_discount":        0.00,
				"total_discount_amount": 0.00,
				"coupon_code":           "",
			})
		}
		database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Update("order_status", "return")
		transaction := models.WalletTransaction{
			UserID:          userID,
			Amount:          paidamount,
			TransactionType: "credit",
			Description:     "Refund for return of single order item",
		}
		database.DB.Create(&transaction)
		now := time.Now()
		today := now.Format("2006-01-02")
		var paymentmethod string
		database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Pluck("payment_method", &paymentmethod)
		payment := models.Payments{
			UserID:        userID,
			OrderID:       orderid,
			OrderItemID:   itemid,
			TotalAmount:   paidamount,
			PaymentDate:   today,
			PaymentType:   paymentmethod,
			PaymentStatus: "refund",
			Description:   "refund for returning single order item",
		}
		database.DB.Create(&payment)

		c.JSON(http.StatusOK, gin.H{
			"message": "order status changed to return successfully",
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cannot return an item until it is delivered",
		})
	}
}

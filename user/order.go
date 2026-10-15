package user

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/helper"
	"github.com/Ansalps/GeZOne/middleware"
	"github.com/Ansalps/GeZOne/models"
	"github.com/Ansalps/GeZOne/requestmodels"
	"github.com/Ansalps/GeZOne/responsemodels"
	"github.com/gin-gonic/gin"
	"github.com/razorpay/razorpay-go"
	"gorm.io/gorm"
)

func shouldCancelEntireOrder(remainingStatuses []string) bool {
	if len(remainingStatuses) == 0 {
		return true
	}
	for _, status := range remainingStatuses {
		if status != "cancelled" && status != "returned" && status != "" {
			return false
		}
	}
	return true
}

func validateCancelQuantity(itemQuantity, requestedQuantity uint) (uint, error) {
	if itemQuantity == 0 {
		return 0, fmt.Errorf("item quantity is invalid")
	}
	if requestedQuantity == 0 {
		return itemQuantity, nil
	}
	if requestedQuantity > itemQuantity {
		return 0, fmt.Errorf("requested quantity %d exceeds available quantity %d", requestedQuantity, itemQuantity)
	}
	return requestedQuantity, nil
}

func validateReturnQuantity(itemQuantity, requestedQuantity uint) (uint, error) {
	if itemQuantity == 0 {
		return 0, fmt.Errorf("item quantity is invalid")
	}
	if requestedQuantity == 0 {
		return itemQuantity, nil
	}
	if requestedQuantity > itemQuantity {
		return 0, fmt.Errorf("requested quantity %d exceeds available quantity %d", requestedQuantity, itemQuantity)
	}
	return requestedQuantity, nil
}

func canCancelByDeliveryStatus(deliveryStatus string) bool {
	return strings.EqualFold(strings.TrimSpace(deliveryStatus), "pending")
}

func canReturnByDeliveryStatus(deliveryStatus string) bool {
	return strings.EqualFold(strings.TrimSpace(deliveryStatus), "delivered")
}

func roundCurrency(value float64) float64 {
	return math.Round(value*100) / 100
}

func proportionalAmount(value float64, quantity uint, totalQuantity uint) float64 {
	if quantity == 0 || totalQuantity == 0 {
		return 0
	}
	return roundCurrency(value * float64(quantity) / float64(totalQuantity))
}

func restoreStockForItem(tx *gorm.DB, item models.OrderItems, quantity uint) error {
	if quantity == 0 {
		return nil
	}
	if item.ProductVariantID != 0 {
		return restoreVariantStockByVariantIDInTx(tx, item.ProductVariantID, quantity)
	}
	return restoreVariantStockInTx(tx, item.ProductID, quantity)
}

func restoreVariantStockInTx(tx *gorm.DB, productID uint, quantity uint) error {
	if quantity == 0 {
		return nil
	}
	var variants []models.ProductVariant
	if err := tx.Where("product_id = ?", productID).Order("id ASC").Find(&variants).Error; err != nil {
		return err
	}
	if len(variants) == 0 {
		return nil
	}
	remaining := int(quantity)
	for i := range variants {
		if remaining <= 0 {
			break
		}
		newStock := variants[i].Stock + uint(remaining)
		if err := tx.Model(&models.ProductVariant{}).Where("id = ?", variants[i].ID).Update("stock", newStock).Error; err != nil {
			return err
		}
		remaining = 0
	}
	if remaining > 0 {
		return fmt.Errorf("unable to restore stock for product %d", productID)
	}
	return nil
}

func restoreVariantStockByVariantIDInTx(tx *gorm.DB, productVariantID uint, quantity uint) error {
	if quantity == 0 || productVariantID == 0 {
		return nil
	}
	var variant models.ProductVariant
	if err := tx.First(&variant, productVariantID).Error; err != nil {
		return err
	}
	newStock := variant.Stock + quantity
	return tx.Model(&models.ProductVariant{}).Where("id = ?", productVariantID).Update("stock", newStock).Error
}

type orderRequest struct {
	requestmodels.OrderAdd
	RazorpayOrderID string `json:"razorpay_order_id"`
	PaymentID       string `json:"payment_id"`
	Signature       string `json:"signature"`
	OrderID         string `json:"order_id"`
}

func verifyRazorpaySignature(orderID string, paymentID string, razorpaySignature string, secret string) bool {
	data := orderID + "|" + paymentID
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	generatedSignature := hex.EncodeToString(h.Sum(nil))
	return hmac.Equal([]byte(generatedSignature), []byte(razorpaySignature))
}

func normalizeRetryPaymentMethod(paymentMethod string) string {
	switch strings.ToLower(strings.TrimSpace(paymentMethod)) {
	case "razorpay", "internet-banking", "internet banking", "netbanking":
		return "RazorPay"
	case "wallet":
		return "Wallet"
	case "cod", "cash on delivery", "cashondelivery":
		return "COD"
	default:
		return "RazorPay"
	}
}

const paymentRetryWindow = time.Minute

func isRetryWindowOpen(createdAt time.Time, now time.Time) bool {
	if createdAt.IsZero() {
		return false
	}
	return now.Sub(createdAt) <= paymentRetryWindow
}

func shouldExpirePayment(createdAt time.Time, now time.Time) bool {
	if createdAt.IsZero() {
		return false
	}
	return now.Sub(createdAt) > paymentRetryWindow
}

func reduceVariantStock(productID uint, quantity uint) error {
	if quantity == 0 {
		return nil
	}
	var variants []models.ProductVariant
	if err := database.DB.Where("product_id = ?", productID).Order("id ASC").Find(&variants).Error; err != nil {
		return err
	}
	if len(variants) == 0 {
		return nil
	}
	remaining := int(quantity)
	for i := range variants {
		if remaining <= 0 {
			break
		}
		if variants[i].Stock == 0 {
			continue
		}
		reduceBy := int(variants[i].Stock)
		if reduceBy > remaining {
			reduceBy = remaining
		}
		newStock := variants[i].Stock - uint(reduceBy)
		if err := database.DB.Model(&models.ProductVariant{}).Where("id = ?", variants[i].ID).Update("stock", newStock).Error; err != nil {
			return err
		}
		remaining -= reduceBy
	}
	if remaining > 0 {
		return fmt.Errorf("insufficient stock available for product %d", productID)
	}
	return nil
}

func reduceVariantStockByVariantID(productVariantID uint, quantity uint) error {
	if quantity == 0 || productVariantID == 0 {
		return nil
	}
	var variant models.ProductVariant
	if err := database.DB.First(&variant, productVariantID).Error; err != nil {
		return err
	}
	if variant.Stock < quantity {
		return fmt.Errorf("insufficient stock available for product variant %d", productVariantID)
	}
	newStock := variant.Stock - quantity
	return database.DB.Model(&models.ProductVariant{}).Where("id = ?", productVariantID).Update("stock", newStock).Error
}

func restoreVariantStock(productID uint, quantity uint) error {
	if quantity == 0 {
		return nil
	}
	var variants []models.ProductVariant
	if err := database.DB.Where("product_id = ?", productID).Order("id ASC").Find(&variants).Error; err != nil {
		return err
	}
	if len(variants) == 0 {
		return nil
	}

	remaining := int(quantity)
	for i := range variants {
		if remaining <= 0 {
			break
		}
		addToVariant := remaining
		newStock := variants[i].Stock + uint(addToVariant)
		if err := database.DB.Model(&models.ProductVariant{}).Where("id = ?", variants[i].ID).Update("stock", newStock).Error; err != nil {
			return err
		}
		remaining = 0
	}

	if remaining > 0 {
		return fmt.Errorf("unable to restore stock for product %d", productID)
	}
	return nil
}

func restoreVariantStockByVariantID(productVariantID uint, quantity uint) error {
	if quantity == 0 || productVariantID == 0 {
		return nil
	}
	var variant models.ProductVariant
	if err := database.DB.First(&variant, productVariantID).Error; err != nil {
		return err
	}
	newStock := variant.Stock + quantity
	return database.DB.Model(&models.ProductVariant{}).Where("id = ?", productVariantID).Update("stock", newStock).Error
}

func reduceStockForOrder(orderID uint) error {
	var orderItems []models.OrderItems
	if err := database.DB.Where("order_id = ?", orderID).Find(&orderItems).Error; err != nil {
		return err
	}
	for _, item := range orderItems {
		if item.ProductVariantID != 0 {
			if err := reduceVariantStockByVariantID(item.ProductVariantID, item.Quantity); err != nil {
				return err
			}
			continue
		}
		if err := reduceVariantStock(item.ProductID, item.Quantity); err != nil {
			return err
		}
	}
	return nil
}

func restoreStockForOrder(orderID uint) error {
	var orderItems []models.OrderItems
	if err := database.DB.Where("order_id = ?", orderID).Find(&orderItems).Error; err != nil {
		return err
	}
	for _, item := range orderItems {
		if item.ProductVariantID != 0 {
			if err := restoreVariantStockByVariantID(item.ProductVariantID, item.Quantity); err != nil {
				return err
			}
			continue
		}
		if err := restoreVariantStock(item.ProductID, item.Quantity); err != nil {
			return err
		}
	}
	return nil
}

func restoreStockForOrderTx(tx *gorm.DB, orderID uint) error {
	var orderItems []models.OrderItems
	if err := tx.Where("order_id = ?", orderID).Find(&orderItems).Error; err != nil {
		return err
	}
	for _, item := range orderItems {
		if item.ProductVariantID != 0 {
			if err := restoreVariantStockByVariantIDInTx(tx, item.ProductVariantID, item.Quantity); err != nil {
				return err
			}
			continue
		}
		if err := restoreVariantStockInTx(tx, item.ProductID, item.Quantity); err != nil {
			return err
		}
	}
	return nil
}

func expirePendingPaymentOrder(order models.Order) error {
	if order.OrderStatus != "pending" || order.PaymentMethod != "RazorPay" {
		return nil
	}
	if !shouldExpirePayment(order.CreatedAt, time.Now()) {
		return nil
	}
	if err := database.DB.Model(&models.Order{}).Where("id = ?", order.ID).Updates(map[string]interface{}{
		"order_status":    "cancelled",
		"delivery_status": "cancelled",
	}).Error; err != nil {
		return err
	}
	if err := database.DB.Model(&models.OrderItems{}).Where("order_id = ?", order.ID).Updates(map[string]interface{}{
		"order_status":    "cancelled",
		"delivery_status": "cancelled",
	}).Error; err != nil {
		return err
	}
	if err := database.DB.Model(&models.Payments{}).Where("order_id = ?", order.ID).Order("created_at DESC").Limit(1).Updates(map[string]interface{}{
		"payment_status": "expired",
		"description":    "Payment expired after 1 minute retry window",
	}).Error; err != nil {
		return err
	}
	return restoreStockForOrder(order.ID)
}

func expirePendingOrdersForUser(userID uint) error {
	var orders []models.Order
	if err := database.DB.Where("user_id = ? AND order_status = ? AND payment_method = ?", userID, "pending", "RazorPay").Find(&orders).Error; err != nil {
		return err
	}
	for _, order := range orders {
		if err := expirePendingPaymentOrder(order); err != nil {
			return err
		}
	}
	return nil
}

func buildRetryOrderUpdates(paymentMethod string, remoteOrderID string) map[string]interface{} {
	method := normalizeRetryPaymentMethod(paymentMethod)
	updates := map[string]interface{}{
		"payment_method":  method,
		"order_status":    "pending",
		"delivery_status": "pending",
	}
	if method == "RazorPay" {
		updates["razorpay_order_id"] = remoteOrderID
	} else {
		updates["razorpay_order_id"] = ""
	}
	return updates
}

func clearCartForUser(userID uint) error {
	var cart models.Cart
	if err := database.DB.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil
		}
		return err
	}
	return database.DB.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error
}

func createRazorpayRemoteOrder(amount float64) (string, int, error) {
	keyID := os.Getenv("RAZORPAY_KEY_ID")
	keySecret := os.Getenv("RAZORPAY_KEY_SECRET")
	if keyID == "" || keySecret == "" {
		return "", 0, fmt.Errorf("razorpay configuration is missing")
	}

	client := razorpay.NewClient(keyID, keySecret)
	data := map[string]interface{}{
		"amount":   int(amount * 100),
		"currency": "INR",
	}
	order, err := client.Order.Create(data, map[string]string{})
	if err != nil {
		return "", 0, err
	}

	remoteOrderID, _ := order["id"].(string)
	checkoutAmountInPaise := 0
	if rawAmount, ok := order["amount"].(float64); ok {
		checkoutAmountInPaise = int(rawAmount)
	}
	return remoteOrderID, checkoutAmountInPaise, nil
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

	var payload orderRequest
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "failed to bind request",
		})
		return
	}

	if payload.PaymentMethod == "" {
		payload.PaymentMethod = "COD"
	}
	if err := helper.Validate(payload.OrderAdd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"message":    err.Error(),
			"error_code": http.StatusBadRequest,
		})
		return
	}

	var count1 int64
	database.DB.Raw(`SELECT COUNT(*) FROM addresses WHERE id = ? AND user_id = ?`, payload.AddressID, userID).Scan(&count1)
	if count1 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "address_id does not exist for this particular user",
		})
		return
	}

	if payload.RazorpayOrderID != "" || payload.PaymentID != "" || payload.Signature != "" || payload.OrderID != "" {
		razorpayOrderID := payload.RazorpayOrderID
		if razorpayOrderID == "" {
			razorpayOrderID = payload.OrderID
		}
		if razorpayOrderID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "razorpay order id is required for payment verification"})
			return
		}
		if payload.PaymentID == "" || payload.Signature == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "payment_id and signature are required for payment verification"})
			return
		}

		var order models.Order
		if err := database.DB.Where("user_id = ? AND razorpay_order_id = ?", userID, razorpayOrderID).First(&order).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "pending razorpay order not found for this user"})
			return
		}

		secret := os.Getenv("RAZORPAY_KEY_SECRET")
		if secret == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Razorpay secret is not configured"})
			return
		}
		if !verifyRazorpaySignature(razorpayOrderID, payload.PaymentID, payload.Signature, secret) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signature"})
			return
		}

		updates := map[string]interface{}{
			"order_status":    "completed",
			"delivery_status": "pending",
		}
		if err := database.DB.Model(&order).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update order status"})
			return
		}
		if err := database.DB.Model(&models.Payments{}).Where("order_id = ?", order.ID).Update("payment_status", "paid").Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update payment status"})
			return
		}
		if err := database.DB.Model(&models.Payments{}).Where("order_id = ?", order.ID).Update("payment_date", time.Now().Format("2006-01-02")).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update payment date"})
			return
		}

		if err := clearCartForUser(userID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "failed to clear cart after payment verification"})
			return
		}

		var responseOrder responsemodels.Order
		var address responsemodels.Address
		var responseOrderItems []responsemodels.OrderItems
		database.DB.Raw(`SELECT orders.id,orders.created_at,orders.updated_at,orders.user_id,orders.address_id,orders.total_amount,orders.payment_method,orders.order_status,orders.offer_discount,orders.coupon_discount,orders.total_discount_amount,orders.final_amount FROM orders JOIN addresses ON orders.address_id = addresses.id WHERE orders.user_id = ? AND orders.id = ? ORDER BY orders.created_at DESC LIMIT 1`, userID, order.ID).Scan(&responseOrder)
		database.DB.Raw(`SELECT * FROM addresses WHERE id = ?`, order.AddressID).Scan(&address)
		responseOrder.Address = address
		database.DB.Raw(`SELECT order_items.id,order_items.created_at,order_items.updated_at,order_items.order_id,order_items.product_id,order_items.product_variant_id,order_items.size,products.product_name,order_items.quantity,order_items.unit_price AS price,order_items.total_amount,order_items.order_status,order_items.payment_method,order_items.coupon_discount,order_items.offer_discount,order_items.total_discount,order_items.final_amount FROM order_items JOIN products ON order_items.product_id = products.id WHERE order_items.order_id = ? ORDER BY order_items.id`, order.ID).Scan(&responseOrderItems)
		c.JSON(http.StatusOK, gin.H{"message": "Payment verified successfully", "order": responseOrder, "order_items": responseOrderItems, "status": "success"})
		return
	}

	var cart models.Cart
	if err := database.DB.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "cart empty, order can't be placed"})
		return
	}

	var cartItems []models.CartItem
	if err := database.DB.Where("cart_id = ?", cart.ID).Find(&cartItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "failed to load cart items"})
		return
	}
	if len(cartItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "cart empty, order can't be placed"})
		return
	}

	totalAmount := 0.0
	totalQuantity := uint(0)
	for _, item := range cartItems {
		if item.Quantity == 0 {
			continue
		}
		if ok, msg := isProductAvailableForUser(item.ProductID, item.Size); !ok {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "cart contains a discontinued product: " + msg})
			return
		}
		totalAmount += item.UnitPrice * float64(item.Quantity)
		totalQuantity += item.Quantity
	}
	if totalQuantity == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "cart empty, order can't be placed"})
		return
	}

	if payload.PaymentMethod == "COD" && totalAmount > 100000 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "orders above Rs.1000 cannot be done through cash on delivery"})
		return
	}

	couponDiscountAmount := 0.0
	var couponID *uint
	if payload.CouponCode != "" {
		var coupon models.Coupon
		if err := database.DB.Where("code = ? AND is_active = true AND start_at <= ? AND end_at >= ?", payload.CouponCode, time.Now(), time.Now()).First(&coupon).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "such coupon does not exists or is not active"})
			return
		}
		if totalAmount < coupon.MinPurchase {
			c.JSON(http.StatusBadRequest, gin.H{"message": "A minimum purchase amount is required to apply this coupon"})
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
		var activeOffer models.Offer
		if err := database.DB.Where("product_id = ? AND start_at <= ? AND end_at >= ?", item.ProductID, time.Now(), time.Now()).Order("created_at DESC").First(&activeOffer).Error; err == nil {
			offerDiscountAmount += (item.UnitPrice * float64(item.Quantity)) * (activeOffer.DiscountPercentage / 100)
		}
	}

	finalAmount := totalAmount - offerDiscountAmount - couponDiscountAmount
	if finalAmount < 0 {
		finalAmount = 0
	}

	orderStatus := "pending"
	deliveryStatus := "pending"
	paymentStatus := "pending"
	paymentType := payload.PaymentMethod
	if payload.PaymentMethod == "Wallet" {
		var wallet models.Wallet
		if err := database.DB.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "wallet is not initialized for this user"})
			return
		}
		if wallet.Balance < finalAmount {
			c.JSON(http.StatusBadRequest, gin.H{"message": "order can't be placed because not enough money in wallet"})
			return
		}
		if err := database.DB.Model(&wallet).Update("balance", wallet.Balance-finalAmount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to deduct wallet balance"})
			return
		}
		paymentStatus = "paid"
		orderStatus = "completed"
	}

	var remoteOrderID string
	var checkoutAmountInPaise int
	if payload.PaymentMethod == "RazorPay" {
		var err error
		remoteOrderID, checkoutAmountInPaise, err = createRazorpayRemoteOrder(finalAmount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Razorpay order"})
			return
		}
	}

	order := models.Order{
		UserID:              userID,
		AddressID:           payload.AddressID,
		TotalAmount:         totalAmount,
		PaymentMethod:       paymentType,
		RazorpayOrderID:     remoteOrderID,
		DeliveryStatus:      deliveryStatus,
		OrderStatus:         orderStatus,
		CouponID:            couponID,
		CouponCode:          payload.CouponCode,
		OfferDiscount:       offerDiscountAmount,
		CouponDiscount:      couponDiscountAmount,
		TotalDiscountAmount: offerDiscountAmount + couponDiscountAmount,
		FinalAmount:         finalAmount,
	}
	if err := database.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "failed to create order"})
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
		var activeOffer *models.Offer
		if offer, ok := helper.GetActiveOfferForProduct(database.DB, item.ProductID); ok {
			activeOffer = offer
			itemOfferDiscount = (unitPrice * float64(item.Quantity)) * (activeOffer.DiscountPercentage / 100)
		}
		itemCouponDiscount := 0.0
		if couponDiscountAmount > 0 && totalAmount > 0 {
			itemCouponDiscount = (unitPrice * float64(item.Quantity) / totalAmount) * couponDiscountAmount
		}
		lineTotalAmount := unitPrice * float64(item.Quantity)
		paidAmount := lineTotalAmount - itemOfferDiscount - itemCouponDiscount
		var offerID *uint
		if activeOffer != nil {
			offerID = &activeOffer.ID
		}
		variantID := uint(0)
		if item.Size != "" {
			var selectedVariant models.ProductVariant
			if err := database.DB.Where("product_id = ? AND size = ?", item.ProductID, item.Size).First(&selectedVariant).Error; err == nil {
				variantID = selectedVariant.ID
			}
		}

		orderItem := models.OrderItems{
			OrderID:          order.ID,
			ProductID:        item.ProductID,
			ProductVariantID: variantID,
			Size:             item.Size,
			OfferID:          offerID,
			Quantity:         item.Quantity,
			UnitPrice:        unitPrice,
			TotalAmount:      lineTotalAmount,
			OrderStatus:      "pending",
			PaymentMethod:    paymentType,
			CouponDiscount:   itemCouponDiscount,
			OfferDiscount:    itemOfferDiscount,
			TotalDiscount:    itemOfferDiscount + itemCouponDiscount,
			PaidAmount:       paidAmount,
		}
		database.DB.Create(&orderItem)
	}
	if err := reduceStockForOrder(order.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "product stock is not available for this order"})
		return
	}

	payment := models.Payments{
		UserID:        userID,
		OrderID:       order.ID,
		TotalAmount:   finalAmount,
		PaymentType:   paymentType,
		PaymentStatus: paymentStatus,
		PaymentDate:   time.Now().Format("2006-01-02"),
	}
	if err := database.DB.Create(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "failed to create payment record"})
		return
	}

	if err := clearCartForUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "failed to clear cart after placing order"})
		return
	}

	if payload.PaymentMethod == "RazorPay" {
		c.JSON(http.StatusOK, gin.H{
			"status":       "success",
			"message":      "Razorpay order created successfully",
			"amount":       checkoutAmountInPaise,
			"currency":     "INR",
			"order_id":     remoteOrderID,
			"receipt":      "receipt_" + remoteOrderID,
			"checkout_url": "/payment-option",
		})
		return
	}

	var responseOrder responsemodels.Order
	var address responsemodels.Address
	var responseOrderItems []responsemodels.OrderItems
	database.DB.Raw(`SELECT orders.id,orders.created_at,orders.updated_at,orders.user_id,orders.address_id,orders.total_amount,orders.payment_method,orders.order_status,orders.offer_discount,orders.coupon_discount,orders.total_discount_amount,orders.final_amount FROM orders JOIN addresses ON orders.address_id = addresses.id WHERE orders.user_id = ? AND orders.id = ? ORDER BY orders.created_at DESC LIMIT 1`, userID, order.ID).Scan(&responseOrder)
	database.DB.Raw(`SELECT * FROM addresses WHERE id = ?`, order.AddressID).Scan(&address)
	responseOrder.Address = address
	database.DB.Raw(`SELECT order_items.id,order_items.created_at,order_items.updated_at,order_items.order_id,order_items.product_id,order_items.product_variant_id,order_items.size,products.product_name,order_items.quantity,order_items.unit_price AS price,order_items.total_amount,order_items.order_status,order_items.payment_method,order_items.coupon_discount,order_items.offer_discount,order_items.total_discount,order_items.final_amount FROM order_items JOIN products ON order_items.product_id = products.id WHERE order_items.order_id = ? ORDER BY order_items.id`, order.ID).Scan(&responseOrderItems)
	c.JSON(http.StatusOK, gin.H{"message": "Order added successfully", "order": responseOrder, "order_items": responseOrderItems})
}

func OrderList(c *gin.Context) {
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
	if err := expirePendingOrdersForUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "failed to refresh order expiry state"})
		return
	}

	listorder := c.Query("list_order")
	var orders []responsemodels.Order
	var address responsemodels.Address

	sql := `SELECT orders.id,orders.created_at,orders.updated_at,orders.user_id,orders.address_id,orders.total_amount,orders.offer_discount,orders.coupon_discount,orders.total_discount_amount,orders.coupon_id,orders.order_status,orders.delivery_status,orders.final_amount,orders.payment_method,COALESCE((SELECT payments.payment_status FROM payments WHERE payments.order_id = orders.id ORDER BY payments.created_at DESC LIMIT 1), 'pending') AS payment_status,addresses.user_id,addresses.country,addresses.state,addresses.street_name,addresses.city,addresses.pin_code,addresses.phone,addresses.default
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

	sql := `SELECT order_items.id,order_items.created_at,order_items.updated_at,order_items.order_id,order_items.product_id,order_items.product_variant_id,order_items.size,products.product_name,order_items.quantity,order_items.unit_price AS price,order_items.total_amount,order_items.order_status,order_items.delivery_status,order_items.payment_method,order_items.coupon_discount,order_items.offer_discount,order_items.total_discount,order_items.final_amount FROM order_items join products on order_items.product_id=products.id WHERE order_items.order_id = ?`

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

func RetryOrderPayment(c *gin.Context) {
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
	if err := expirePendingOrdersForUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "failed to refresh payment expiry state"})
		return
	}
	orderID := c.Param("order_id")

	var payload struct {
		AddressID     uint   `json:"address_id"`
		CouponCode    string `json:"coupon_code"`
		PaymentMethod string `json:"payment_method"`
	}
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "failed to bind retry request"})
		return
	}

	var order models.Order
	if err := database.DB.Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "order not found for this user"})
		return
	}
	if order.OrderStatus == "cancelled" {
		var latestPayment models.Payments
		if err := database.DB.Where("order_id = ?", order.ID).Order("created_at DESC").First(&latestPayment).Error; err == nil && latestPayment.PaymentStatus == "expired" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Payment expired. Retry window has closed."})
			return
		}
	}
	if order.OrderStatus == "completed" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "this order is already completed"})
		return
	}
	if order.PaymentMethod != "RazorPay" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "retry payment is only available for Razorpay orders"})
		return
	}
	if shouldExpirePayment(order.CreatedAt, time.Now()) {
		if err := expirePendingPaymentOrder(order); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "failed to expire payment window"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"message": "Payment expired. Retry window has closed."})
		return
	}

	if payload.AddressID != 0 {
		order.AddressID = payload.AddressID
	}
	if payload.CouponCode != "" {
		order.CouponCode = payload.CouponCode
	}

	method := normalizeRetryPaymentMethod(payload.PaymentMethod)
	updates := buildRetryOrderUpdates(method, "")
	if payload.AddressID != 0 {
		updates["address_id"] = payload.AddressID
	}
	if payload.CouponCode != "" {
		updates["coupon_code"] = payload.CouponCode
	}

	if method == "RazorPay" {
		remoteOrderID, checkoutAmountInPaise, err := createRazorpayRemoteOrder(order.FinalAmount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create Razorpay retry order"})
			return
		}
		updates = buildRetryOrderUpdates(method, remoteOrderID)
		if payload.AddressID != 0 {
			updates["address_id"] = payload.AddressID
		}
		if payload.CouponCode != "" {
			updates["coupon_code"] = payload.CouponCode
		}
		if err := database.DB.Model(&order).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update order for retry"})
			return
		}
		if err := database.DB.Model(&models.Payments{}).Where("order_id = ?", order.ID).Updates(map[string]interface{}{
			"payment_type":   "RazorPay",
			"payment_status": "pending",
			"transaction_id": remoteOrderID,
			"description":    "Razorpay retry initiated",
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update payment retry status"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":            "success",
			"message":           "Retry payment initialized",
			"order_id":          order.ID,
			"razorpay_order_id": remoteOrderID,
			"amount":            checkoutAmountInPaise,
			"currency":          "INR",
		})
		return
	}

	if method == "Wallet" {
		var wallet models.Wallet
		if err := database.DB.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "wallet is not initialized for this user"})
			return
		}
		if wallet.Balance < order.FinalAmount {
			c.JSON(http.StatusBadRequest, gin.H{"message": "wallet balance is insufficient for this order retry"})
			return
		}

		if err := database.DB.Model(&wallet).Update("balance", wallet.Balance-order.FinalAmount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to deduct wallet balance for retry"})
			return
		}
		updates["payment_method"] = "Wallet"
		updates["order_status"] = "completed"
		updates["delivery_status"] = "pending"
	}

	if method == "COD" {
		updates["payment_method"] = "COD"
		updates["order_status"] = "pending"
		updates["delivery_status"] = "pending"
	}

	if err := database.DB.Model(&order).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update order for retry"})
		return
	}
	paymentStatus := "pending"
	if method == "Wallet" {
		paymentStatus = "paid"
	}
	if err := database.DB.Model(&models.Payments{}).Where("order_id = ?", order.ID).Updates(map[string]interface{}{
		"payment_type":   method,
		"payment_status": paymentStatus,
		"description":    "Retry updated to selected payment method",
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update payment retry method"})
		return
	}
	if err := clearCartForUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "failed to clear cart after retry"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Retry payment updated", "order_id": order.ID, "payment_method": method})
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

	var deliveryStatus string
	database.DB.Model(&models.Order{}).Where("id = ?", orderID).Pluck("delivery_status", &deliveryStatus)
	if !canCancelByDeliveryStatus(deliveryStatus) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "order cancellation is only allowed when delivery status is pending",
		})
		return
	}

	database.DB.Model(&models.Order{}).Where("id = ?", orderID).Updates(map[string]interface{}{"order_status": "cancelled", "delivery_status": "cancelled"})
	database.DB.Model(&models.OrderItems{}).Where("order_id = ?", orderID).Updates(map[string]interface{}{"order_status": "cancelled", "delivery_status": "cancelled"})
	if err := restoreStockForOrder(orderid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to restore stock for cancelled order"})
		return
	}
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
	itemID := c.Param("orderitem_id")

	var payload requestmodels.CancelSingleOrderItem
	if err := c.ShouldBindJSON(&payload); err != nil && !strings.Contains(err.Error(), "EOF") {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request payload"})
		return
	}

	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		var item models.OrderItems
		if err := tx.First(&item, itemID).Error; err != nil {
			return fmt.Errorf("order item id does not exist for this particular user")
		}

		var itemBelongsToUser int64
		if err := tx.Raw(`SELECT COUNT(*) FROM order_items oi JOIN orders o ON o.id = oi.order_id WHERE oi.id = ? AND o.user_id = ?`, itemID, userID).Scan(&itemBelongsToUser).Error; err != nil {
			return err
		}
		if itemBelongsToUser == 0 {
			return fmt.Errorf("order item id does not exist for this particular user")
		}

		cancelQuantity, err := validateCancelQuantity(item.Quantity, payload.Quantity)
		if err != nil {
			return err
		}

		var orderDeliveryStatus string
		if err := tx.Model(&models.Order{}).Where("id = ?", item.OrderID).Pluck("delivery_status", &orderDeliveryStatus).Error; err != nil {
			return err
		}
		if !canCancelByDeliveryStatus(item.DeliveryStatus) || !canCancelByDeliveryStatus(orderDeliveryStatus) {
			return fmt.Errorf("order item cancellation is only allowed when delivery status is pending")
		}

		refundAmount := item.PaidAmount
		if refundAmount == 0 {
			refundAmount = item.TotalAmount
		}
		if refundAmount == 0 && item.UnitPrice > 0 {
			refundAmount = item.UnitPrice * float64(item.Quantity)
		}
		if item.Quantity > 0 && cancelQuantity < item.Quantity {
			refundAmount = proportionalAmount(refundAmount, cancelQuantity, item.Quantity)
		}

		remainingQuantity := item.Quantity - cancelQuantity
		remainingTotalAmount := item.TotalAmount
		remainingFinalAmount := item.PaidAmount
		if cancelQuantity == item.Quantity {
			remainingTotalAmount = 0
			remainingFinalAmount = 0
		} else {
			remainingTotalAmount = proportionalAmount(item.TotalAmount, remainingQuantity, item.Quantity)
			if item.PaidAmount > 0 {
				remainingFinalAmount = proportionalAmount(item.PaidAmount, remainingQuantity, item.Quantity)
			} else {
				remainingFinalAmount = remainingTotalAmount
			}
		}

		var remainingStatuses []string
		if err := tx.Raw(`SELECT delivery_status FROM order_items WHERE order_id = ? AND id != ?`, item.OrderID, item.ID).Scan(&remainingStatuses).Error; err != nil {
			return err
		}
		finalizeOrderCancellation := shouldCancelEntireOrder(remainingStatuses)

		var order models.Order
		if err := tx.First(&order, item.OrderID).Error; err == nil {
			if cancelQuantity == item.Quantity {
				if order.FinalAmount > item.PaidAmount {
					if err := tx.Model(&models.Order{}).Where("id = ?", item.OrderID).Update("final_amount", order.FinalAmount-item.PaidAmount).Error; err != nil {
						return err
					}
				}
				if finalizeOrderCancellation {
					if err := tx.Model(&models.Order{}).Where("id = ?", item.OrderID).Updates(map[string]interface{}{
						"order_status":          "cancelled",
						"delivery_status":       "cancelled",
						"final_amount":          0.0,
						"coupon_code":           "",
						"coupon_discount":       0.0,
						"offer_discount":        0.0,
						"total_discount_amount": 0.0,
					}).Error; err != nil {
						return err
					}
					if err := restoreStockForOrderTx(tx, item.OrderID); err != nil {
						return fmt.Errorf("failed to restore stock for cancelled order: %w", err)
					}
				} else if err := restoreStockForItem(tx, item, cancelQuantity); err != nil {
					return fmt.Errorf("failed to restore stock for cancelled item: %w", err)
				}
			} else {
				if order.FinalAmount > refundAmount {
					if err := tx.Model(&models.Order{}).Where("id = ?", item.OrderID).Update("final_amount", order.FinalAmount-refundAmount).Error; err != nil {
						return err
					}
				}
				if err := restoreStockForItem(tx, item, cancelQuantity); err != nil {
					return fmt.Errorf("failed to restore stock for cancelled item: %w", err)
				}
			}
		}

		if cancelQuantity == item.Quantity {
			if err := tx.Model(&models.OrderItems{}).Where("id = ?", item.ID).Updates(map[string]interface{}{
				"quantity":        0,
				"order_status":    "cancelled",
				"delivery_status": "cancelled",
				"total_amount":    0.0,
				"final_amount":    0.0,
				"coupon_discount": 0.0,
				"offer_discount":  0.0,
				"total_discount":  0.0,
			}).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Model(&models.OrderItems{}).Where("id = ?", item.ID).Updates(map[string]interface{}{
				"quantity":        remainingQuantity,
				"total_amount":    remainingTotalAmount,
				"final_amount":    remainingFinalAmount,
				"order_status":    "pending",
				"delivery_status": item.DeliveryStatus,
			}).Error; err != nil {
				return err
			}

			cancelledItem := models.OrderItems{
				OrderID:          item.OrderID,
				ProductID:        item.ProductID,
				ProductVariantID: item.ProductVariantID,
				Size:             item.Size,
				OfferID:          item.OfferID,
				Quantity:         cancelQuantity,
				UnitPrice:        item.UnitPrice,
				TotalAmount:      proportionalAmount(item.TotalAmount, cancelQuantity, item.Quantity),
				DeliveryStatus:   "cancelled",
				OrderStatus:      "cancelled",
				PaymentMethod:    item.PaymentMethod,
				CouponDiscount:   proportionalAmount(item.CouponDiscount, cancelQuantity, item.Quantity),
				OfferDiscount:    proportionalAmount(item.OfferDiscount, cancelQuantity, item.Quantity),
				TotalDiscount:    proportionalAmount(item.TotalDiscount, cancelQuantity, item.Quantity),
				PaidAmount:       proportionalAmount(item.PaidAmount, cancelQuantity, item.Quantity),
			}
			if err := tx.Create(&cancelledItem).Error; err != nil {
				return err
			}
		}

		if item.PaymentMethod == "RazorPay" || item.PaymentMethod == "Wallet" {
			var walletCount int64
			if err := tx.Raw(`SELECT COUNT(*) FROM wallets WHERE user_id = ?`, userID).Scan(&walletCount).Error; err != nil {
				return err
			}
			if walletCount == 0 {
				wallet := models.Wallet{UserID: userID, Balance: refundAmount}
				if err := tx.Create(&wallet).Error; err != nil {
					return err
				}
			} else {
				var balance float64
				if err := tx.Model(&models.Wallet{}).Where("user_id = ?", userID).Pluck("balance", &balance).Error; err != nil {
					return err
				}
				wallet := models.Wallet{UserID: userID, Balance: balance + refundAmount}
				if err := tx.Model(&models.Wallet{}).Where("user_id = ?", userID).Updates(&wallet).Error; err != nil {
					return err
				}
			}
			transaction := models.WalletTransaction{
				UserID:          userID,
				Amount:          refundAmount,
				TransactionType: "Credit",
				Description:     "Refund for cancelling single order item",
			}
			if err := tx.Create(&transaction).Error; err != nil {
				return err
			}
			now := time.Now()
			today := now.Format("2006-01-02")
			payment := models.Payments{
				UserID:        userID,
				OrderID:       item.OrderID,
				OrderItemID:   item.ID,
				TotalAmount:   refundAmount,
				PaymentDate:   today,
				PaymentType:   item.PaymentMethod,
				PaymentStatus: "refund",
				Description:   "refund for cancelling single order item",
			}
			if err := tx.Create(&payment).Error; err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order item cancelled successfully"})
}

func ReturnOrder(c *gin.Context) {
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

	var order models.Order
	if err := database.DB.Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "order id does not exist for this particular user"})
		return
	}
	if !canReturnByDeliveryStatus(order.DeliveryStatus) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "return is only allowed when delivery status is delivered"})
		return
	}

	refundAmount := order.FinalAmount
	if refundAmount <= 0 {
		refundAmount = order.TotalAmount
	}
	if refundAmount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "nothing to refund for this order"})
		return
	}

	if err := restoreStockForOrder(order.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to restore stock for returned order"})
		return
	}

	database.DB.Model(&models.Order{}).Where("id = ?", order.ID).Updates(map[string]interface{}{
		"delivery_status":       "returned",
		"order_status":          "completed",
		"final_amount":          0.0,
		"coupon_discount":       0.0,
		"offer_discount":        0.0,
		"total_discount_amount": 0.0,
		"coupon_code":           "",
	})
	database.DB.Model(&models.OrderItems{}).Where("order_id = ?", order.ID).Updates(map[string]interface{}{
		"delivery_status": "returned",
		"order_status":    "completed",
	})

	if order.PaymentMethod == "RazorPay" || order.PaymentMethod == "Wallet" {
		var walletCount int64
		database.DB.Raw(`SELECT COUNT(*) FROM wallets WHERE user_id = ?`, userID).Scan(&walletCount)
		if walletCount == 0 {
			wallet := models.Wallet{UserID: userID, Balance: refundAmount}
			database.DB.Create(&wallet)
		} else {
			var balance float64
			database.DB.Model(&models.Wallet{}).Where("user_id = ?", userID).Pluck("balance", &balance)
			wallet := models.Wallet{UserID: userID, Balance: balance + refundAmount}
			database.DB.Model(&models.Wallet{}).Where("user_id = ?", userID).Updates(&wallet)
		}
		transaction := models.WalletTransaction{
			UserID:          userID,
			Amount:          refundAmount,
			TransactionType: "Credit",
			Description:     "Refund for returning order",
		}
		database.DB.Create(&transaction)
	}

	now := time.Now()
	today := now.Format("2006-01-02")
	payment := models.Payments{
		UserID:        userID,
		OrderID:       order.ID,
		TotalAmount:   refundAmount,
		PaymentDate:   today,
		PaymentType:   order.PaymentMethod,
		PaymentStatus: "refund",
		Description:   "refund for returning order",
	}
	database.DB.Create(&payment)

	c.JSON(http.StatusOK, gin.H{"message": "order return requested successfully"})
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
	itemID := c.Param("orderitem_id")

	var payload requestmodels.ReturnSingleOrderItem
	if err := c.ShouldBindJSON(&payload); err != nil && !strings.Contains(err.Error(), "EOF") {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request payload"})
		return
	}

	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		var item models.OrderItems
		if err := tx.Where("id = ?", itemID).First(&item).Error; err != nil {
			return fmt.Errorf("order item id does not exist for this particular user")
		}

		var userOrderCount int64
		if err := tx.Raw(`SELECT COUNT(*) FROM order_items oi JOIN orders o ON o.id = oi.order_id WHERE oi.id = ? AND o.user_id = ?`, itemID, userID).Scan(&userOrderCount).Error; err != nil {
			return err
		}
		if userOrderCount == 0 {
			return fmt.Errorf("order item id does not exist for this particular user")
		}

		returnQuantity, err := validateReturnQuantity(item.Quantity, payload.Quantity)
		if err != nil {
			return err
		}
		if !canReturnByDeliveryStatus(item.DeliveryStatus) {
			return fmt.Errorf("cannot return the order item until it is delivered")
		}

		var order models.Order
		if err := tx.Where("id = ?", item.OrderID).First(&order).Error; err == nil && !canReturnByDeliveryStatus(order.DeliveryStatus) {
			return fmt.Errorf("return is only allowed when delivery status is delivered")
		}

		refundAmount := item.PaidAmount
		if refundAmount == 0 {
			refundAmount = item.TotalAmount
		}
		if refundAmount == 0 && item.UnitPrice > 0 {
			refundAmount = item.UnitPrice * float64(item.Quantity)
		}
		if item.Quantity > 0 && returnQuantity < item.Quantity {
			refundAmount = proportionalAmount(refundAmount, returnQuantity, item.Quantity)
		}

		remainingQuantity := item.Quantity - returnQuantity
		remainingTotalAmount := item.TotalAmount
		remainingFinalAmount := item.PaidAmount
		if returnQuantity == item.Quantity {
			remainingTotalAmount = 0
			remainingFinalAmount = 0
		} else {
			remainingTotalAmount = proportionalAmount(item.TotalAmount, remainingQuantity, item.Quantity)
			if item.PaidAmount > 0 {
				remainingFinalAmount = proportionalAmount(item.PaidAmount, remainingQuantity, item.Quantity)
			} else {
				remainingFinalAmount = remainingTotalAmount
			}
		}

		if math.Abs(order.FinalAmount-refundAmount) > 0.000001 {
			newFinalAmount := order.FinalAmount - refundAmount
			if newFinalAmount < 0 {
				newFinalAmount = 0
			}
			if err := tx.Model(&models.Order{}).Where("id = ?", item.OrderID).Update("final_amount", newFinalAmount).Error; err != nil {
				return err
			}
		}

		if returnQuantity == item.Quantity {
			if err := tx.Model(&models.OrderItems{}).Where("id = ?", item.ID).Updates(map[string]interface{}{
				"quantity":        0,
				"order_status":    "completed",
				"delivery_status": "returned",
				"total_amount":    0.0,
				"final_amount":    0.0,
				"coupon_discount": 0.0,
				"offer_discount":  0.0,
				"total_discount":  0.0,
			}).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Model(&models.OrderItems{}).Where("id = ?", item.ID).Updates(map[string]interface{}{
				"quantity":        remainingQuantity,
				"total_amount":    remainingTotalAmount,
				"final_amount":    remainingFinalAmount,
				"order_status":    "completed",
				"delivery_status": item.DeliveryStatus,
			}).Error; err != nil {
				return err
			}

			returnedItem := models.OrderItems{
				OrderID:          item.OrderID,
				ProductID:        item.ProductID,
				ProductVariantID: item.ProductVariantID,
				Size:             item.Size,
				OfferID:          item.OfferID,
				Quantity:         returnQuantity,
				UnitPrice:        item.UnitPrice,
				TotalAmount:      proportionalAmount(item.TotalAmount, returnQuantity, item.Quantity),
				DeliveryStatus:   "returned",
				OrderStatus:      "completed",
				PaymentMethod:    item.PaymentMethod,
				CouponDiscount:   proportionalAmount(item.CouponDiscount, returnQuantity, item.Quantity),
				OfferDiscount:    proportionalAmount(item.OfferDiscount, returnQuantity, item.Quantity),
				TotalDiscount:    proportionalAmount(item.TotalDiscount, returnQuantity, item.Quantity),
				PaidAmount:       proportionalAmount(item.PaidAmount, returnQuantity, item.Quantity),
			}
			if err := tx.Create(&returnedItem).Error; err != nil {
				return err
			}
		}

		if err := restoreStockForItem(tx, item, returnQuantity); err != nil {
			return fmt.Errorf("failed to restore stock for returned item: %w", err)
		}

		var remainingItemCount int64
		if err := tx.Raw(`SELECT COUNT(*) FROM order_items WHERE order_id = ? AND delivery_status != 'cancelled' AND delivery_status != 'returned'`, item.OrderID).Scan(&remainingItemCount).Error; err != nil {
			return err
		}
		if remainingItemCount == 0 {
			if err := tx.Model(&models.Order{}).Where("id = ?", item.OrderID).Updates(map[string]interface{}{
				"delivery_status": "returned",
				"order_status":    "completed",
			}).Error; err != nil {
				return err
			}
		}

		if item.PaymentMethod == "RazorPay" || item.PaymentMethod == "Wallet" {
			var walletCount int64
			if err := tx.Raw(`SELECT COUNT(*) FROM wallets WHERE user_id = ?`, userID).Scan(&walletCount).Error; err != nil {
				return err
			}
			if walletCount == 0 {
				wallet := models.Wallet{UserID: userID, Balance: refundAmount}
				if err := tx.Create(&wallet).Error; err != nil {
					return err
				}
			} else {
				var balance float64
				if err := tx.Model(&models.Wallet{}).Where("user_id = ?", userID).Pluck("balance", &balance).Error; err != nil {
					return err
				}
				wallet := models.Wallet{UserID: userID, Balance: balance + refundAmount}
				if err := tx.Model(&models.Wallet{}).Where("user_id = ?", userID).Updates(&wallet).Error; err != nil {
					return err
				}
			}
			transaction := models.WalletTransaction{
				UserID:          userID,
				Amount:          refundAmount,
				TransactionType: "Credit",
				Description:     "Refund for returning single order item",
			}
			if err := tx.Create(&transaction).Error; err != nil {
				return err
			}
		}

		now := time.Now()
		today := now.Format("2006-01-02")
		payment := models.Payments{
			UserID:        userID,
			OrderID:       item.OrderID,
			OrderItemID:   item.ID,
			TotalAmount:   refundAmount,
			PaymentDate:   today,
			PaymentType:   item.PaymentMethod,
			PaymentStatus: "refund",
			Description:   "refund for returning single order item",
		}
		if err := tx.Create(&payment).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order item return processed successfully"})
}

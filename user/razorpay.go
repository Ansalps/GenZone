package user

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/helper"
	"github.com/Ansalps/GeZOne/middleware"
	"github.com/Ansalps/GeZOne/models"
	"github.com/Ansalps/GeZOne/responsemodels"
	"github.com/gin-gonic/gin"
	"github.com/razorpay/razorpay-go"
	"gorm.io/gorm"
)

var RazorpayClient *razorpay.Client

type OrderID struct {
	Orderid string
}

func getUserCartSummary(db *gorm.DB, userID uint) (count int64, totalQuantity uint, totalAmount float64, err error) {
	if db == nil {
		db = database.DB
	}

	if err = db.Table("cart_items").
		Select("COUNT(*)").
		Joins("JOIN carts ON carts.id = cart_items.cart_id").
		Where("carts.user_id = ?", userID).
		Scan(&count).Error; err != nil {
		return 0, 0, 0, err
	}

	var qty struct {
		Total uint
	}
	if err = db.Table("cart_items").
		Select("COALESCE(SUM(cart_items.quantity), 0) AS total").
		Joins("JOIN carts ON carts.id = cart_items.cart_id").
		Where("carts.user_id = ?", userID).
		Scan(&qty).Error; err != nil {
		return 0, 0, 0, err
	}
	totalQuantity = qty.Total

	var amount struct {
		Total float64
	}
	if err = db.Table("cart_items").
		Select("COALESCE(SUM(cart_items.quantity * cart_items.unit_price), 0) AS total").
		Joins("JOIN carts ON carts.id = cart_items.cart_id").
		Where("carts.user_id = ?", userID).
		Scan(&amount).Error; err != nil {
		return 0, totalQuantity, 0, err
	}
	totalAmount = amount.Total
	return count, totalQuantity, totalAmount, nil
}

func CreateOrder(c *gin.Context) {
	razorpayKeyId := os.Getenv("RAZORPAY_KEY_ID")
	razorpayKeySecret := os.Getenv("RAZORPAY_KEY_SECRET")
	fmt.Println("--------", razorpayKeyId, razorpayKeySecret)
	RazorpayClient := razorpay.NewClient(razorpayKeyId, razorpayKeySecret)
	fmt.Println("RazorpayClient")
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
	//addressid verifying
	var tempaddress models.TempAddress
	err := c.BindJSON(&tempaddress)
	response := gin.H{
		"status":  false,
		"message": "failed to bind request",
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// Validate the content of the JSON
	if err := helper.Validate(tempaddress); err != nil {
		fmt.Println("", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"message":    err.Error(),
			"error_code": http.StatusBadRequest,
		})
		return
	}

	var count1 int64
	database.DB.Raw(`SELECT COUNT(*) FROM addresses where id = ? AND user_id = ? AND deleted_at IS NULL`, tempaddress.AddressID, userID).Scan(&count1)
	if count1 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "address_id does not exist for this particular user",
		})
		return
	}

	count, totalquantity, totalamount, err := getUserCartSummary(database.DB, userID)
	if err != nil {
		log.Println("failed to load cart summary:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate cart total"})
		return
	}

	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "false",
			"message": "cart empty, order can't be placed",
		})
		return
	}

	if totalquantity == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "false",
			"message": "cart empty, order can't be placed",
		})
		return
	}

	var Finalamount float64
	var discountamount float64
	var totalamount1 = totalamount
	if tempaddress.CouponCode != "" {
		var count2 int64
		database.DB.Raw(`SELECT COUNT(*) FROM coupons where code = ? AND deleted_at IS NULL`, tempaddress.CouponCode).Scan(&count2)
		if count2 == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "such coupon does not exists",
			})
			return
		}
		var minpurchase float64
		database.DB.Model(&models.Coupon{}).Where("code = ?", tempaddress.CouponCode).Pluck("min_purchase", &minpurchase)
		if totalamount1 > minpurchase {

			database.DB.Model(&models.Coupon{}).Where("code = ?", tempaddress.CouponCode).Pluck("discount", &discountamount)

		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": "there is a minimum purchase amount to apply the coupon",
			})
			return
		}
	}
	//added code

	//code
	Finalamount = totalamount - discountamount
	amountInPaise := int(Finalamount * 100)
	data := map[string]interface{}{
		"amount":   amountInPaise, // amount in smallest currency unit (e.g., 50000 paise = 500 INR)
		"currency": "INR",
	}

	headers := map[string]string{} // Optional headers if any

	order, err := RazorpayClient.Order.Create(data, headers)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	tempaddress1 := models.TempAddress{
		AddressID:  tempaddress.AddressID,
		CouponCode: tempaddress.CouponCode,
	}
	database.DB.Create(&tempaddress1)

	orderID := order["id"].(string)
	checkoutAmountInPaise := 0
	if rawAmount, ok := order["amount"].(float64); ok {
		checkoutAmountInPaise = int(rawAmount)
	}

	
	c.JSON(http.StatusOK, gin.H{
		"status":       "success",
		"message":      "Razorpay order created successfully",
		"amount":       checkoutAmountInPaise,
		"currency":     "INR",
		"order_id":     orderID,
		"receipt":      "receipt_" + orderID,
		"checkout_url": "/payment-option",
	})
}

type Payload struct {
	OrderID   string `json:"order_id"`
	PaymentID string `json:"payment_id"`
	Signature string `json:"signature"`
}

func verifySignature(orderID string, paymentID string, razorpaySignature string, secret string) bool {
	fmt.Println("order ID: ", orderID, " payment ID: ", paymentID)
	data := orderID + "|" + paymentID
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	generatedSignature := hex.EncodeToString(h.Sum(nil))
	fmt.Println("generatedSignature:", generatedSignature)
	fmt.Println("razorpaySignature:", razorpaySignature)
	return hmac.Equal([]byte(generatedSignature), []byte(razorpaySignature))
}

func PaymentWebhook(c *gin.Context) {
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
	var payload Payload
	if err := c.BindJSON(&payload); err != nil {
		log.Println("Error reading request body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	orderID := payload.OrderID
	paymentID := payload.PaymentID
	signature := payload.Signature
	secret := os.Getenv("RAZORPAY_KEY_SECRET")
	if secret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Razorpay secret is not configured"})
		return
	}

	if verifySignature(orderID, paymentID, signature, secret) {
		// Process the payment event
		_, _, totalamount, err := getUserCartSummary(database.DB, userID)
		if err != nil {
			log.Println("failed to load cart summary for verification:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate cart total"})
			return
		}
		var totalamount1 = totalamount
		var addressid uint
		database.DB.Raw(`SELECT address_id from temp_addresses`).Scan(&addressid)
		var couponcode string
		var count int64
		result := database.DB.Raw(`SELECT COUNT(*) FROM temp_addresses WHERE coupon_code != ''`).Scan(&count)
		if result.Error != nil {
			panic(result.Error)
		}

		if count != 0 {
			database.DB.Raw(`SELECT coupon_code from temp_addresses`).Scan(&couponcode)
		}
		var Finalamount float64
		var discountamount float64
		if couponcode != "" {
			var count2 int64
			database.DB.Raw(`SELECT COUNT(*) FROM coupons where code = ?`, couponcode).Scan(&count2)
			if count2 == 0 {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "such coupon does not exists",
				})
				return
			}
			var minpurchase float64
			database.DB.Model(&models.Coupon{}).Where("code = ?", couponcode).Pluck("min_purchase", &minpurchase)
			if totalamount1 > minpurchase {

				database.DB.Model(&models.Coupon{}).Where("code = ?", couponcode).Pluck("discount", &discountamount)

			}

		}
		var offerapplied float64
		var couponID *uint
		if couponcode != "" {
			var couponIDValue uint
			database.DB.Model(&models.Coupon{}).Where("code = ?", couponcode).Select("id").Scan(&couponIDValue)
			couponID = &couponIDValue
		}
		Finalamount = totalamount - discountamount
		order := models.Order{
			UserID:              userID,
			AddressID:           addressid,
			TotalAmount:         totalamount,
			PaymentMethod:       "RazorPay",
			OrderStatus:         "pending",
			CouponID:            couponID,
			CouponCode:          couponcode,
			TotalDiscountAmount: offerapplied + discountamount,
			FinalAmount:         Finalamount,
		}
		database.DB.Create(&order)
		var CartItems []models.CartItem
		if err := database.DB.Joins("JOIN carts ON carts.id = cart_items.cart_id").Where("carts.user_id = ?", userID).Find(&CartItems).Error; err != nil {
			log.Println("failed to load cart items after payment verification:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart items"})
			return
		}

		var ID uint

		database.DB.Raw(`SELECT id FROM orders where user_id = ? ORDER BY created_at DESC LIMIT 1`, userID).Scan(&ID)

		for _, v := range CartItems {

			if v.Quantity == 0 {
				continue
			}
			for i := 0; i < int(v.Quantity); i++ {
				var price float64
				database.DB.Model(&models.CartItem{}).Where("product_id = ?", v.ProductID).Pluck("unit_price", &price)

				var offerdiscount float64
				var coupondiscount float64
				if discountPercentage, hasOffer := helper.GetActiveOfferDiscount(database.DB, v.ProductID); hasOffer {
					offerdiscount = price * discountPercentage / 100
				}
				var totalamount float64
				database.DB.Model(&models.Order{}).Where("id = ?", ID).Pluck("total_amount", &totalamount)
				coupondiscount = (price / totalamount1) * discountamount
				coupondiscount = math.Round(coupondiscount*100) / 100
				totaldiscount := offerdiscount + coupondiscount
				paidamount := price - totaldiscount
				lineTotalAmount := price * float64(v.Quantity)
				orderItem := models.OrderItems{
					OrderID:        ID,
					ProductID:      v.ProductID,
					Quantity:       v.Quantity,
					UnitPrice:      price,
					TotalAmount:    lineTotalAmount,
					PaymentMethod:  "RazorPay",
					OfferDiscount:  offerdiscount,
					CouponDiscount: coupondiscount,
					TotalDiscount:  totaldiscount,
					PaidAmount:     paidamount,
				}

				database.DB.Create(&orderItem)

			}

		}
		now := time.Now()
		today := now.Format("2006-01-02")
		Payment := models.Payments{
			UserID:        userID,
			OrderID:       order.ID,
			TotalAmount:   Finalamount,
			PaymentDate:   today,
			PaymentType:   "RazorPay",
			PaymentStatus: "paid",
		}
		database.DB.Create(&Payment)
		var currentCart models.Cart
		if err := database.DB.Where("user_id = ?", userID).First(&currentCart).Error; err == nil {
			database.DB.Where("cart_id = ?", currentCart.ID).Delete(&models.CartItem{})
		}
		var order1 responsemodels.Order
		var address responsemodels.Address
		var orderitems1 []responsemodels.OrderItems
		database.DB.Raw(`SELECT orders.id,orders.created_at,orders.updated_at,orders.deleted_at,orders.user_id,orders.address_id,orders.total_amount,orders.offer_discount AS offer_applied,orders.coupon_code,orders.coupon_discount AS discount_amount,orders.final_amount,orders.order_status,orders.payment_method FROM orders join addresses on orders.address_id=addresses.id WHERE orders.user_id = ? ORDER BY orders.created_at desc LIMIT 1`, userID).Scan(&order1)

		var orderid uint
		database.DB.Raw(`SELECT id FROM orders WHERE user_id = ? ORDER BY created_at desc limit 1`, userID).Scan(&orderid)

		database.DB.Raw(`SELECT address_id FROM orders WHERE user_id = ? ORDER BY created_at desc limit 1`, userID).Scan(&addressid)

		database.DB.Raw(`SELECT * FROM addresses WHERE id = ?`, addressid).Scan(&address)
		order1.Address = address
		database.DB.Raw(`SELECT order_items.id,order_items.created_at,order_items.updated_at,order_items.deleted_at,order_items.order_id,order_items.product_id,products.product_name,order_items.quantity,order_items.unit_price AS price,order_items.total_amount,order_items.order_status,order_items.payment_method,order_items.coupon_discount,order_items.offer_discount,order_items.total_discount,order_items.paid_amount FROM order_items join products on order_items.product_id=products.id WHERE order_items.order_id = ? ORDER BY order_items.id`, orderid).Scan(&orderitems1)
		result = database.DB.Exec("TRUNCATE temp_addresses")
		if result.Error != nil {
			panic(result.Error)
		}
		c.JSON(http.StatusOK, gin.H{"message": "Order added successfully",
			"order":       order1,
			"order_items": orderitems1,
			"status":      "success"})

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signature"})
	}
}

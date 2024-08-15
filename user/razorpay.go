package user

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/helper"
	"github.com/Ansalps/GeZOne/middleware"
	"github.com/Ansalps/GeZOne/models"
	"github.com/Ansalps/GeZOne/responsemodels"
	"github.com/gin-gonic/gin"
	"github.com/razorpay/razorpay-go"
)

var RazorpayClient *razorpay.Client

func CreateOrder(c *gin.Context) {
	RazorpayClient := razorpay.NewClient("rzp_test_I0KzQyB0QDtMcB", "5nCtZw13gRp79G3ptqHut3Fl")
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
	fmt.Println("tempaddress.AddressID---", tempaddress.AddressID)
	fmt.Println("", userID)
	var count1 int64
	database.DB.Raw(`SELECT COUNT(*) FROM addresses where id = ? AND user_id = ? AND deleted_at IS NULL`, tempaddress.AddressID, userID).Scan(&count1)
	if count1 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "address_id does not exist for this particular user",
		})
		return
	}
	// var requestData struct {
	// 	Amount   int    `json:"amount" binding:"required"`
	// 	Currency string `json:"currency" binding:"required"`
	// 	Receipt  string `json:"receipt" binding:"required"`
	// }

	// if err := c.ShouldBindJSON(&requestData); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM cart_items WHERE user_id=? and deleted_at IS NULL`, userID).Scan(&count)
	fmt.Println("count ", count)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "false",
			"message": "cart empty, order can't be placed",
		})
		return
	}
	var totalquantity uint
	database.DB.Raw(`SELECT SUM(qty) FROM cart_items WHERE user_id=? and deleted_at IS NULL`, userID).Scan(&totalquantity)
	fmt.Println("total quantity", totalquantity)
	if totalquantity == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "false",
			"message": "cart empty, order can't be placed",
		})
		return
	}

	var totalamount float64
	database.DB.Raw("SELECT SUM(final_amount) from cart_items where user_id = ? and deleted_at IS NULL", userID).Scan(&totalamount)
	var totalamount1 float64
	database.DB.Raw("SELECT SUM(total_amount) from cart_items where user_id = ? and deleted_at IS NULL", userID).Scan(&totalamount1)
	var Finalamount float64
	var discountamount float64
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
	//fmt.Println("order from razorpay")
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

	c.JSON(http.StatusOK, order)
}

// type Order struct {
// 	gorm.Model
// 	UserID uint ` gorm:"default:1"`
// 	//OrderDate   time.Time
// 	TotalAmount float64 `gorm:"default:100.00"`
// 	// OrderStatus string `gorm:"type:varchar(10); check(order_status IN ('pending', 'delivered', 'cancelled')) ;default:'pending'" json:"order_status" validate:"required"`
// 	OrderStatus string `gorm:"type:varchar(10);check:order_status IN ('pending','shipped', 'delivered', 'cancelled','failed');default:'pending'" json:"order_status"`
// }
// type OrderRequest struct {
// 	TotalAmount float64 `gorm:"default:100.00" json:"total_amount"`
// }

//	func CreateOrder(c *gin.Context) {
//		var orderrequest OrderRequest
//		err := c.BindJSON(&orderrequest)
//		response := gin.H{
//			"status":  false,
//			"message": "failed to bind request",
//		}
//		if err != nil {
//			c.JSON(http.StatusBadRequest, response)
//			return
//		}
//		order := Order{
//			TotalAmount: orderrequest.TotalAmount,
//		}
//		DB.Create(&order)
//	}
type Payload struct {
	OrderID   string `json:"order_id"`
	PaymentID string `json:"payment_id"`
	Signature string `json:"signature"`
}

func verifySignature(orderID string, paymentID string, razorpaySignature string, secret string) bool {
	data := orderID + "|" + paymentID
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	generatedSignature := hex.EncodeToString(h.Sum(nil))
	fmt.Println(generatedSignature)
	fmt.Println(razorpaySignature)
	return generatedSignature == razorpaySignature
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
	secret := "5nCtZw13gRp79G3ptqHut3Fl"
	fmt.Println("signature", signature)
	if verifySignature(orderID, paymentID, signature, secret) {
		// Process the payment event
		fmt.Println("Payment verified:", payload)
		var totalamount float64
		database.DB.Raw("SELECT SUM(final_amount) from cart_items where user_id = ? and deleted_at IS NULL", userID).Scan(&totalamount)
		var totalamount1 float64
		database.DB.Raw("SELECT SUM(total_amount) from cart_items where user_id = ? and deleted_at IS NULL", userID).Scan(&totalamount1)
		var addressid uint
		database.DB.Raw(`SELECT address_id from temp_addresses`).Scan(&addressid)
		var couponcode string
		var count int64
		result := database.DB.Raw(`SELECT COUNT(*) FROM temp_addresses WHERE coupon_code != ''`).Scan(&count)
		if result.Error != nil {
			panic(result.Error)
		}
		fmt.Println("count printing--", count)

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
		database.DB.Raw(`SELECT SUM(discount) FROM cart_items WHERE deleted_at IS NULL`).Scan(&offerapplied)
		Finalamount = totalamount - discountamount
		order := models.Order{
			UserID:         userID,
			AddressID:      addressid,
			TotalAmount:    totalamount,
			PaymentMethod:  "RazorPay",
			OfferApplied:   offerapplied,
			CouponCode:     couponcode,
			DiscountAmount: discountamount,
			FinalAmount:    Finalamount,
		}
		database.DB.Create(&order)
		var CartItems []models.CartItems
		database.DB.Where("user_id = ?", userID).Find(&CartItems)

		var ID uint
		//database.DB.Model(&models.Order{}).Where("user_id = ?", userID).Pluck("id", &ID)
		database.DB.Raw(`SELECT id FROM orders where user_id = ? ORDER BY created_at DESC LIMIT 1`, userID).Scan(&ID)
		fmt.Println("latest order id ", ID)
		//var orderItem models.OrderItems
		for _, v := range CartItems {
			//var Product models.Product
			//database.DB.Where("id = ?", v.ProductID).First(&Product)
			//database.DB.Where("price=?",v.)
			fmt.Println("qty", v.Qty)
			if v.Qty == 0 {
				continue
			}
			for i := 0; i < int(v.Qty); i++ {
				var price float64
				database.DB.Model(&models.CartItems{}).Where("product_id = ?", v.ProductID).Pluck("price", &price)
				fmt.Println("id", ID)
				fmt.Println("price printing--", price)
				var offerdiscount float64
				var coupondiscount float64
				var hasoffer bool
				database.DB.Model(&models.Product{}).Where("id = ?", v.ProductID).Pluck("has_offer", &hasoffer)
				if hasoffer {
					var discountpercentage uint
					database.DB.Model(&models.Offer{}).Where("product_id = ?", v.ProductID).Pluck("discount_percentage", &discountpercentage)
					offerdiscount = price * float64(discountpercentage) / 100
				}
				var totalamount float64
				database.DB.Model(&models.Order{}).Where("id = ?", ID).Pluck("total_amount", &totalamount)
				coupondiscount = (price / totalamount1) * discountamount
				coupondiscount = math.Round(coupondiscount*100) / 100
				totaldiscount := offerdiscount + coupondiscount
				paidamount := price - totaldiscount
				orderItem := models.OrderItems{
					OrderID:   ID,
					ProductID: v.ProductID,
					//Qty:         v.Qty,
					Price: price,
					//TotalAmount: float64(v.Qty) * price,
					PaymentMethod:  "RazorPay",
					OfferDiscount:  offerdiscount,
					CouponDiscount: coupondiscount,
					TotalDiscount:  totaldiscount,
					PaidAmount:     paidamount,
				}
				fmt.Println("order id", orderItem.OrderID)
				fmt.Println("order item create hi")
				database.DB.Create(&orderItem)
				fmt.Println("order item create hello")
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
		database.DB.Where("user_id = ?", userID).Delete(&models.CartItems{})
		var order1 responsemodels.Order
		var address responsemodels.Address
		var orderitems1 []responsemodels.OrderItems
		database.DB.Raw(`SELECT orders.id,orders.created_at,orders.updated_at,orders.deleted_at,orders.user_id,orders.address_id,orders.total_amount,orders.offer_applied,orders.coupon_code,orders.discount_amount,orders.final_amount,orders.order_status,orders.payment_method FROM orders join addresses on orders.address_id=addresses.id WHERE orders.user_id = ? ORDER BY orders.created_at desc LIMIT 1`, userID).Scan(&order1)
		fmt.Println("-----------------")
		fmt.Println("user id ", userID)
		var orderid uint
		database.DB.Raw(`SELECT id FROM orders WHERE user_id = ? ORDER BY created_at desc limit 1`, userID).Scan(&orderid)
		fmt.Println("order id ", orderid)
		//var addressid uint
		database.DB.Raw(`SELECT address_id FROM orders WHERE user_id = ? ORDER BY created_at desc limit 1`, userID).Scan(&addressid)
		fmt.Println("address id", addressid)
		database.DB.Raw(`SELECT * FROM addresses WHERE id = ?`, addressid).Scan(&address)
		order1.Address = address
		database.DB.Raw(`SELECT order_items.id,order_items.created_at,order_items.updated_at,order_items.deleted_at,order_items.order_id,order_items.product_id,products.product_name,order_items.price,order_items.order_status,order_items.payment_method,order_items.coupon_discount,order_items.offer_discount,order_items.total_discount,order_items.paid_amount FROM order_items join products on order_items.product_id=products.id WHERE order_items.order_id = ? ORDER BY order_items.id`, orderid).Scan(&orderitems1)
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

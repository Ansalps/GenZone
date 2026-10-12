package user

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/helper"
	"github.com/Ansalps/GeZOne/middleware"
	"github.com/Ansalps/GeZOne/models"
	"github.com/Ansalps/GeZOne/requestmodels"
	"github.com/Ansalps/GeZOne/responsemodels"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func WalletListing(c *gin.Context) {
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
	var wallet responsemodels.Wallet
	database.DB.Raw(`SELECT * from wallets where user_id = ?`, userID).Scan(&wallet)
	c.JSON(http.StatusOK, gin.H{
		"data":    wallet,
		"message": "listing wallet",
	})
}

func WalletTransactionListing(c *gin.Context) {
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
	var wallettransaction []responsemodels.WalletTransaction
	database.DB.Raw(`SELECT * from wallet_transactions where user_id = ?`, userID).Scan(&wallettransaction)
	c.JSON(http.StatusOK, gin.H{
		"data":    wallettransaction,
		"message": "listing wallettransactions",
	})

}

func AddMoney(c *gin.Context) {
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
	var payload struct {
		Amount float64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "invalid request payload",
		})
		return
	}

	if payload.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "amount must be greater than zero",
		})
		return
	}

	var wallet models.Wallet
	if err := database.DB.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			wallet = models.Wallet{UserID: userID, Balance: 0}
			if createErr := database.DB.Create(&wallet).Error; createErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  false,
					"message": "failed to initialize wallet",
				})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "failed to load wallet",
			})
			return
		}
	}

	newBalance := wallet.Balance + payload.Amount
	if err := database.DB.Model(&wallet).Update("balance", newBalance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "failed to update wallet balance",
		})
		return
	}

	transaction := models.WalletTransaction{
		UserID:          userID,
		Amount:          payload.Amount,
		TransactionType: "Credit",
		Description:     "Wallet top-up",
	}
	if err := database.DB.Create(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "wallet updated but transaction log failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "money added to wallet successfully",
		"data": gin.H{
			"balance": newBalance,
		},
	})
}

func WalletOrder(c *gin.Context) {
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
	var OrderAdd requestmodels.OrderAdd
	err := c.BindJSON(&OrderAdd)
	response := gin.H{
		"status":  false,
		"message": "failed to bind request",
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// Validate the content of the JSON
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
	database.DB.Raw(`SELECT COUNT(*) FROM addresses where id = ? AND user_id = ? AND deleted_at IS NULL`, OrderAdd.AddressID, userID).Scan(&count1)
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

	if OrderAdd.CouponCode != "" {

		var count2 int64
		database.DB.Raw(`SELECT COUNT(*) FROM coupons where code = ? AND deleted_at IS NULL`, OrderAdd.CouponCode).Scan(&count2)
		if count2 == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "such coupon does not exists",
			})
			return
		}
		var minpurchase float64
		database.DB.Model(&models.Coupon{}).Where("code = ?", OrderAdd.CouponCode).Pluck("min_purchase", &minpurchase)
		if totalamount1 > minpurchase {

			database.DB.Model(&models.Coupon{}).Where("code = ?", OrderAdd.CouponCode).Pluck("discount", &discountamount)

		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "coupon cannot be applied, because there is a minimum purchase amount to apply this coupon",
			})
			return
		}
	}
	var offerapplied float64
	var couponID *uint
	if OrderAdd.CouponCode != "" {
		var couponIDValue uint
		database.DB.Model(&models.Coupon{}).Where("code = ?", OrderAdd.CouponCode).Select("id").Scan(&couponIDValue)
		couponID = &couponIDValue
	}
	Finalamount = totalamount - discountamount
	var balance float64
	database.DB.Raw(`SELECT balance FROM wallets WHERE user_id = ?`, userID).Scan(&balance)

	if balance-Finalamount < 0.00 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "order can't be placed because not enough money in wallet",
		})
		return
	}
	order := models.Order{
		UserID:              userID,
		AddressID:           OrderAdd.AddressID,
		TotalAmount:         totalamount,
		PaymentMethod:       "Wallet",
		OrderStatus:         "pending",
		CouponID:            couponID,
		CouponCode:          OrderAdd.CouponCode,
		TotalDiscountAmount: offerapplied + discountamount,
		FinalAmount:         Finalamount,
	}
	database.DB.Create(&order)
	balance = balance - Finalamount

	result := database.DB.Exec(`UPDATE wallets SET balance = ? WHERE user_id = ?`, balance, userID)
	if result.Error != nil {
		fmt.Println("--", result.Error)
	}
	transaction := models.WalletTransaction{
		UserID:          userID,
		Amount:          Finalamount,
		TransactionType: "Debit",
		Description:     "Purchase through Wallet",
	}
	database.DB.Create(&transaction)
	var CartItems []models.CartItem
	if err := database.DB.Joins("JOIN carts ON carts.id = cart_items.cart_id").Where("carts.user_id = ?", userID).Find(&CartItems).Error; err != nil {
		log.Println("failed to load cart items for wallet order:", err)
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
				PaymentMethod:  "Wallet",
				CouponDiscount: coupondiscount,
				OfferDiscount:  offerdiscount,
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
		PaymentType:   "Wallet",
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
	database.DB.Raw(`SELECT orders.id,orders.created_at,orders.updated_at,orders.deleted_at,orders.user_id,orders.address_id,orders.total_amount,orders.offer_discount AS offer_applied,orders.payment_method,orders.order_status,orders.coupon_code,orders.coupon_discount AS discount_amount,orders.final_amount FROM orders join addresses on orders.address_id=addresses.id WHERE orders.user_id = ? ORDER BY orders.created_at desc LIMIT 1`, userID).Scan(&order1)

	var orderid uint
	database.DB.Raw(`SELECT id FROM orders WHERE user_id = ? ORDER BY created_at desc limit 1`, userID).Scan(&orderid)

	var addressid uint
	database.DB.Raw(`SELECT address_id FROM orders WHERE user_id = ? ORDER BY created_at desc limit 1`, userID).Scan(&addressid)

	database.DB.Raw(`SELECT * FROM addresses WHERE id = ?`, addressid).Scan(&address)
	order1.Address = address
	database.DB.Raw(`SELECT order_items.id,order_items.created_at,order_items.updated_at,order_items.deleted_at,order_items.order_id,order_items.product_id,products.product_name,order_items.quantity,order_items.unit_price AS price,order_items.total_amount,order_items.order_status,order_items.payment_method,order_items.coupon_discount,order_items.offer_discount,order_items.total_discount,order_items.paid_amount FROM order_items join products on order_items.product_id=products.id WHERE order_items.order_id = ? ORDER BY order_items.id`, orderid).Scan(&orderitems1)
	c.JSON(http.StatusOK, gin.H{"message": "Order added successfully",
		"order":       order1,
		"order_items": orderitems1})
}

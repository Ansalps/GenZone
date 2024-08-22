package user

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/helper"
	"github.com/Ansalps/GeZOne/middleware"
	"github.com/Ansalps/GeZOne/models"
	"github.com/Ansalps/GeZOne/responsemodels"
	"github.com/gin-gonic/gin"
)

func Profile(c *gin.Context) {
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
	fmt.Println("print user id : ", userID)
	var User []responsemodels.User
	database.DB.Where("id = ?", userID).First(&User)
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "successfully retrieved user informations",
		"data": gin.H{
			"User": User,
		},
	})
}

func ProfileEdit(c *gin.Context) {
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
	fmt.Println("print user id : ", userID)
	var Profile models.ProfileEdit
	err := c.BindJSON(&Profile)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "false",
			"message": "failed to bind request",
		})
		return
	}
	// Validate the content of the JSON
	if err := helper.Validate(Profile); err != nil {
		fmt.Println("", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"message":    err.Error(),
			"error_code": http.StatusBadRequest,
		})
		return
	}
	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM users where phone = ?`, Profile.Phone).Scan(&count)
	if count != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "already registered mobile",
		})
		return
	}
	user := models.User{
		FirstName: Profile.FirstName,
		LastName:  Profile.LastName,
		Phone:     Profile.Phone,
	}
	//database.DB.Where("id = ?", userID).Updates(&user)
	database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(&user)
	c.JSON(http.StatusOK, gin.H{"message": "updated user profile"})
}

func PasswordChange(c *gin.Context) {
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
	fmt.Println("print user id : ", userID)
	var PasswordChange models.PasswordChange
	err := c.BindJSON(&PasswordChange)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "false",
			"message": "failed to bind request",
		})
		return
	}
	// Validate the content of the JSON
	if err := helper.Validate(PasswordChange); err != nil {
		fmt.Println("", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"message":    err.Error(),
			"error_code": http.StatusBadRequest,
		})
		return
	}
	if PasswordChange.Password != PasswordChange.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"message":    "password should match",
			"error_code": http.StatusBadRequest,
		})
		return
	}
	passwordchange := models.User{
		Password: PasswordChange.Password,
	}
	database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(&passwordchange)
	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
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
	fmt.Println("print user id : ", userID)
	listorder := c.Query("list_order")
	var orders []responsemodels.Order
	var address responsemodels.Address

	sql := `SELECT orders.id,orders.created_at,orders.updated_at,orders.deleted_at,orders.user_id,orders.address_id,orders.total_amount,orders.offer_applied,orders.order_status,orders.coupon_code,orders.discount_amount,orders.final_amount,orders.payment_method,addresses.user_id,addresses.country,addresses.state,addresses.street_name,addresses.district,addresses.pin_code,addresses.phone,addresses.default
	FROM orders
	JOIN addresses ON orders.address_id = addresses.id where orders.user_id = ?`
	if listorder == "" || listorder == "ASC" {
		sql += ` ORDER BY orders.id ASC`
	} else if listorder == "DSC" {
		sql += ` ORDER BY orders.id DESC`
	}
	database.DB.Raw(sql, userID).Scan(&orders)
	// database.DB.Raw(`SELECT *
	//         FROM orders
	//         JOIN addresses ON orders.address_id = addresses.id
	//         WHERE orders.user_id = ?`, userID).Scan(&Order)

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
	fmt.Println("print user id : ", userID)
	orderId := c.Param("order_id")
	listorder := c.Query("list_order")
	//var order models.Order
	// t := database.DB.Where("id = ?", orderId).Find(&order)
	// if t.Error != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"message": "order id does not exist",
	// 	})
	// }
	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM orders where id = ? AND user_id = ?`, orderId, userID).Scan(&count)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "order id does not exist for this particular user",
		})
		return
	}
	var orderitems []responsemodels.OrderItems
	//tx := database.DB.Where("order_id = ?", orderId).Find(&orderitems)
	// tx := database.DB.Raw(`SELECT order_items.id,order_items.created_at,order_items.updated_at,order_items.deleted_at,order_items.order_id,order_items.product_id,products.product_name,order_items.price,order_items.order_status FROM order_items join products on order_items.product_id=products.id WHERE order_items.order_id = ? ORDER BY order_items.id`, orderId).Scan(&orderitems)
	// if tx.Error != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"message": "order_id does not exist",
	// 	})
	// }
	sql := `SELECT order_items.id,order_items.created_at,order_items.updated_at,order_items.deleted_at,order_items.order_id,order_items.product_id,products.product_name,order_items.price,order_items.order_status,order_items.payment_method,order_items.coupon_discount,order_items.offer_discount,order_items.total_discount,order_items.paid_amount FROM order_items join products on order_items.product_id=products.id WHERE order_items.order_id = ?`
	if listorder == "" || listorder == "ASC" {
		sql += ` ORDER BY order_items.id ASC`
	} else if listorder == "DSC" {
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
	fmt.Println("print user id : ", userID)
	//userID:=c.Param("user_id")
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
	fmt.Println("print user id : ", userID)
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
	var orderstatus string
	database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Pluck("order_status", &orderstatus)
	fmt.Println("order status", orderstatus)
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

	var price float64
	database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Pluck("price", &price)
	var orderid uint
	database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Pluck("order_id", &orderid)
	fmt.Println("price :", price)
	var offerdiscount float64
	database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Pluck("offer_discount", &offerdiscount)
	fmt.Println("offer disount", offerdiscount)
	offerdiscounted := price - offerdiscount
	var totalamount float64
	database.DB.Model(&models.Order{}).Where("id = ?", orderid).Pluck("total_amount", &totalamount)
	//adding code
	var offerapplied float64
	database.DB.Model(&models.Order{}).Where("id = ?", orderid).Pluck("offer_applied", &offerapplied)
	total := totalamount + offerapplied
	var couponcode string
	database.DB.Raw(`select coupon_code from orders where id = ?`, orderid).Scan(&couponcode)
	var minpurchase float64
	database.DB.Raw(`select min_purchase from coupons where code = ?`, couponcode).Scan(&minpurchase)
	//added code
	totalamount = totalamount - offerdiscounted
	database.DB.Model(&models.Order{}).Where("id = ?", orderid).Update("total_amount", totalamount)
	fmt.Println("total amount---", totalamount)
	//adding code
	var paidamount float64
	var paymentmethod string
	// database.DB.Model(&models.)
	var coupondiscount float64
	if total-price < minpurchase {
		database.DB.Model(&models.Order{}).Where("id = ?", orderid).Update("final_amount", totalamount)
		database.DB.Model(&models.Order{}).Where("id = ?", orderid).Update("coupon_code", "")
		database.DB.Model(&models.Order{}).Where("id = ?", orderid).Update("discount_amount", 0.00)
		database.DB.Model(&models.OrderItems{}).Where("order_id = ?", orderid).Update("coupon_discount", 0.00)
		var OrderItems []models.OrderItems
		database.DB.Where("order_id = ?", orderid).Find(&OrderItems)

		for _, v := range OrderItems {
			database.DB.Model(&models.OrderItems{}).Where("id = ?", v.ID).Update("total_discount", v.OfferDiscount)
			fmt.Println("what about here?")
			//fmt.Println("payment method ---", paymentmethod)
			if v.PaymentMethod == "RazorPay" || v.PaymentMethod == "Wallet" {
				fmt.Println("is it entering here====")
				fmt.Println("price==", v.Price)
				fmt.Println("total dicount", v.TotalDiscount)
				var orderitem models.OrderItems
				database.DB.Raw(`UPDATE order_items SET paid_amount = ? WHERE id = ?`, v.Price-v.OfferDiscount, v.ID).Scan(&orderitem)
				fmt.Println("what is it?")
			}
			fmt.Println("nice===")
		}
	} else {

		database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Pluck("paid_amount", &paidamount)
		fmt.Println("paid amount ", paidamount)

		database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Pluck("payment_method", &paymentmethod)
		var totaldiscount float64
		database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Pluck("total_discount", &totaldiscount)
		if paymentmethod == "COD" {
			paidamount = price - totaldiscount
		}
		var finalamount float64
		database.DB.Model(&models.Order{}).Where("id = ?", orderid).Pluck("final_amount", &finalamount)
		finalamount = finalamount - paidamount
		database.DB.Model(&models.Order{}).Where("id = ?", orderid).Update("final_amount", finalamount)
		fmt.Println("final amount ", finalamount)

		database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Pluck("coupon_discount", &coupondiscount)
		var discountamount float64
		database.DB.Model(&models.Order{}).Where("id = ?", orderid).Pluck("discount_amount", &discountamount)
		discountamount = discountamount - coupondiscount
		database.DB.Model(&models.Order{}).Where("id = ?", orderid).Update("discount_amount", discountamount)
	}
	// if paymentmethod != "RazorPay" && paymentmethod != "Wallet" {
	// 	database.DB.Model(&models.Payments{}).Where("order_id = ?", orderid).Update("total_amount", finalamount)
	// }
	var count1 int
	database.DB.Raw(`SELECT COUNT(*) FROM order_items WHERE order_id = ? AND order_status != 'cancelled' AND order_status != 'return'`, orderid).Scan(&count1)
	fmt.Println("count1", count1)
	if count1 == 1 {
		database.DB.Model(&models.Order{}).Where("id = ?", orderid).Update("order_status", "cancelled")
		database.DB.Model(&models.Order{}).Where("id = ?", orderid).Update("final_amount", 0.00)
		database.DB.Model(&models.Order{}).Where("id = ?", orderid).Update("discount_amount", 0.00)
		database.DB.Model(&models.Order{}).Where("id = ?", orderid).Update("offer_applied", 0.00)
		//database.DB.Model(&models.OrderItems{}).Where("id = ?", orderid).Update("final_amount", 0.00)
	}
	database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Update("order_status", "cancelled")

	// var price float64
	// database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Pluck("price", &price)

	database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Pluck("payment_method", &paymentmethod)
	fmt.Println("Hello hi payment method---", paymentmethod)
	if paymentmethod == "RazorPay" || paymentmethod == "Wallet" {
		if total-price < minpurchase {
			var discount float64
			database.DB.Raw(`SELECT discount from coupons where code = ? and deleted_at is null`, couponcode).Scan((&discount))
			paidamount = offerdiscounted - discount
		}
		//database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Update("paid_amount", 0.00)
		var count int64
		database.DB.Raw(`SELECT COUNT(*) FROM wallets WHERE user_id = ?`, userID).Scan(&count)
		if count == 0 {
			wallet := models.Wallet{
				UserID:  userID,
				Balance: paidamount,
			}
			database.DB.Create(&wallet)
		} else {
			var balance float64
			database.DB.Model(&models.Wallet{}).Where("user_id = ?", userID).Pluck("balance", &balance)
			balance = balance + paidamount
			wallet := models.Wallet{
				UserID:  userID,
				Balance: balance,
			}
			database.DB.Model(&models.Wallet{}).Where("user_id = ?", userID).Updates(&wallet)
		}
		transaction := models.WalletTransaction{
			UserID:          userID,
			Amount:          paidamount,
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
			TotalAmount:   paidamount,
			PaymentDate:   today,
			PaymentType:   paymentmethod,
			PaymentStatus: "refund",
			Description:   "refund for cancelling single order item",
		}
		database.DB.Create(&payment)
	}
	database.DB.Model(&models.Order{}).Where("id = ?", orderid).Update("offer_applied", offerapplied-offerdiscount)
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
	fmt.Println("print user id : ", userID)
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
	fmt.Println("order status", orderstatus)
	if orderstatus == "delivered" {

		// var paymentmethod string
		// database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Pluck("payment_method", paymentmethod)

		var paidamount float64
		database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Pluck("paid_amount", &paidamount)
		var orderid uint
		database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Pluck("order_id", &orderid)
		var totalamount float64
		database.DB.Model(&models.Order{}).Where("id = ?", orderid).Pluck("total_amount", &totalamount)
		var price float64
		database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Pluck("price", &price)
		var offerdiscount float64
		database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Pluck("offer_discount", &offerdiscount)
		fmt.Println("offer disount", offerdiscount)
		offerdiscounted := price - offerdiscount
		fmt.Println("offer discounted", offerdiscounted)
		database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Pluck("price", &price)
		totalamount = totalamount - offerdiscounted
		fmt.Println("total--amount", totalamount)
		database.DB.Model(&models.Order{}).Where("id = ?", orderid).Update("total_amount", totalamount)
		var finalamount float64
		database.DB.Model(&models.Order{}).Where("id = ?", orderid).Pluck("final_amount", &finalamount)
		finalamount = finalamount - paidamount
		database.DB.Model(&models.Order{}).Where("id = ?", orderid).Update("final_amount", finalamount)
		var coupondiscount float64
		database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Pluck("coupon_discount", &coupondiscount)
		var discountamount float64
		database.DB.Model(&models.Order{}).Where("id = ?", orderid).Pluck("discount_amount", &discountamount)
		discountamount = discountamount - coupondiscount
		database.DB.Model(&models.Order{}).Where("id = ?", orderid).Update("discount_amount", discountamount)
		var balance float64
		database.DB.Model(&models.Wallet{}).Where("user_id = ?", userID).Pluck("balance", &balance)
		balance = balance + paidamount
		database.DB.Model(&models.Wallet{}).Where("user_id = ?", userID).Update("balance", balance)
		//database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Update("paid_amount", 0.00)
		var count1 int
		database.DB.Raw(`SELECT COUNT(*) FROM order_items WHERE order_id = ? AND order_status != 'cancelled' AND order_status != 'return'`, orderid).Scan(&count1)
		fmt.Println("count1", count1)
		if count1 == 1 {
			database.DB.Model(&models.Order{}).Where("id = ?", orderid).Update("order_status", "cancelled")
			database.DB.Model(&models.Order{}).Where("id = ?", orderid).Update("final_amount", 0.00)
			database.DB.Model(&models.Order{}).Where("id = ?", orderid).Update("discount_amount", 0.00)
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
	fmt.Println("print user id : ", userID)
	var wishlist []responsemodels.Wishlist
	database.DB.Raw(`SELECT wishlists.id,wishlists.created_at,wishlists.updated_at,wishlists.deleted_at,wishlists.user_id,wishlists.product_id,products.product_name,categories.category_name,products.description,products.image_url,products.price,products.stock,products.popular,products.size,products.has_offer,products.offer_discount_percent FROM wishlists JOIN products ON wishlists.product_id = products.id JOIN categories ON categories.id = products.category_id  WHERE wishlists.user_id = ? AND wishlists.deleted_at IS NULL`, userID).Scan(&wishlist)
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
	fmt.Println("print user id : ", userID)
	var wishlistadd models.WishlistAdd
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
	fmt.Println("print user id : ", userID)
	var wishlistadd models.WishlistAdd
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

package user

import (
	"fmt"
	"net/http"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/helper"
	"github.com/Ansalps/GeZOne/middleware"
	"github.com/Ansalps/GeZOne/models"
	"github.com/Ansalps/GeZOne/responsemodels"
	"github.com/gin-gonic/gin"
)

func CheckOut(c *gin.Context) {
	//userID := c.Param("user_id")
	//var CartItem models.CartItems
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
	var couponcheckout models.CouponCheckout
	err := c.BindJSON(&couponcheckout)
	response := gin.H{
		"status":  false,
		"message": "failed to bind request",
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, response)
		return
	}
	var coupondiscount float64
	if couponcheckout.CouponCode != "" {
		var count1 int64
		database.DB.Raw(`select count(*) from coupons where code = ? and deleted_at is null`, couponcheckout.CouponCode).Scan(&count1)
		if count1 == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "such coupon does not exist",
			})
			return
		}

		//var coupondiscount float64
		database.DB.Model(&models.Coupon{}).Where("code = ?", couponcheckout.CouponCode).Pluck("discount", &coupondiscount)
	}

	type mix struct {
		CartItem            []responsemodels.CartItems
		Totalamount         float64
		OfferApplied        float64
		CouponDiscount      float64
		CouponAppliedAmount float64
		Address             []responsemodels.Address
	}
	//var CartItem models.CartItems
	var Mix mix
	//database.DB.Where("user_id = ? AND qty != 0 AND deleted_at IS NULL", userID).Find(&Mix.CartItem)
	database.DB.Raw(`select cart_items.user_id,cart_items.product_id,products.product_name,cart_items.total_amount,cart_items.qty,cart_items.price,cart_items.discount,cart_items.final_amount from cart_items join products on cart_items.product_id = products.id where cart_items.user_id = ? and cart_items.qty != 0 and cart_items.deleted_at is null`, userID).Scan(&Mix.CartItem)
	//var totalamount float64
	//database.DB.Model(&models.CartItems{}).Where("user_id = ?", userID).Pluck("total_amount", &totalamount)
	var count int64
	database.DB.Raw("SELECT COUNT(*) from cart_items where user_id = ? and deleted_at IS NULL", userID).Scan(&count)
	if count != 0 {
		err := database.DB.Raw("SELECT SUM(final_amount) from cart_items where user_id = ? and deleted_at IS NULL", userID).Scan(&Mix.Totalamount).Error
		fmt.Println("-----------", Mix.Totalamount)
		if err != nil {
			fmt.Println("failed to execute query", err)
		}
		err = database.DB.Raw("SELECT SUM(discount) from cart_items where user_id = ? and deleted_at IS NULL", userID).Scan(&Mix.OfferApplied).Error
		if err != nil {
			fmt.Println("failed to execute query", err)
		}
	}
	var minpurchase float64
	database.DB.Raw(`select min_purchase from coupons where code = ? and deleted_at is null`, couponcheckout.CouponCode).Scan(&minpurchase)
	if Mix.Totalamount+Mix.OfferApplied < float64(minpurchase) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "coupon can be applied if you purchase with a minimum amount",
		})
		return
	}
	Mix.CouponDiscount = coupondiscount
	Mix.CouponAppliedAmount = Mix.Totalamount - coupondiscount
	//var Address []responsemodels.Address
	database.DB.Where("user_id = ?", userID).Find(&Mix.Address)

	finalResult := helper.Responses("Showing CheckOut Page", Mix, nil)
	c.JSON(http.StatusOK, finalResult)
	// c.JSON(http.StatusOK, gin.H{
	// 	"status":  true,
	// 	"message": "successfully retrieved total amount",
	// 	"data": gin.H{
	// 		"cart item":    Mix.CartItem,
	// 		"Total Amount": Mix.totalamount,
	// 		"Address":      Address,
	// 	},
	// })
}

// func CheckOutAddress(c *gin.Context) {
// 	userID := c.Param("user_id")
// 	var Address []models.Address
// 	database.DB.Where("user_id = ?", userID).Find(&Address)
// 	c.JSON(http.StatusOK, gin.H{
// 		"status":  true,
// 		"message": "successfully retrieved user informations",
// 		"data": gin.H{
// 			"Address": Address,
// 		},
// 	})
// }

func CheckOutAddressEdit(c *gin.Context) {
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
	addressID := c.Param("address_id")
	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM addresses where id = ? AND user_id = ? and deleted_at IS NULL`, addressID, userID).Scan(&count)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "no such address_id exist for this particular user",
		})
		return
	}
	var Address models.AddressAdd
	err := c.BindJSON(&Address)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "false",
			"message": "failed to bind request",
		})
	}
	//validate the content of JSON
	if err := helper.Validate(Address); err != nil {
		fmt.Println("", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"message":    err.Error(),
			"error_code": http.StatusBadRequest,
		})
		return
	}
	address := models.Address{
		//UserID:     UserID,
		Country:    Address.Country,
		State:      Address.State,
		District:   Address.District,
		StreetName: Address.StreetName,
		PinCode:    Address.PinCode,
		Phone:      Address.Phone,
		Default:    Address.Default,
	}
	database.DB.Model(&models.Address{}).Where("id = ? and user_id = ?", addressID, userID).Updates(&address)
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "Address updated successfully"})
}

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
	type mix struct {
		CartItem    []responsemodels.CartItems
		totalamount float64
	}
	//var CartItem models.CartItems
	var Mix mix
	database.DB.Where("user_id = ? AND qty != 0 AND deleted_at IS NULL", userID).Find(&Mix.CartItem)
	//var totalamount float64
	//database.DB.Model(&models.CartItems{}).Where("user_id = ?", userID).Pluck("total_amount", &totalamount)
	var count int64
	database.DB.Raw("SELECT COUNT(*) from cart_items where user_id = ? and deleted_at IS NULL", userID).Scan(&count)
	if count != 0 {
		err := database.DB.Raw("SELECT SUM(total_amount) from cart_items where user_id = ? and deleted_at IS NULL", userID).Scan(&Mix.totalamount).Error
		if err != nil {
			fmt.Println("failed to execute query", err)
		}
	}

	var Address []responsemodels.Address
	database.DB.Where("user_id = ?", userID).Find(&Address)

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "successfully retrieved total amount",
		"data": gin.H{
			"cart item":    Mix.CartItem,
			"Total Amount": Mix.totalamount,
			"Address":      Address,
		},
	})
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

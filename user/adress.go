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

// func Address(c *gin.Context) {
// 	var address []models.Address
// 	tx := database.DB.Find(&address)
// 	if tx.Error != nil {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"status":  false,
// 			"message": "failed to retrieve data from the database, or the data doesn't exists",
// 		})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{
// 		"status":  true,
// 		"message": "successfully retrieved user informations",
// 		"data": gin.H{
// 			"Addresses": address,
// 		},
// 	})
// }

func AddressList(c *gin.Context) {
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
	var Address []responsemodels.Address
	//database.DB.Where("user_id = ? and deleted_at IS NULL", userID).Find(&Address)
	sql := `SELECT * FROM addresses WHERE user_id = ? AND deleted_at IS NULL`
	if listorder == "" || listorder == "ASC" {
		sql += ` ORDER BY addresses.id ASC`
	} else if listorder == "DSC" {
		sql += ` ORDER BY addresses.id DESC`
	}
	database.DB.Raw(sql, userID).Scan(&Address)
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "successfully retrieved user informations",
		"data": gin.H{
			"address": Address,
		},
	})
}

func AddressAdd(c *gin.Context) {
	//UserID := c.Param("user_id")
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
	var Address models.AddressAdd
	err := c.BindJSON(&Address)
	response := gin.H{
		"status":  false,
		"message": "failed to bind request",
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, response)
		return
	}
	//validate the content of the JSON
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
		UserID:     userID,
		Country:    Address.Country,
		State:      Address.State,
		District:   Address.District,
		StreetName: Address.StreetName,
		PinCode:    Address.PinCode,
		Phone:      Address.Phone,
		Default:    Address.Default,
	}
	database.DB.Create(&address)
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "address added successfully"})
}

func AddressEdit(c *gin.Context) {
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
	fmt.Println("user id ", userID)
	AddressID := c.Param("address_id")
	fmt.Println("address id ", AddressID)
	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM addresses WHERE id = ? AND user_id = ? and deleted_at IS NULL`, AddressID, userID).Scan(&count)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "particular address id does not exist for this user",
		})
		return
	}
	var Address models.AddressAdd
	err := c.BindJSON(&Address)
	response := gin.H{
		"status":  false,
		"message": "failed to bind request",
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, response)
		return
	}
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
	database.DB.Model(&models.Address{}).Where("id = ?", AddressID).Updates(&address)
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "Address updated successfully"})
}

func AddressDelete(c *gin.Context) {
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
	fmt.Println("user_id", userID)
	AddressID := c.Param("address_id")
	fmt.Println("Address id : ", AddressID)
	//var Address models.Address
	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM addresses WHERE id = ? AND user_id = ? and deleted_at IS NULL`, AddressID, userID).Scan(&count)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "particular address id does not exist for this user",
		})
		return
	}
	//err := database.DB.Where("id =? AND deleted_at IS NULL", AddressID).First(&Address)
	// err := database.DB.Raw(`SELECT * FROM addresses WHERE id = ? AND user_id = ? AND deleted_at IS NULL`, AddressID, userID).Scan(&Address).Error
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"status":  false,
	// 		"message": "no such address id exists for this user",
	// 		"data":    gin.H{},
	// 	})
	// 	return
	// }
	database.DB.Where("id = ?", AddressID).Delete(&models.Address{})
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "Address deleted succesfully"})
}

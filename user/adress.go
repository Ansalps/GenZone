package user

import (
	"fmt"
	"net/http"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/helper"
	"github.com/Ansalps/GeZOne/models"
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

func AddressAdd(c *gin.Context) {
	//UserID := c.Param("user_id")
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Claims not found"})
		return
	}

	customClaims, ok := claims.(*helper.CustomClaims)
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
	AddressID := c.Param("address_id")
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
	AddressID := c.Param("address_id")
	fmt.Println("Address id : ", AddressID)
	var Address models.Address
	//err := database.DB.Where("id =? AND deleted_at IS NULL", AddressID).First(&Address)
	err := database.DB.Raw(`SELECT * FROM addresses WHERE id = ? AND deleted_at IS NULL`, AddressID).Scan(&Address).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "no such address id exists",
			"data":    gin.H{},
		})
		return
	}
	database.DB.Where("id = ?", AddressID).Delete(&models.Address{})
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "Address deleted succesfully"})
}

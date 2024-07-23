package admin

import (
	"fmt"
	"net/http"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/helper"
	"github.com/Ansalps/GeZOne/models"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var AdminLogin models.AdminLogin
	err := c.BindJSON(&AdminLogin)
	response := gin.H{
		"status":  false,
		"message": "failed to bind request",
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, response)
		return
	}
	err = helper.Validate(AdminLogin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
			"data":    gin.H{},
		})
		return
	}
	var Admin models.Admin
	tx := database.DB.Where("email =? AND deleted_at IS NULL", AdminLogin.Email).First(&Admin)
	if tx.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "invalid email or password",
			"data":    gin.H{},
		})
		return
	}
	var password string
	database.DB.Model(&models.Admin{}).Where("email = ?", AdminLogin.Email).Pluck("password", &password)
	if password != AdminLogin.Password {
		// Return success response
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "invalid email or password",
			"data":    gin.H{},
		})
		return
	}
	var id uint
	database.DB.Model(&models.Admin{}).Where("email = ?", AdminLogin.Email).Pluck("id", &id)
	token, err := helper.CreateToken("admin", AdminLogin.Email, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}
	fmt.Println("", token)
	// Set token as cookie
	c.SetCookie("jwt_token", token, 3600, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{"message": "Admin Login successful", "token": token})
}

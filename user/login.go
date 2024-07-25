package user

import (
	"fmt"
	"net/http"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/helper"
	"github.com/Ansalps/GeZOne/middleware"
	"github.com/Ansalps/GeZOne/models"
	"github.com/gin-gonic/gin"
)

func UserLogin(c *gin.Context) {
	var UserLogin models.UserLogin
	//get the json from the request
	if err := c.BindJSON(&UserLogin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "failed to process the incoming request",
			"data":    gin.H{},
		})
		return
	}
	//validate the content of the json
	err := helper.Validate(UserLogin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
			"data":    gin.H{},
		})
		return
	}
	//chekc whether the email exist on the database, if not return an error
	var count int64
	err = database.DB.Raw(`SELECT COUNT(*) FROM users WHERE email=?`, UserLogin.Email).Scan(&count).Error
	if err != nil {
		fmt.Println("failed to execute query", err)
	}
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User does not exist"})
		return
	} else {
		fmt.Println("user exists")
	}

	var password, status string
	database.DB.Model(&models.User{}).Where("email = ?", UserLogin.Email).Pluck("password", &password)
	database.DB.Model(&models.User{}).Where("email = ?", UserLogin.Email).Pluck("status", &status)
	if password != UserLogin.Password {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "invalid email or password",
			"data":    gin.H{},
		})
		return
	}
	if status == "Blocked" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "User is blocked by Admin",
			"data":    gin.H{},
		})
		return
	}
	var id uint
	database.DB.Model(&models.User{}).Where("email = ?", UserLogin.Email).Pluck("id", &id)
	token, err := middleware.CreateToken("user", UserLogin.Email, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}
	fmt.Println("", token)
	// Set token as cookie
	c.SetCookie("jwt_token", token, 3600, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{"message": "User Login successful", "token": token})

}

package public

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/helper"
	"github.com/Ansalps/GeZOne/models"
	requestmodemodels "github.com/Ansalps/GeZOne/requestmodels"

	"github.com/Ansalps/GeZOne/utils"
	"github.com/gin-gonic/gin"
)

func UserSignUp(c *gin.Context) {
	var userSignUp requestmodemodels.UserSignUp
	if err := c.BindJSON(&userSignUp); err != nil {
		fmt.Println("binding error",err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid request payload",
		})
		return
	}

	if err := helper.Validate(userSignUp); err != nil {
		fmt.Println("Validation error:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	if userSignUp.Password != userSignUp.ConfirmPassword {
		fmt.Println("Password mismatch error")
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Passwords do not match",
		})
		return
	}

	// Check existing phone
	var phoneCount int64
	database.DB.Model(&models.User{}).Where("phone = ?", userSignUp.Phone).Count(&phoneCount)
	if phoneCount > 0 {
		fmt.Println("Phone number already registered")
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Phone number is already registered",
		})
		return
	}

	// Check existing user
	var userCount int64
	database.DB.Model(&models.User{}).Where("email = ?", userSignUp.Email).Count(&userCount)
	if userCount > 0 {
		fmt.Println("Email already registered")
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "User with this email already exists",
		})
		return
	}

	// Generate & Send OTP
	otpCode := utils.GenerateOTP()
	if err := helper.SendOTPEmail(userSignUp.Email, otpCode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to send OTP email",
		})
		return
	}

	expiryTime := time.Now().Add(2 * time.Minute)

	// Save/Update OTP
	var otpCount int64
	database.DB.Model(&models.OTP{}).Where("email = ?", userSignUp.Email).Count(&otpCount)
	if otpCount > 0 {
		database.DB.Model(&models.OTP{}).Where("email = ?", userSignUp.Email).Updates(models.OTP{
			OTP:       otpCode,
			OtpExpiry: expiryTime,
		})
	} else {
		database.DB.Create(&models.OTP{
			Email:     userSignUp.Email,
			OTP:       otpCode,
			OtpExpiry: expiryTime,
		})
	}

	// Save/Update TempUser
	var tempCount int64
	database.DB.Model(&models.TempUser{}).Where("email = ?", userSignUp.Email).Count(&tempCount)
	if tempCount > 0 {
		database.DB.Model(&models.TempUser{}).Where("email = ?", userSignUp.Email).Updates(models.TempUser{
			FirstName: userSignUp.FirstName,
			LastName:  userSignUp.LastName,
			Password:  userSignUp.Password,
			Phone:     userSignUp.Phone,
		})
	} else {
		database.DB.Create(&models.TempUser{
			FirstName: userSignUp.FirstName,
			LastName:  userSignUp.LastName,
			Email:     userSignUp.Email,
			Password:  userSignUp.Password,
			Phone:     userSignUp.Phone,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "OTP generated and sent successfully",
	})
}

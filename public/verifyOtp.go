package public

import (
	"net/http"
	"time"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/models"
	requestmodemodels "github.com/Ansalps/GeZOne/requestmodels"

	"github.com/gin-gonic/gin"
)

func VerifyOTPHandler(c *gin.Context) {
	email := c.Query("email") // Get email from URL Query (?email=...)
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Email query parameter is required"})
		return
	}

	var verifyOTP requestmodemodels.VerifyOTP
	if err := c.BindJSON(&verifyOTP); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid request body"})
		return
	}

	// Get OTP record from DB
	var otpRecord models.OTP
	if err := database.DB.Where("email = ?", email).First(&otpRecord).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "Invalid or expired OTP"})
		return
	}

	// Check validity
	if verifyOTP.Otp != otpRecord.OTP || time.Now().After(otpRecord.OtpExpiry) {
		c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "Invalid or expired OTP"})
		return
	}

	// Retrieve TempUser
	var tempUser models.TempUser
	if err := database.DB.Where("email = ?", email).First(&tempUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Temp user data not found"})
		return
	}

	// Migrate TempUser to User
	newUser := models.User{
		FirstName: tempUser.FirstName,
		LastName:  tempUser.LastName,
		Email:     tempUser.Email,
		Password:  tempUser.Password,
		Phone:     tempUser.Phone,
	}
	if err := database.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to create user account"})
		return
	}

	// Track Login Method
	userLoginMethod := models.UserLoginMethod{
		UserLoginMethodEmail: email,
		LoginMethod:          "Manual",
	}
	database.DB.Create(&userLoginMethod)

	// Cleanup Temp records
	database.DB.Where("email = ?", email).Delete(&models.OTP{})
	database.DB.Where("email = ?", email).Delete(&models.TempUser{})

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Account created successfully",
	})
}

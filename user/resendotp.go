package user

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/helper"
	"github.com/Ansalps/GeZOne/models"
	"github.com/gin-gonic/gin"
)

func ResendOtp(c *gin.Context) {
	Email := c.Param("email")
	Otp := helper.GenerateOTP()
	//var VerifyOTP models.VerifyOTP
	//VerifyOTP.Otp=Otp
	// Send OTP via email
	fmt.Println("", Otp)
	err := helper.SendOTPEmail(Email, Otp)
	// if err != nil {
	// 	c.JSON(http.StatusOK, gin.H{"message": "Failed to send otp to mail"})
	// } else {
	// 	c.JSON(http.StatusOK, gin.H{"message": "Otp generated successfully"})
	// }
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Failed to send otp to mail"})
		return
	} else {
		Time := time.Now().Add(1 * time.Minute)
		otp := models.OTP{
			Email:     Email,
			OTP:       Otp,
			OtpExpiry: Time,
		}
		//database.DB.Create(&otp)
		database.DB.Model(&models.OTP{}).Where("email = ?", Email).Updates(&otp)
		c.JSON(http.StatusOK, gin.H{"message": "Otp generated successfully"})
		return
	}
}

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

// var Otp1 string
//var User models.User

// sign up
func UserSignUp(c *gin.Context) {
	var UserSignUp models.UserSignUp
	err := c.BindJSON(&UserSignUp)
	response := gin.H{
		"status":  false,
		"message": "failed to bind request",
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// Validate the content of the JSON
	if err := helper.Validate(UserSignUp); err != nil {
		fmt.Println("", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"message":    err.Error(),
			"error_code": http.StatusBadRequest,
		})
		return
	}
	var count4 int64
	database.DB.Raw(`SELECT COUNT(*) FROM users where phone = ?`, UserSignUp.Phone).Scan(&count4)
	if count4 != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "already registered mobile",
		})
		return
	}
	if UserSignUp.Password != UserSignUp.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"message":    "password should match",
			"error_code": http.StatusBadRequest,
		})
		return
	}

	var count3 int64
	err = database.DB.Raw(`SELECT COUNT(*) FROM user_login_methods WHERE user_login_method_email=?`, UserSignUp.Email).Scan(&count3).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "failed to retreive information from the database",
			"data":    gin.H{},
		})
		return
	}
	if count3 != 0 {
		var loginmethod string
		database.DB.Model(&models.UserLoginMethod{}).Where("user_login_method_email = ?", UserSignUp.Email).Pluck("login_method", &loginmethod)
		if loginmethod == "Google Authentication" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "please log in through google authentication",
				"data":    gin.H{},
			})
			return
		}
	}

	var count int64
	err = database.DB.Raw(`SELECT COUNT(*) FROM users WHERE email=?`, UserSignUp.Email).Scan(&count).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "failed to retreive information from the database",
			"data":    gin.H{},
		})
		return
	}
	if count == 0 {
		// Generate OTP
		fmt.Println("print email")
		Otp1 := helper.GenerateOTP()

		// Send OTP via email
		err := helper.SendOTPEmail(UserSignUp.Email, Otp1)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"message": "Failed to send otp to mail"})
			return
		} else {
			Time := time.Now().Add(2 * time.Minute)
			var count4 int64
			database.DB.Raw(`SELECT COUNT(*) FROM otps WHERE email=?`, UserSignUp.Email).Scan(&count4)
			if count4 > 0 {
				database.DB.Model(&models.OTP{}).Where("email = ?", UserSignUp.Email).Updates(models.OTP{OTP: Otp1, OtpExpiry: Time})
			} else {
				otp := models.OTP{
					Email:     UserSignUp.Email,
					OTP:       Otp1,
					OtpExpiry: Time,
				}
				database.DB.Create(&otp)
			}

			var count1 int64
			database.DB.Raw(`SELECT COUNT(*) FROM temp_users WHERE email=?`, UserSignUp.Email).Scan(&count1)
			if count1 > 0 {
				database.DB.Model(&models.TempUser{}).Where("email = ?", UserSignUp.Email).Updates(models.TempUser{FirstName: UserSignUp.FirstName, LastName: UserSignUp.LastName, Password: UserSignUp.Password, Phone: UserSignUp.Phone})
			} else {
				User := models.TempUser{
					FirstName: UserSignUp.FirstName,
					LastName:  UserSignUp.LastName,
					Email:     UserSignUp.Email,
					Password:  UserSignUp.Password,
					Phone:     UserSignUp.Phone,
				}
				database.DB.Create(&User)
			}

			c.JSON(http.StatusOK, gin.H{"message": "Otp generated successfully"})
			return
		}

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "user already exists",
			"data":    gin.H{},
		})
		return

	}

}

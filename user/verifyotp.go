package user

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/models"
	"github.com/gin-gonic/gin"
)

func VerifyOTPHandler(c *gin.Context) {
	fmt.Println("HI")
	Email := c.Param("email")
	var VerifyOTP models.VerifyOTP
	if err := c.BindJSON(&VerifyOTP); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("", VerifyOTP.Otp)
	fmt.Println("HELLO")
	var otp string
	database.DB.Model(&models.OTP{}).Where("email = ?", Email).Pluck("otp", &otp)
	fmt.Println("", otp)
	var otptime time.Time
	database.DB.Model(&models.OTP{}).Where("email = ?", Email).Pluck("otp_expiry", &otptime)
	if VerifyOTP.Otp != otp || time.Now().After(otptime) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired OTP"})
		fmt.Println("hi hello")
		return
	} else {
		// var User models.User
		// database.DB.Model(&models.TempUser{}).Select("first_name,last_name,email,password,phone").Where("email = ?", Email).Scan(&User)
		var tempUser models.TempUser
		if err := database.DB.Model(&models.TempUser{}).
			Select("first_name, last_name, email, password, phone").
			Where("email = ?", Email).
			Scan(&tempUser).Error; err != nil {
			// Handle error
			panic(err)
		}
		newUser := models.User{
			FirstName: tempUser.FirstName,
			LastName:  tempUser.LastName,
			Email:     tempUser.Email,
			Password:  tempUser.Password,
			Phone:     tempUser.Phone,
		}
		database.DB.Create(&newUser)
		//var UserLoginMethod models.UserLoginMethod
		UserLoginMethod := models.UserLoginMethod{
			UserLoginMethodEmail: Email,
			LoginMethod:          "Manual",
		}
		//fmt.Println("this will work")

		//fmt.Println("User creates")
		var count int64
		err := database.DB.Raw(`SELECT COUNT(*) FROM user_login_methods WHERE user_login_method_email=? group by user_login_method_email,login_method`, Email).Scan(&count).Error
		if err != nil {
			fmt.Println("failed to execute query", err)
		}
		if count == 0 {
			database.DB.Create(&UserLoginMethod)
			//c.JSON(http.StatusBadRequest, gin.H{"error": "User does not exist"})
		} else {
			fmt.Println("user exists")
		}
		//database.DB.Create(&UserLoginMethod)
		database.DB.Where("email = ?", Email).Delete(&models.OTP{})
		database.DB.Where("email = ?", Email).Delete(&models.TempUser{})
		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "created a new user",
			"data":    gin.H{},
		})

		//c.JSON(http.StatusOK, gin.H{"message": "Signup successful"})
	}

}

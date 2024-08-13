package database

import (
	"fmt"

	"github.com/Ansalps/GeZOne/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Initialize() {
	var err error
	dsn := "postgres://postgres:123@localhost:5432/genzone"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("connection failed due to ", err)
	}
}

func AutoMigrate() {
	DB.AutoMigrate(&models.Category{})
	DB.AutoMigrate(&models.Product{})
	DB.AutoMigrate(&models.Admin{})
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.OTP{})
	DB.AutoMigrate(&models.TempUser{})
	DB.AutoMigrate(&models.UserLoginMethod{})
	DB.AutoMigrate(&models.Address{})
	DB.AutoMigrate(&models.CartItems{})
	DB.AutoMigrate(&models.Order{})
	DB.AutoMigrate(&models.OrderItems{})
	DB.AutoMigrate(&models.Payments{})
	DB.AutoMigrate(&models.TempAddress{})
	DB.AutoMigrate(&models.Wallet{})
	DB.AutoMigrate(&models.Wishlist{})
	DB.AutoMigrate(&models.Coupon{})
	DB.AutoMigrate(&models.Offer{})
	DB.AutoMigrate(&models.SalesReportItem{})
	DB.AutoMigrate(&models.WalletTransaction{})
}

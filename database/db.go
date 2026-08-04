package database

import (
	"fmt"
	"log"
	"os"

	"github.com/Ansalps/GeZOne/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Initialize() {

	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	// Reads your string variable
	dsn := os.Getenv("DATABASE_URL")
	fmt.Println("Your DSN is:", dsn)
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
	DB.AutoMigrate(&models.CartItem{})
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
	DB.AutoMigrate(&models.Invoice{})
	admin:=models.Admin{
		Email: "admin@example.com",
		Password: "admin",
	}
	DB.Model(&models.Admin{}).Create(&admin)
}

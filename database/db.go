package database

import (
	"fmt"
	"log"
	"os"

	"github.com/Ansalps/GeZOne/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Initialize() {
	// Load the .env file if present; some local/test environments do not have one.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found; continuing with environment variables or local sqlite fallback")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Println("DATABASE_URL not set; using local sqlite database for development/testing")
		var err error
		DB, err = gorm.Open(sqlite.Open("genzone.db"), &gorm.Config{})
		if err != nil {
			fmt.Println("sqlite fallback connection failed due to ", err)
		}
		return
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("connection failed due to ", err)
	}
}

func AutoSeed() error {
	return Seed(DB)
}
func AutoMigrate() {
	DB.AutoMigrate(&models.Category{})
	DB.AutoMigrate(&models.Product{})
	DB.AutoMigrate(&models.ProductVariant{})
	DB.AutoMigrate(&models.Admin{})
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.OTP{})
	DB.AutoMigrate(&models.TempUser{})
	DB.AutoMigrate(&models.UserLoginMethod{})
	DB.AutoMigrate(&models.Address{})
	DB.AutoMigrate(&models.Cart{})
	DB.AutoMigrate(&models.CartItem{})
	DB.AutoMigrate(&models.Order{})
	DB.AutoMigrate(&models.OrderItems{})
	DB.AutoMigrate(&models.Payments{})
	DB.AutoMigrate(&models.Wallet{})
	DB.AutoMigrate(&models.Wishlist{})
	DB.AutoMigrate(&models.Coupon{})
	DB.AutoMigrate(&models.Offer{})
	DB.AutoMigrate(&models.SalesReportItem{})
	DB.AutoMigrate(&models.WalletTransaction{})
	DB.AutoMigrate(&models.Invoice{})

}

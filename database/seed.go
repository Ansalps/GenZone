package database

import (
	"github.com/Ansalps/GeZOne/models"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {

	if err:= seedAdmins(db); err!=nil{
		return err
	}
	if err := seedCategories(db); err != nil {
		return err
	}

	if err := seedProducts(db); err != nil {
		return err
	}

	return nil
}

func seedAdmins(db *gorm.DB) error{
	admins:= []models.Admin{
		{
			Email: "admin@example.com",
			Password: "admin1234",
		},
	}
	for _, admin := range admins {
		if err := db.
			Where("email = ?", admin.Email).
			FirstOrCreate(&admin).
			Error; err != nil {
			return err
		}
	}
	return nil
}

func seedCategories(db *gorm.DB) error {
	topLevelCategories := []models.Category{
		{
			CategoryName: "Men's Clothing",
			Description:  "Men's clothing",
			ImageURL:     "https://genzone-public-assets.s3.ap-south-1.amazonaws.com/categories/83937e62-705d-46cb-8972-bad0d66c4f2f.jpeg",
		},
		{
			CategoryName: "Men's Footwear",
			Description:  "Men's footwear",
			ImageURL:     "https://genzone-public-assets.s3.ap-south-1.amazonaws.com/categories/04b002bb-b405-4dc3-9afa-ec1e0c2146c5.avif",
		},
		{
			CategoryName: "Women's Clothing",
			Description:  "Women's clothing",
			ImageURL:     "https://genzone-public-assets.s3.ap-south-1.amazonaws.com/categories/01e90844-636e-4f3f-a1a2-108161bc3973.webp",
		},
		{
			CategoryName: "Women's Footwear",
			Description:  "Women's footwear",
			ImageURL:     "https://genzone-public-assets.s3.ap-south-1.amazonaws.com/categories/2e19e85c-86ed-49bb-87b1-944826e85bb7.jpeg",
		},
		{
			CategoryName: "Kids' Clothing",
			Description:  "Kids' clothing",
			ImageURL:     "https://genzone-public-assets.s3.ap-south-1.amazonaws.com/categories/913fc20e-af3d-4c8f-9af5-dca0f9f4c424.jpeg",
		},
		{
			CategoryName: "Kids' Footwear",
			Description:  "Kids' footwear",
			ImageURL:     "https://genzone-public-assets.s3.ap-south-1.amazonaws.com/categories/018053ba-367d-4733-be2b-c60aa3002e0c.jpeg",
		},
	}

	for _, category := range topLevelCategories {
		if err := db.
			Where("category_name = ?", category.CategoryName).
			FirstOrCreate(&category).
			Error; err != nil {
			return err
		}
	}

	return nil
}

func seedProducts(db *gorm.DB) error {
	var mensClothing models.Category

	if err := db.
		Where("category_name = ?", "Men's Clothing").
		First(&mensClothing).
		Error; err != nil {
		return err
	}

	product := models.Product{
		ProductName: "Classic T-Shirt",
		Description: "Classic cotton t-shirt for everyday wear.",
		ImageURL:    "https://genzone-public-assets.s3.ap-south-1.amazonaws.com/categories/f3ea3094-2d9a-4796-874d-550d511ef772.webp",
		Price:       999.00,
		Popular:     true,
		CategoryID:  mensClothing.ID,
	}

	if err := db.
		Where("product_name = ?", product.ProductName).
		FirstOrCreate(&product).
		Error; err != nil {
		return err
	}

	variants := []models.ProductVariant{
		{
			ProductID: product.ID,
			Size:      "Small",
			Stock:     10,
		},
		{
			ProductID: product.ID,
			Size:      "Medium",
			Stock:     20,
		},
		{
			ProductID: product.ID,
			Size:      "Large",
			Stock:     15,
		},
	}

	for _, variant := range variants {
		if err := db.
			Where(
				"product_id = ? AND size = ?",
				variant.ProductID,
				variant.Size,
			).
			FirstOrCreate(&variant).
			Error; err != nil {
			return err
		}
	}

	return nil
}

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
			CategoryName: "Men",
			Description:  "Men's fashion and accessories",
			ImageURL:     "/images/categories/men.jpg",
		},
		{
			CategoryName: "Women",
			Description:  "Women's fashion and accessories",
			ImageURL:     "/images/categories/women.jpg",
		},
		{
			CategoryName: "Kids",
			Description:  "Kids' fashion and accessories",
			ImageURL:     "/images/categories/kids.jpg",
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

	// Get parent categories
	var men, women, kids models.Category

	if err := db.Where("category_name = ?", "Men").First(&men).Error; err != nil {
		return err
	}

	if err := db.Where("category_name = ?", "Women").First(&women).Error; err != nil {
		return err
	}

	if err := db.Where("category_name = ?", "Kids").First(&kids).Error; err != nil {
		return err
	}

	// Create child categories
	childCategories := []models.Category{
		{
			CategoryName: "Men's Clothing",
			Description:  "Men's clothing",
			ParentID:     &men.ID,
			ImageURL:     "/images/categories/mens-clothing.jpg",
		},
		{
			CategoryName: "Men's Footwear",
			Description:  "Men's footwear",
			ParentID:     &men.ID,
			ImageURL:     "/images/categories/mens-footwear.jpg",
		},
		{
			CategoryName: "Women's Clothing",
			Description:  "Women's clothing",
			ParentID:     &women.ID,
			ImageURL:     "/images/categories/womens-clothing.jpg",
		},
		{
			CategoryName: "Women's Footwear",
			Description:  "Women's footwear",
			ParentID:     &women.ID,
			ImageURL:     "/images/categories/womens-footwear.jpg",
		},
		{
			CategoryName: "Kids' Clothing",
			Description:  "Kids' clothing",
			ParentID:     &kids.ID,
			ImageURL:     "/images/categories/kids-clothing.jpg",
		},
		{
			CategoryName: "Kids' Footwear",
			Description:  "Kids' footwear",
			ParentID:     &kids.ID,
			ImageURL:     "/images/categories/kids-footwear.jpg",
		},
	}

	for _, category := range childCategories {
		if err := db.
			Where(
				"category_name = ? AND parent_id = ?",
				category.CategoryName,
				*category.ParentID,
			).
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
		ImageURL:    "/images/products/classic-tshirt.jpg",
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

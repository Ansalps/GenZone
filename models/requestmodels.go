package models

import "gorm.io/gorm"

type AdminLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Category struct {
	gorm.Model
	ID           uint   `gorm:"primary key" json:"id"`
	CategoryName string `json:"category_name"`
	Description  string `json:"category_description"`
	ImageUrl     string `json:"category_imageUrl"`
}

type Product struct {
	gorm.Model
	ID          uint     `gorm:"primary key" json:"id"`
	CategoryID  uint     `json:"category_id"`
	Category    Category `gorm:"foriegnkey:CategoryID;references:ID"`
	ProductName string   `json:"product_name"`
	Description string   `json:"product_description"`
	ImageUrl    string   `json:"product_imageUrl"`
	Price       float64  `gorm:"type:decimal(10,2)" json:"price"`
	Stock       int      `json:"stock"`
	Popular     bool     `gorm:"type:boolean;default:false" json:"popular"`
	Size        string   `gorm:"type:varchar(10); check(size IN ('Medium', 'Small', 'Large'))" json:"size"`
}

//user
type User struct {
}

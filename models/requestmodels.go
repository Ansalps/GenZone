package models

import "gorm.io/gorm"

type AdminLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AdminCategory struct {
	gorm.Model
	CategoryID   uint   `gorm:"primaryKey" json:"category_id"`
	CategoryName string `json:"name"`
	Description  string `json:"description"`
	ImageUrl     string `json:"imageUrl"`
}

type AdminProduct struct {
	gorm.Model
	ProductID   uint          `gorm:"primaryKey"`
	CategoryID  uint          `json:"category_id"`
	Category    AdminCategory `gorm:"foreignKey:CategoryID"`
	ProductName string        `json:"name"`
	Description string        `json:"description"`
	Price       float64       `gorm:"type:decimal(10,2)" json:"price"`
	Stock       int           `json:"stock"`
	Popular     bool          `gorm:"type:boolean;default:false" json:"popular"`
	Size        string        `gorm:"type:varchar(10);default:'Medium' check(size IN ('Medium', 'Small', 'Large'))" json:"size"`
}

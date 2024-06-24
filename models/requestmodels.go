package models

import "gorm.io/gorm"

type AdminLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AdminCategory struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageUrl    string `json:"imageUrl"`
}

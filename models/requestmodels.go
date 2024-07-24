package models

type AdminLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password"`
}
type ProductAdd struct {
	//ID         uint ` json:"id"`
	CategoryID uint `json:"category_id" validate:"required"`
	//Category    Category `gorm:"foriegnkey:CategoryID;references:ID"`
	ProductName string  `json:"product_name" validate:"required"`
	Description string  `json:"product_description" validate:"required"`
	ImageUrl    string  `json:"product_imageUrl" validate:"required"`
	Price       float64 ` json:"price" validate:"required"`
	Stock       uint    `json:"stock"`
	Popular     bool    `json:"popular" validate:"required"`
	Size        string  ` json:"size" validate:"required"`
}
type ProductEdit struct {
	//ID         uint ` json:"id"`
	CategoryID uint `json:"category_id" validate:"required"`
	//Category    Category `gorm:"foriegnkey:CategoryID;references:ID"`
	ProductName string  `json:"product_name" validate:"required"`
	Description string  `json:"product_description" validate:"required"`
	ImageUrl    string  `json:"product_imageUrl" validate:"required"`
	Price       float64 ` json:"price" validate:"required"`
	Stock       int     `json:"stock"`
	Popular     bool    `json:"popular" validate:"required"`
	Size        string  ` json:"size" validate:"required"`
}
type CategoryEdit struct {
	//ID           uint   `gorm:"primary key" json:"id"`
	CategoryName string ` gorm:"unique" json:"category_name" validate:"required"`
	Description  string `json:"category_description" validate:"required"`
	ImageUrl     string `json:"category_imageUrl" validate:"required"`
}
type UserSignUp struct {
	FirstName       string `validate:"required,excludesall= " json:"name"`
	LastName        string `validate:"required,nameOrInitials" json:"last_name"`
	Email           string `gorm:"unique" validate:"required,email" json:"email"`
	Password        string `validate:"required,min=8,password" json:"password"`
	ConfirmPassword string `validate:"required" json:"confirmpassword"`
	Phone           string `json:"phone" validate:"required,numeric,len=10"`
}

type VerifyOTP struct {
	Otp string `json:"otp"`
}

type UserLogin struct {
	Email    string `gorm:"unique" validate:"required,email" json:"email"`
	Password string `validate:"required" json:"password"`
}

type BlockUser struct {
	UserID uint `validate:"required" json:"user_id"`
}

type AddressAdd struct {
	//UserID     uint   `validate:"required"`
	//User       User   `gorm:"foriegnkey:UserID;references:ID"`
	Country    string `json:"country" validate:"required"`
	State      string `json:"state" validate:"required"`
	District   string `json:"district" validate:"required"`
	StreetName string `json:"street_name" validate:"required"`
	PinCode    string `json:"pin_code" validate:"required,numeric"`
	Phone      string `json:"phone" validate:"required,numeric,len=10"`
	Default    bool   `json:"Default" `
}

type CartAdd struct {
	ProductID string `json:"product_id" validate:"required,numeric"`
}

type OrderAdd struct {
	AddressID uint `json:"address_id" validate:"required"`
}

type ProfileEdit struct {
	FirstName string `validate:"required" json:"name"`
	LastName  string `validate:"required" json:"last_name"`
	//Email           string `gorm:"unique" validate:"required,email" json:"email"`
	//Password        string `validate:"required" json:"password"`
	//ConfirmPassword string `validate:"required" json:"confirmpassword"`
	Phone string `json:"phone" validate:"required,numeric,len=10"`
}

type PasswordChange struct {
	Password        string `validate:"required" json:"password"`
	ConfirmPassword string `validate:"required" json:"confirmpassword"`
}

type CancelOrder struct {
	OrderStatus string `json:"order_status" validate:"required"`
}

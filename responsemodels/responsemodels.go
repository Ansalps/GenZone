package responsemodels

import "gorm.io/gorm"

// type Category struct {
// 	gorm.Model
// 	//ID           uint   `gorm:"primary key" json:"id"`
// 	CategoryName string ` gorm:"unique" json:"category_name" validate:"required"`
// 	Description  string `json:"category_description" validate:"required"`
// 	ImageUrl     string `json:"category_imageUrl" validate:"required"`
// }

type Product struct {
	gorm.Model
	//ID          uint     `gorm:"primary key" json:"id"`
	CategoryID   uint   `json:"category_id" validate:"required"`
	CategoryName string `json:"category_name" validate:"required"`
	//Category     Category `gorm:"foriegnkey:CategoryID;references:ID" json:"category,omitempty"`
	ProductName string  `json:"product_name" validate:"required"`
	Description string  `json:"product_description" validate:"required"`
	ImageUrl    string  `json:"product_imageUrl" validate:"required"`
	Price       float64 `gorm:"type:decimal(10,2)" json:"price" validate:"required"`
	Stock       int     `json:"stock"`
	Popular     bool    `gorm:"type:boolean;default:false" json:"popular" validate:"required"`
	Size        string  `gorm:"type:varchar(10); check(size IN ('Medium', 'Small', 'Large'))" json:"size" validate:"required"`
}

type CartItems struct {
	//gorm.Model
	UserID string `validate:"required,numeric"`
	//User   User   `gorm:"foriegnkey:UserID;references:ID"`
	// CartID      string  `validate:"required,numeric"`
	// Cart        Cart    `gorm:"foriegnkey:CartID;references:ID"`
	ProductID   string `validate:"required,numeric"`
	ProductName string `json:"product_name" validate:"required"`
	//Product     Product `gorm:"foriegnkey:ProductID;references:ID"`
	TotalAmount float64 `gorm:"type:decimal(10,2);default:0.00" json:"price" validate:"required"`
	Qty         uint    `gorm:"default:0"`
}

type Address struct {
	gorm.Model
	UserID string `validate:"required,numeric"`
	//User       User   `gorm:"foriegnkey:UserID;references:ID"`
	Country    string `validate:"required"`
	State      string `validate:"required"`
	District   string `validate:"required"`
	StreetName string `validate:"required"`
	PinCode    string `validate:"required,numeric"`
	Phone      string `validate:"required,numeric,len=10"`
	Default    bool   `gorm:"default:false" validate:"required"`
}

type Order struct {
	gorm.Model
	UserID string `validate:"required,numeric"`
	//OrderDate   time.Time
	AddressID uint
	//Address     Address `gorm:"foriegnkey:AddressID;references:ID"`
	TotalAmount float64
	OrderStatus string `gorm:"type:varchar(10); check(status IN ('pending', 'delivered', 'cancelled')) ;default:'pending'" json:"order_status" validate:"required"`
}

type User struct {
	//gorm.Model
	//ID        uint   `gorm:"primary key" json:"id"`
	FirstName string `validate:"required"`
	LastName  string `validate:"required"`
	Email     string `gorm:"unique" validate:"required"`
	Password  string `validate:"required"`
	Phone     string `json:"phone" validate:"required,numeric,len=10"`
	Status    string `gorm:"type:varchar(10); check(status IN ('Active', 'Blocked', 'Deleted')) ;default:'Active'" json:"status" validate:"required"`
}

package models

import (
	"time"

	"gorm.io/gorm"
)

type Admin struct {
	gorm.Model
	Email    string `json:"email"`
	Password string `json:"password"`
}


type TempUser struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
	Phone     string
}
type OTP struct {
	Email     string `gorm:"primary key" json:"email"`
	OTP       string
	OtpExpiry time.Time
}

type UserLoginMethod struct {
	UserLoginMethodEmail string
	LoginMethod          string
}

type Address struct {
	gorm.Model
	UserID     uint   `validate:"required"`
	User       User   `gorm:"foriegnkey:UserID;references:ID"`
	Country    string `validate:"required"`
	State      string `validate:"required"`
	District   string `validate:"required"`
	StreetName string `validate:"required"`
	PinCode    string `validate:"required,numeric"`
	Phone      string `validate:"required,numeric,len=10"`
	Default    bool   `gorm:"default:false" validate:"required"`
}

type TempAddress struct {
	AddressID  uint   `json:"address_id" validate:"required"`
	CouponCode string `json:"coupon_code"`
}


type User struct {
	gorm.Model
	FirstName string `validate:"required"`
	LastName  string `validate:"required"`
	Email     string `gorm:"unique" validate:"required"`
	Password  string `validate:"required"`
	Phone     string `json:"phone" validate:"required,numeric,len=10"`
	Status    string `gorm:"type:varchar(10); check(status IN ('Active', 'Blocked', 'Deleted')) ;default:'Active'" json:"status" validate:"required"`
}

type Category struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	CategoryName string `gorm:"unique;not null" json:"category_name" validate:"required"`
	Description  string `json:"category_description" validate:"required"`
	ImageURL     string `json:"category_image_url" validate:"required"`

	Products []Product
}

type Product struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	CategoryID uint     `json:"category_id"`
	Category   Category `json:"category,omitempty"`

	ProductName string `json:"product_name" validate:"required"`
	Description string `json:"product_description" validate:"required"`
	ImageURL    string `json:"product_image_url" validate:"required"`

	Price float64 `gorm:"type:decimal(10,2)" json:"price"`

	Stock uint `json:"stock"`

	Popular bool `gorm:"default:false" json:"popular"`

	Size string `gorm:"type:varchar(10);check:size IN ('Small','Medium','Large')" json:"size"`
}

type Offer struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	ProductID          uint `gorm:"not null"`
	DiscountPercentage float64 `gorm:"not null"`
	StartAt time.Time
	EndAt time.Time
}

type Cart struct {
	gorm.Model

	UserID uint
	User   User

	CartItems []CartItem
}

type CartItem struct {
	gorm.Model

	CartID uint
	Cart   Cart

	ProductID uint
	Product   Product

	Quantity  uint
	UnitPrice float64 `gorm:"type:decimal(10,2)"`
}

type Order struct {
	gorm.Model
	UserID uint `validate:"required"`
	//OrderDate   time.Time
	AddressID   uint
	Address     Address `gorm:"foriegnkey:AddressID;references:ID"`
	TotalAmount float64
	// OrderStatus string `gorm:"type:varchar(10); check(order_status IN ('pending', 'delivered', 'cancelled')) ;default:'pending'" json:"order_status" validate:"required"`
	PaymentMethod  string  `gorm:"type:varchar(10); check(order_status IN ('COD', 'RazorPay')) ;default:'COD'" json:"payment_method" validate:"required"`
	OrderStatus    string  `gorm:"type:varchar(10);check:order_status IN ('pending','shipped', 'delivered', 'cancelled','failed');default:'pending'" json:"order_status" validate:"required,oneof=pending delivered shipped cancelled failed"`
	OfferApplied   float64 `gorm:"default:0.00"`
	CouponCode     string
	DiscountAmount float64 `gorm:"type:decimal(10,2);default:0.00"`
	FinalAmount    float64 `gorm:"type:decimal(10,2);not null"`
}

type OrderItems struct {
	gorm.Model
	OrderID   uint    `validate:"required"`
	Order     Order   `gorm:"foriegnkey:OrderID;references:ID"`
	ProductID uint  `validate:"required,numeric"`
	Product   Product `gorm:"foriegnkey:ProductID;references:ID"`
	//Qty         uint
	Price float64
	//TotalAmount float64
	OrderStatus    string  `gorm:"type:varchar(10);check:order_status IN ('pending','shipped', 'delivered', 'cancelled','failed','return');default:'pending'" json:"order_status" validate:"required,oneof=pending delivered shipped cancelled failed return"`
	PaymentMethod  string  `gorm:"type:varchar(10); check(order_status IN ('COD', 'RazorPay','Wallet')) ;default:'COD'" json:"payment_method" validate:"required"`
	CouponDiscount float64 `gorm:"default:0.00"`
	OfferDiscount  float64 `gorm:"default:0.00"`
	TotalDiscount  float64 `gorm:"default:0.00"`
	PaidAmount     float64 `gorm:"default:0.00"`
	DeliveredDate  string
}
type SalesReportItem struct {
	OrderID        uint
	ProductID      uint
	ProductName    string
	Qty            uint
	Price          float64
	OrderStatus    string
	PaymentMethod  string
	CouponDiscount float64
	OfferDiscount  float64
	TotalDiscount  float64
	PaidAmount     float64
	OrderDate      time.Time
	DeliveredDate  string
}
type Payments struct {
	gorm.Model
	UserID        uint `validate:"required"`
	OrderID       uint `validate:"required"`
	OrderItemID   uint
	TotalAmount   float64 `validate:"required,numeric"`
	TransactionID string
	PaymentDate   string
	PaymentType   string `gorm:"type:varchar(10); check(status IN ('COD', 'RazorPay')) ;default:'COD'" json:"payment_type" validate:"required"`
	PaymentStatus string `gorm:"type:varchar(10); check(status IN ('pending', 'paid', 'refund')) ;default:'pending'" json:"payment_status" validate:"required"`
	Description   string
}

type Wallet struct {
	gorm.Model
	UserID  uint    `gorm:"not null"`
	Balance float64 `gorm:"type:decimal(10,2);default:0.00"`
}

type Wishlist struct {
	gorm.Model
	UserID      uint `gorm:"not null"`
	ProductID   uint `gorm:"not null"`
	ProductName string
}

type Coupon struct {
	gorm.Model
	Code        string  `gorm:"not null" json:"code"`
	Discount    float64 `gorm:"type:decimal(5,2);not null" json:"discount"`
	MinPurchase float64 `gorm:"type:decimal(10,2)" json:"min_purchase"`
}



type WalletTransaction struct {
	gorm.Model
	UserID          uint    `gorm:"not null"`
	Amount          float64 `gorm:"not null"`
	TransactionType string  `gorm:"size:50;not null"`
	Description     string  `gorm:"size:255"`
}

type Invoice struct {
	No             int
	ProductID      uint
	ProductName    string
	Quantity       uint
	MRP            float64
	CouponDiscount float64
	OfferDiscount  float64
	TotalDiscount  float64
	FinalPrice     float64
}

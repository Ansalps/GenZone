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
	City       string `validate:"required"`
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
	FirstName      string `validate:"required"`
	LastName       string `validate:"required"`
	Email          string `gorm:"unique" validate:"required"`
	Password       string `validate:"required"`
	Phone          string `json:"phone" validate:"required,numeric,len=10"`
	ProfilePicture string `gorm:"default:null" json:"profile_picture"`
	Status         string `gorm:"type:varchar(10); check(status IN ('Active', 'Blocked', 'Deleted')) ;default:'Active'" json:"status" validate:"required"`
}

type Category struct {
	ID uint `gorm:"primaryKey" json:"id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	CategoryName string `gorm:"not null;unique" json:"category_name" validate:"required"`
	Description  string `gorm:"not null" json:"category_description" validate:"required"`
	ImageURL     string `gorm:"not null" json:"category_image_url" validate:"required"`

	// Self-referencing relationship
	ParentID *uint      `gorm:"index" json:"parent_id"`
	Parent   *Category  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`

	Products []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}

type Product struct {
	ID uint `gorm:"primaryKey" json:"id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	CategoryID uint     `gorm:"not null;index" json:"category_id"`
	Category   Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`

	ProductName string `gorm:"not null" json:"product_name" validate:"required"`
	Description string `gorm:"not null" json:"product_description" validate:"required"`
	ImageURL    string `gorm:"not null" json:"product_image_url" validate:"required"`

	Price float64 `gorm:"type:decimal(10,2);not null" json:"price"`

	Popular bool `gorm:"not null;default:false" json:"popular"`

	Variants []ProductVariant `gorm:"foreignKey:ProductID" json:"variants,omitempty"`
}

type ProductVariant struct {
	ID uint `gorm:"primaryKey" json:"id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	ProductID uint    `gorm:"not null;index;uniqueIndex:idx_product_size" json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`

	Size  string `gorm:"type:varchar(10);not null;uniqueIndex:idx_product_size" json:"size" validate:"required"`
	Stock uint   `gorm:"not null;default:0" json:"stock"`
}

type Offer struct {
	ID                 uint `gorm:"primarykey"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	ProductID          uint    `gorm:"not null"`
	DiscountPercentage float64 `gorm:"not null"`
	StartAt            time.Time
	EndAt              time.Time
}

// Cart represents a user's shopping cart.
type Cart struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Foreign Key to User model with UNIQUE constraint (One Cart per User)
	UserID uint `gorm:"not null;uniqueIndex;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"user,omitempty"`

	// Relationship to CartItems
	CartItems []CartItem `gorm:"foreignKey:CartID" json:"cart_items,omitempty"`
}

// CartItem represents individual product entries within a user's cart.
type CartItem struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Foreign Key linking to Cart
	CartID uint `gorm:"not null;index;uniqueIndex:idx_cart_product_size;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"cart_id"`
	Cart   Cart `gorm:"foreignKey:CartID" json:"cart,omitempty"`

	// Foreign Key linking to Product
	ProductID uint    `gorm:"not null;index;uniqueIndex:idx_cart_product_size;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`

	Size string `gorm:"type:varchar(32);not null;uniqueIndex:idx_cart_product_size" json:"size" validate:"required"`

	Quantity  uint    `gorm:"not null;default:1" json:"quantity" validate:"required,gt=0"`
	UnitPrice float64 `gorm:"type:decimal(10,2);not null" json:"unit_price" validate:"required,gt=0"`
}

type Order struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID uint `validate:"required"`
	User User `gorm:"foriegnkey:UserID;reference:ID"`
	AddressID   uint
	Address     Address `gorm:"foriegnkey:AddressID;references:ID"`
	TotalAmount float64
	PaymentMethod  string  `gorm:"type:varchar(10); check(order_status IN ('COD', 'RazorPay')) ;default:'COD'" json:"payment_method" validate:"required"`
	OrderStatus    string  `gorm:"type:varchar(10);check:order_status IN ('pending','shipped', 'delivered', 'cancelled','failed');default:'pending'" json:"order_status" validate:"required,oneof=pending delivered shipped cancelled failed"`
	CouponID    uint
	Coupon Coupon `gorm:"foriegnkey:CouponID;references:ID"`
	OfferDiscount float64
	CouponDiscount float64
	TotalDiscountAmount float64 `gorm:"type:decimal(10,2);default:0.00"`
	FinalAmount    float64 `gorm:"type:decimal(10,2);not null"`
}

type OrderItems struct {
	gorm.Model
	OrderID   uint    `validate:"required"`
	Order     Order   `gorm:"foriegnkey:OrderID;references:ID"`
	ProductID uint    `validate:"required,numeric"`
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
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Code        string  `gorm:"not null" json:"code"`
	Discount    float64 `gorm:"type:decimal(5,2);not null" json:"discount"`
	MinPurchase float64 `gorm:"type:decimal(10,2)" json:"min_purchase"`
	StartAt            time.Time
	EndAt              time.Time
	IsActive bool

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

package responsemodels

import "time"

type Category struct {
	//gorm.Model
	ID           uint   ` json:"id"`
	CategoryName string ` gorm:"unique" json:"category_name" validate:"required"`
	Description  string `json:"category_description" validate:"required"`
	ImageUrl     string `json:"category_imageUrl" validate:"required"`
}

type Product struct {
	//gorm.Model
	ID           uint   `json:"id"`
	CategoryID   uint   `json:"category_id" validate:"required"`
	CategoryName string `json:"category_name" validate:"required"`
	//Category     Category `gorm:"foriegnkey:CategoryID;references:ID" json:"category,omitempty"`
	ProductName          string  `json:"product_name" validate:"required"`
	Description          string  `json:"product_description" validate:"required"`
	ImageUrl             string  `json:"product_imageUrl" validate:"required"`
	Price                float64 `gorm:"type:decimal(10,2)" json:"price" validate:"required"`
	Stock                int     `json:"stock"`
	Popular              bool    `gorm:"type:boolean;default:false" json:"popular" validate:"required"`
	Size                 string  `gorm:"type:varchar(10); check(size IN ('Medium', 'Small', 'Large'))" json:"size" validate:"required"`
	HasOffer             bool    `json:"has_offer"`
	OfferDiscountPercent uint    `json:"offer_discount_percent"`
}

type CartItems struct {
	//gorm.Model
	UserID string `validate:"required,numeric" json:"user_id"`
	//User   User   `gorm:"foriegnkey:UserID;references:ID"`
	// CartID      string  `validate:"required,numeric"`
	// Cart        Cart    `gorm:"foriegnkey:CartID;references:ID"`
	ProductID   string `validate:"required,numeric" json:"product_id"`
	ProductName string `json:"product_name" validate:"required"`
	//Product     Product `gorm:"foriegnkey:ProductID;references:ID"`
	TotalAmount float64 `gorm:"type:decimal(10,2);default:0.00" json:"total_amount" validate:"required"`
	Qty         uint    `gorm:"default:0" json:"qty"`
	Price       float64 `gorm:"type:decimal(10,2)" json:"price" validate:"required"`
	Discount    float64 `json:"discount"`
	FinalAmount float64 `json:"final_amount"`
}

type Address struct {
	//gorm.Model
	ID uint `json:"id"`
	// CreatedAt time.Time    `json:"created_at"`
	// UpdatedAt time.Time    `json:"updated_at"`
	// DeletedAt sql.NullTime `json:"deleted_at"`
	UserID string `validate:"required,numeric" json:"user_id"`
	//User       User   `gorm:"foriegnkey:UserID;references:ID"`
	Country    string `validate:"required" json:"country"`
	State      string `validate:"required" json:"state"`
	District   string `validate:"required" json:"district"`
	StreetName string `validate:"required" json:"street_name"`
	PinCode    string `validate:"required,numeric" json:"pin_code"`
	Phone      string `validate:"required,numeric,len=10" json:"phone"`
	Default    bool   `gorm:"default:false" validate:"required" json:"default"`
}

type Order struct {
	//gorm.Model
	ID uint `json:"id"`
	// CreatedAt time.Time
	// UpdatedAt time.Time
	// DeletedAt sql.NullTime `gorm:"index"`
	UserID string `validate:"required,numeric" json:"user_id"`
	//OrderDate   time.Time
	AddressID      uint    `json:"address_id"`
	Address        Address `gorm:"foriegnkey:AddressID;references:ID" json:"address"`
	TotalAmount    float64 `json:"total_amount"`
	PaymentMethod  string  `json:"payment_method"`
	OrderStatus    string  `gorm:"type:varchar(10); check(status IN ('pending', 'delivered', 'cancelled')) ;default:'pending'" json:"order_status" validate:"required"`
	OfferApplied   float64 `json:"offer_applied"`
	CouponCode     string  `json:"coupon_code"`
	DiscountAmount float64 `json:"discount_amount"`
	FinalAmount    float64 `json:"final_amount"`
}
type OrderItems struct {
	//gorm.Model
	ID      uint   `json:"id"`
	OrderID string `validate:"required,numeric" json:"order_id"`
	//Order       Order   `gorm:"foriegnkey:OrderID;references:ID"`
	ProductID string `validate:"required,numeric" json:"product_id"`
	//Product     Product `gorm:"foriegnkey:ProductID;references:ID"`
	ProductName string `json:"product_name"`
	//Qty         uint
	Price       float64 `json:"price"`
	OrderStatus string  `json:"order_staus"`
	//TotalAmount float64
	PaymentMethod  string  `json:"payment_method"`
	CouponDiscount float64 `json:"coupon_discount"`
	OfferDiscount  float64 `json:"offer_discount"`
	TotalDiscount  float64 `json:"total_discount"`
	PaidAmount     float64 `json:"paid_amount"`
}

type User struct {
	//gorm.Model
	ID        uint   ` json:"id"`
	FirstName string `validate:"required" json:"first_name"`
	LastName  string `validate:"required" json:"last_name"`
	Email     string `gorm:"unique" validate:"required" json:"email"`
	Password  string `validate:"required" json:"password"`
	Phone     string `json:"phone" validate:"required,numeric,len=10"`
	Status    string `gorm:"type:varchar(10); check(status IN ('Active', 'Blocked', 'Deleted')) ;default:'Active'" json:"status" validate:"required"`
}

type Wishlist struct {
	ID uint `json:"id"`
	//gorm.Model
	UserID               uint    `gorm:"not null" json:"user_id"`
	ProductID            uint    `gorm:"not null" json:"product_id"`
	ProductName          string  `json:"product_name"`
	CategoryName         string  `json:"category_name"`
	Description          string  `json:"desciption"`
	ImageUrl             string  `json:"image_url"`
	Price                float64 `json:"price"`
	Stock                uint    `json:"stock"`
	Popular              bool    `json:"popular"`
	Size                 string  `json:"size"`
	HasOffer             bool    `json:"has_offer"`
	OfferDiscountPercent uint    `json:"offer_discount_amount"`
}

type Offer struct {
	//gorm.Model
	ID                   uint    `json:"id"`
	ProductID            uint    `gorm:"not null" json:"product_id"`
	DiscountPercentage   uint    `gorm:"not null" json:"discount_percentage"`
	ProductName          string  `json:"product_name"`
	CategoryName         string  `json:"category_name"`
	Description          string  `json:"description"`
	ImageUrl             string  `json:"image_url"`
	Price                float64 `json:"price"`
	Stock                uint    `json:"stock"`
	Popular              bool    `json:"popular"`
	Size                 string  `json:"size"`
	HasOffer             bool    `json:"has_offer"`
	OfferDiscountPercent uint    `json:"offer_discount_percent"`
}
type BestSelling struct {
	Count        int
	CategoryName string `json:"category_name"`
}

type Coupon struct {
	//gorm.Model
	ID          uint    `json:"id"`
	Code        string  `gorm:"not null" json:"code"`
	Discount    float64 `gorm:"type:decimal(5,2);not null" json:"discount"`
	MinPurchase float64 `gorm:"type:decimal(10,2)" json:"min_purchase"`
}
type SalesReportItem struct {
	OrderID        uint      `json:"order_id"`
	ProductID      string    `json:"product_id"`
	ProductName    string    `json:"product_name"`
	Qty            uint      `json:"qty"`
	Price          float64   `json:"price"`
	OrderStatus    string    `json:"order_status"`
	PaymentMethod  string    `json:"payment_method"`
	CouponDiscount float64   `json:"coupon_discount"`
	OfferDiscount  float64   `json:"offer_discount"`
	TotalDiscount  float64   `json:"total_discount"`
	PaidAmount     float64   `json:"paid_amount"`
	OrderDate      time.Time `json:"order_date"`
	DeliveredDate  string    `json:"delivered_date"`
}
type Wallet struct {
	//gorm.Model
	ID      uint    `json:"id"`
	UserID  uint    `gorm:"not null" json:"user_id"`
	Balance float64 `gorm:"type:decimal(10,2);default:0.00" json:"balance"`
}
type WalletTransaction struct {
	//gorm.Model
	ID              uint    `json:"id"`
	UserID          uint    `gorm:"not null" json:"user_id"`
	Amount          float64 `gorm:"not null" json:"amount"`
	TransactionType string  `gorm:"size:50;not null" json:"transaction_type"`
	Description     string  `gorm:"size:255" json:"description"`
}

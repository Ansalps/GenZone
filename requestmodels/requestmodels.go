package requestmodemodels

type AdminLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type Category struct {
	CategoryName string `json:"category_name" validate:"required,no_leading_trailing_spaces,no_repeating_spaces,max=50"`
	Description  string `json:"category_description" validate:"required,no_leading_trailing_spaces,no_repeating_spaces,max=100"`
}

type Product struct {
	
	CategoryName string `json:"category_name" validate:"required"`
	ProductName string  `json:"product_name" validate:"required,no_leading_trailing_spaces,no_repeating_spaces,max=50"`
	Description string  `json:"product_description" validate:"required,no_leading_trailing_spaces,no_repeating_spaces,max=100"`
	ImageUrl    string  `json:"product_image_url" validate:"required,excludesall= "`
	Price       float64 ` json:"price" validate:"required,gt=0"`
	Stock       uint    `json:"stock" validate:"required"`
	Popular     bool    `json:"popular"`
	Size        string  ` json:"size" validate:"required"`
	DiscountPercentage float64 `json:"discount_percentage"`
	
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
	Country    string `json:"country" validate:"required,no_leading_trailing_spaces,no_repeating_spaces,max=50,alpha"`
	State      string `json:"state" validate:"required,no_leading_trailing_spaces,no_repeating_spaces,max=50,alpha"`
	District   string `json:"district" validate:"required,no_leading_trailing_spaces,no_repeating_spaces,max=50,alpha"`
	StreetName string `json:"street_name" validate:"required,no_leading_trailing_spaces,no_repeating_spaces,max=50,alpha"`
	PinCode    string `json:"pin_code" validate:"required,numeric"`
	Phone      string `json:"phone" validate:"required,numeric,len=10"`
	Default    bool   `json:"Default" `
}

type CartAdd struct {
	ProductID string `json:"product_id" validate:"required,numeric"`
}

type OrderAdd struct {
	AddressID  uint   `json:"address_id" validate:"required"`
	CouponCode string `json:"coupon_code"`
}

type ProfileEdit struct {
	FirstName string `validate:"required,excludesall= " json:"name"`
	LastName  string `validate:"required,nameOrInitials" json:"last_name"`
	//Email           string `gorm:"unique" validate:"required,email" json:"email"`
	//Password        string `validate:"required" json:"password"`
	//ConfirmPassword string `validate:"required" json:"confirmpassword"`
	Phone string `json:"phone" validate:"required,numeric,len=10"`
}

type PasswordChange struct {
	Password        string `validate:"required,min=8,password" json:"password"`
	ConfirmPassword string `validate:"required" json:"confirmpassword"`
}

type CancelOrder struct {
	OrderStatus string `json:"order_status" validate:"required"`
}

// type RazorPayOrder struct {
// 	UserID    uint `validate:"required"`
// 	AddressID uint `json:"address_id" validate:"required"`
// }
type WishlistAdd struct {
	ProductID uint `json:"product_id" validate:"required"`
}

type CouponAdd struct {
	Code        string  `validate:"required" json:"code"`
	Discount    float64 `validate:"required" json:"discount"`
	MinPurchase float64 `validate:"required" json:"min_purchase"`
}
type Offer struct {
	ProductID          uint `validate:"required" json:"product_id"`
	DiscountPercentage float64 `validate:"required" json:"discount_percentage"`
}
type CouponCheckout struct {
	CouponCode string `json:"coupon_code"`
}

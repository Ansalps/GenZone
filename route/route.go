package route

import (
	"github.com/Ansalps/GeZOne/admin"
	"github.com/Ansalps/GeZOne/middleware"
	"github.com/Ansalps/GeZOne/user"
	"github.com/gin-gonic/gin"
)

func RegisterUrls(router *gin.Engine) {

	//users
	adminGroup := router.Group("admin/")
	adminGroup.GET("listusers", middleware.AuthMiddleware("admin"), admin.ListUsers)
	adminGroup.PUT("listusers/blockuser", middleware.AuthMiddleware("admin"), admin.BlockUser)
	adminGroup.PUT("listusers/unblockuser", middleware.AuthMiddleware("admin"), admin.UnblockUser)
	//category

	adminGroup.POST("/", admin.Login)
	//adminCategory := router.Group("admin/Category/")
	adminGroup.GET("category", middleware.AuthMiddleware("admin"), admin.Category)
	adminGroup.POST("category", middleware.AuthMiddleware("admin"), admin.CategoryAdd)
	adminGroup.PUT("category/:id", middleware.AuthMiddleware("admin"), admin.CategoryEdit)
	adminGroup.DELETE("category/:id", middleware.AuthMiddleware("admin"), admin.CategoryDelete)

	//products
	adminGroup.GET("product", middleware.AuthMiddleware("admin"), admin.Product)
	adminGroup.POST("product", middleware.AuthMiddleware("admin"), admin.ProductAdd)
	adminGroup.PUT("product/:id", middleware.AuthMiddleware("admin"), admin.ProductEdit)
	adminGroup.DELETE("product/:id", middleware.AuthMiddleware("admin"), admin.ProductDelete)

	//order
	adminGroup.GET("orderlist", middleware.AuthMiddleware("admin"), admin.OrderList)
	adminGroup.GET("orderlist/items/:order_id", middleware.AuthMiddleware("admin"), admin.OrderItemsList)
	adminGroup.PUT("order/changestatus/:id", middleware.AuthMiddleware("admin"), admin.ChangeOrderStatus)

	//coupon
	adminGroup.GET("coupon", middleware.AuthMiddleware("admin"), admin.CouponList)
	adminGroup.POST("coupon", middleware.AuthMiddleware("admin"), admin.CouponAdd)
	adminGroup.DELETE("coupon/:id", middleware.AuthMiddleware("admin"), admin.CouponRemove)

	//productoffer
	adminGroup.GET("offer", middleware.AuthMiddleware("admin"), admin.OfferList)
	adminGroup.POST("offer", middleware.AuthMiddleware("admin"), admin.OfferAdd)
	adminGroup.DELETE("offer/:id", middleware.AuthMiddleware("admin"), admin.OfferRemove)

	//salesreport
	adminGroup.POST("salesreport", middleware.AuthMiddleware("admin"), admin.GenerateSalesReport)
	adminGroup.GET("salesreport", middleware.AuthMiddleware("admin"), admin.FilterSalesReport)
	adminGroup.GET("salesreportdownload", middleware.AuthMiddleware("admin"), admin.FilterSalesReportPdfExcel)

	//best selling
	adminGroup.GET("bestselling", middleware.AuthMiddleware("admin"), admin.BestSelling)
	adminGroup.GET("invoice/:order_id", middleware.AuthMiddleware("admin"), admin.GenerateInvoice)

	//user

	router.POST("/signup/", user.UserSignUp)
	router.POST("/login/", user.UserLogin)
	router.POST("/signup/verifyotp/:email", user.VerifyOTPHandler)
	router.POST("/signup/resendotp/:email", user.ResendOtp)
	router.GET("/", middleware.AuthMiddleware("user"), user.ListProducts)
	router.GET("/auth/google/login", user.HandleGoogleLogin)
	router.GET("/auth/google/callback", user.HandleGoogleCallback)

	router.GET("searchproduct", middleware.AuthMiddleware("user"), user.SearchProduct)

	router.GET("profile", middleware.AuthMiddleware("user"), user.Profile)
	router.PUT("profile", middleware.AuthMiddleware("user"), user.ProfileEdit)
	router.GET("profile/userorders", middleware.AuthMiddleware("user"), user.OrderList)
	router.GET("profile/userorders/items/:order_id", middleware.AuthMiddleware("user"), user.OrderItemsList)
	router.PUT("profile/userorders/cancelorder/:order_id", middleware.AuthMiddleware("user"), user.CancelOrder)
	router.PUT("profile/userorders/cancelsingleorderitem/:orderitem_id", middleware.AuthMiddleware("user"), user.CancelSingleOrderItem)
	router.PUT("profile/userorders/returnsingleorderitem/:orderitem_id", middleware.AuthMiddleware("user"), user.ReturnSingleOrderItem)
	//wallet listing
	router.GET("proflie/wallet", middleware.AuthMiddleware("user"), user.WalletListing)
	router.GET("profile/wallettransaction", middleware.AuthMiddleware("user"), user.WalletTransactionListing)
	router.GET("profile/wishlist", middleware.AuthMiddleware("user"), user.Wishlist)
	router.POST("profile/wishlist", middleware.AuthMiddleware("user"), user.WishlistAdd)
	router.DELETE("profile/wishlist", middleware.AuthMiddleware("user"), user.WishlistRemove)
	router.PUT("profile/changepassword", middleware.AuthMiddleware("user"), user.PasswordChange)

	//router.GET("address", helper.AuthMiddleware("user"), user.Address)
	router.GET("profile/useraddress", middleware.AuthMiddleware("user"), user.AddressList)
	router.POST("profile/useraddress", middleware.AuthMiddleware("user"), user.AddressAdd)
	router.PUT("profile/useraddress/:address_id", middleware.AuthMiddleware("user"), user.AddressEdit)
	router.DELETE("profile/useraddress/:address_id", middleware.AuthMiddleware("user"), user.AddressDelete)

	router.GET("cart", middleware.AuthMiddleware("user"), user.Cart)
	router.POST("cart", middleware.AuthMiddleware("user"), user.CartAdd)
	router.DELETE("cart", middleware.AuthMiddleware("user"), user.CartRemove)

	router.GET("checkout", middleware.AuthMiddleware("user"), user.CheckOut)
	//router.GET("checkout/:user_id/address", helper.AuthMiddleware("user"), user.CheckOutAddress)
	router.PUT("checkout/address/:address_id", middleware.AuthMiddleware("user"), user.CheckOutAddressEdit)
	router.POST("checkout/order", middleware.AuthMiddleware("user"), user.Order)
	router.POST("checkout/razorpay", middleware.AuthMiddleware("user"), user.CreateOrder)
	router.POST("checkout/razorpay/paymentverification", middleware.AuthMiddleware("user"), user.PaymentWebhook)
	router.POST("checkout/wallet", middleware.AuthMiddleware("user"), user.WalletOrder)

}

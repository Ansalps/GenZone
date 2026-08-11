package route

import (
	"github.com/Ansalps/GeZOne/admin"
	"github.com/Ansalps/GeZOne/middleware"
	"github.com/Ansalps/GeZOne/public"
	"github.com/Ansalps/GeZOne/user"
	"github.com/gin-gonic/gin"
)

func RegisterUrls(router *gin.Engine) {

	adminGroup := router.Group("admin")
	
	//admin login/logout
	adminGroup.POST("/login", admin.Login)
	adminGroup.POST("/logout", admin.Logout)

	//admin category management
	adminGroup.GET("category", middleware.AuthMiddleware("admin"), admin.ReadCategory)
	adminGroup.GET("category/:id",middleware.AuthMiddleware("admin"),admin.ReadCategoryById)
	adminGroup.POST("category", middleware.AuthMiddleware("admin"), admin.AddCategory)
	adminGroup.PUT("category/:id", middleware.AuthMiddleware("admin"), admin.EditCategory)
	adminGroup.DELETE("category/:id", middleware.AuthMiddleware("admin"), admin.CategoryDelete)

	//admin products management
	adminGroup.GET("product", middleware.AuthMiddleware("admin"), admin.ReadProducts)
	adminGroup.GET("product/:id", middleware.AuthMiddleware("admin"), admin.ReadProductById)
	adminGroup.POST("product", middleware.AuthMiddleware("admin"), admin.AddProduct)
	adminGroup.PUT("product/:id", middleware.AuthMiddleware("admin"), admin.EditProduct)
	adminGroup.DELETE("product/:id", middleware.AuthMiddleware("admin"), admin.ProductDelete)

	//admin user management
	adminGroup.GET("listusers", middleware.AuthMiddleware("admin"), admin.ListUsers)
	adminGroup.PUT("listusers/blockuser", middleware.AuthMiddleware("admin"), admin.BlockUser)
	adminGroup.PUT("listusers/unblockuser", middleware.AuthMiddleware("admin"), admin.UnblockUser)
	
	
	//order management
	adminGroup.GET("orderlist", middleware.AuthMiddleware("admin"), admin.OrderList)
	adminGroup.GET("orderlist/items/:order_id", middleware.AuthMiddleware("admin"), admin.OrderItemsList)
	adminGroup.PUT("order/changestatus/:id", middleware.AuthMiddleware("admin"), admin.ChangeOrderStatus)

	//coupon management
	adminGroup.GET("coupon", middleware.AuthMiddleware("admin"), admin.CouponList)
	adminGroup.POST("coupon", middleware.AuthMiddleware("admin"), admin.CouponAdd)
	adminGroup.DELETE("coupon/:id", middleware.AuthMiddleware("admin"), admin.CouponRemove)

	//productoffer management
	adminGroup.GET("offer", middleware.AuthMiddleware("admin"), admin.OfferList)
	adminGroup.POST("offer", middleware.AuthMiddleware("admin"), admin.OfferAdd)
	adminGroup.DELETE("offer/:id", middleware.AuthMiddleware("admin"), admin.OfferRemove)

	//salesreport generation
	adminGroup.POST("salesreport", middleware.AuthMiddleware("admin"), admin.GenerateSalesReport)
	adminGroup.GET("salesreport", middleware.AuthMiddleware("admin"), admin.FilterSalesReport)
	adminGroup.GET("salesreportdownload", middleware.AuthMiddleware("admin"), admin.FilterSalesReportPdfExcel)

	//best selling, invoice generation
	adminGroup.GET("bestselling", middleware.AuthMiddleware("admin"), admin.BestSelling)
	adminGroup.GET("invoice/:order_id", middleware.AuthMiddleware("admin"), admin.GenerateInvoice)

	//public
	//router.GET("", public.ListProducts)
	router.GET("",public.ReadCategory)
	router.POST("signup", user.UserSignUp)
	router.POST("signup/verifyotp/:email", user.VerifyOTPHandler)
	router.POST("signup/resendotp/:email", user.ResendOtp)
	router.POST("login", user.UserLogin)
	router.GET("auth/google/login", user.HandleGoogleLogin)
	router.GET("auth/google/callback", user.HandleGoogleCallback)
	router.GET("searchproduct",  user.SearchProduct)

	//user
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

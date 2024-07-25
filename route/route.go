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
	adminGroup.PUT("order/changestatus/:id", middleware.AuthMiddleware("admin"), admin.ChangeOrderStatus)

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
	router.GET("profile/useraddress", middleware.AuthMiddleware("user"), user.AddressList)
	router.GET("profile/userorders", middleware.AuthMiddleware("user"), user.OrderList)
	router.GET("profile/userorders/items/:order_id", middleware.AuthMiddleware("user"), user.OrderItemsList)
	router.PUT("profile/userorders/cancelorder/:order_id", middleware.AuthMiddleware("user"), user.CancelOrder)
	router.PUT("profile/userorders/cancelsingleorderitem/:orderitem_id", middleware.AuthMiddleware("user"), user.CancelSingleOrderItem)
	router.PUT("profile/changepassword", middleware.AuthMiddleware("user"), user.PasswordChange)

	//router.GET("address", helper.AuthMiddleware("user"), user.Address)
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

}

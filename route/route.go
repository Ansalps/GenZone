package route

import (
	"github.com/Ansalps/GeZOne/admin"
	"github.com/Ansalps/GeZOne/helper"
	"github.com/Ansalps/GeZOne/user"
	"github.com/gin-gonic/gin"
)

func RegisterUrls(router *gin.Engine) {

	//users
	adminGroup := router.Group("admin/")
	adminGroup.GET("listusers", helper.AuthMiddleware("admin"), admin.ListUsers)
	adminGroup.PUT("listusers/blockuser/:id", helper.AuthMiddleware("admin"), admin.BlockUser)
	adminGroup.PUT("listusers/unblockuser/:id", helper.AuthMiddleware("admin"), admin.UnblockUser)
	//category

	adminGroup.POST("/", admin.Login)
	//adminCategory := router.Group("admin/Category/")
	adminGroup.GET("category", helper.AuthMiddleware("admin"), admin.Category)
	adminGroup.POST("category", helper.AuthMiddleware("admin"), admin.CategoryAdd)
	adminGroup.PUT("category/:id", helper.AuthMiddleware("admin"), admin.CategoryEdit)
	adminGroup.DELETE("category/:id", helper.AuthMiddleware("admin"), admin.CategoryDelete)

	//products
	adminGroup.GET("product", helper.AuthMiddleware("admin"), admin.Product)
	adminGroup.POST("product", helper.AuthMiddleware("admin"), admin.ProductAdd)
	adminGroup.PUT("product/:id", helper.AuthMiddleware("admin"), admin.ProductEdit)
	adminGroup.DELETE("product/:id", helper.AuthMiddleware("admin"), admin.ProductDelete)

	//order
	adminGroup.GET("orderlist", helper.AuthMiddleware("admin"), admin.OrderList)
	adminGroup.PUT("order/changestatus/:id", helper.AuthMiddleware("admin"), admin.ChangeOrderStatus)

	//user

	router.POST("/signup/", user.UserSignUp)
	router.POST("/login/", user.UserLogin)
	router.POST("/signup/verifyotp/:email", user.VerifyOTPHandler)
	router.POST("/signup/resendotp/:email", user.ResendOtp)
	router.GET("/", helper.AuthMiddleware("user"), user.ListProducts)
	router.GET("/auth/google/login", user.HandleGoogleLogin)
	router.GET("/auth/google/callback", user.HandleGoogleCallback)

	router.GET("searchproduct", helper.AuthMiddleware("user"), user.SearchProduct)

	router.GET("profile", helper.AuthMiddleware("user"), user.Profile)
	router.PUT("profile", helper.AuthMiddleware("user"), user.ProfileEdit)
	router.GET("profile/useraddress", helper.AuthMiddleware("user"), user.AddressList)
	router.GET("profile/userorders", helper.AuthMiddleware("user"), user.OrderList)
	router.PUT("profile/userorders/cancelorder/:order_id", helper.AuthMiddleware("user"), user.CancelOrder)
	router.PUT("profile/changepassword", helper.AuthMiddleware("user"), user.PasswordChange)

	//router.GET("address", helper.AuthMiddleware("user"), user.Address)
	router.POST("profile/useraddress", helper.AuthMiddleware("user"), user.AddressAdd)
	router.PUT("profile/useraddress/:address_id", helper.AuthMiddleware("user"), user.AddressEdit)
	router.DELETE("profile/useraddress/:address_id", helper.AuthMiddleware("user"), user.AddressDelete)

	router.GET("cart", helper.AuthMiddleware("user"), user.Cart)
	router.POST("cart", helper.AuthMiddleware("user"), user.CartAdd)
	router.DELETE("cart", helper.AuthMiddleware("user"), user.CartRemove)

	router.GET("checkout", helper.AuthMiddleware("user"), user.CheckOut)
	//router.GET("checkout/:user_id/address", helper.AuthMiddleware("user"), user.CheckOutAddress)
	router.PUT("checkout/address/:address_id", helper.AuthMiddleware("user"), user.CheckOutAddressEdit)
	router.POST("checkout/order", helper.AuthMiddleware("user"), user.Order)

}

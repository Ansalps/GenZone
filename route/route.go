package route

import (
	"github.com/Ansalps/GeZOne/admin"
	"github.com/gin-gonic/gin"
)

func RegisterUrls(router *gin.Engine) {
	//category
	adminGroup := router.Group("admin/")
	adminGroup.POST("/", admin.Login)
	//adminCategory := router.Group("admin/Category/")
	adminGroup.GET("category/", admin.Category)
	adminGroup.GET("category/add/", admin.CategoryAdd)
	adminGroup.PUT("category/edit/:id", admin.CategoryEdit)
	adminGroup.DELETE("category/delete/:id", admin.CategoryDelete)

	//products
	adminGroup.GET("product/", admin.Product)
	adminGroup.GET("product/add/", admin.ProductAdd)
	adminGroup.PUT("product/edit/:id", admin.ProductEdit)
	adminGroup.DELETE("product/delete/:id", admin.ProductDelete)
}

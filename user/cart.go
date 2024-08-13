package user

import (
	"fmt"
	"net/http"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/helper"
	"github.com/Ansalps/GeZOne/middleware"
	"github.com/Ansalps/GeZOne/models"
	"github.com/Ansalps/GeZOne/responsemodels"
	"github.com/gin-gonic/gin"
)

func Cart(c *gin.Context) {
	//UserID := c.Param("user_id")
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Claims not found"})
		return
	}

	customClaims, ok := claims.(*middleware.CustomClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
		return
	}

	userID := customClaims.ID
	fmt.Println("print user id : ", userID)
	var cart []responsemodels.CartItems
	//tx := database.DB.Where("user_id = ?", UserID).Find(&cart)
	tx := database.DB.Raw("SELECT cart_items.user_id,cart_items.product_id,products.product_name,cart_items.total_amount,cart_items.qty,cart_items.price,cart_items.discount,cart_items.final_amount FROM cart_items join products on cart_items.product_id=products.id where user_id = ? AND cart_items.deleted_at IS NULL AND cart_items.qty != 0", userID).Scan(&cart)
	if tx.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "failed to retrieve data from the database, or the data doesn't exists",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "successfully retrieved user informations",
		"data": gin.H{
			"Cart Items": cart,
		},
	})
}

func CartAdd(c *gin.Context) {
	//UserID := c.Param("user_id")
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Claims not found"})
		return
	}

	customClaims, ok := claims.(*middleware.CustomClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
		return
	}

	userID := customClaims.ID
	fmt.Println("print user id : ", userID)
	var Cart models.CartAdd
	err := c.BindJSON(&Cart)
	response := gin.H{
		"status":  false,
		"message": "failed to bind request",
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, response)
		return
	}
	//validate the content of the JSON
	if err := helper.Validate(Cart); err != nil {
		fmt.Println("", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"message":    err.Error(),
			"error_code": http.StatusBadRequest,
		})
		return
	}
	var count1 int64
	database.DB.Raw("SELECT COUNT(*) FROM products where id=? AND deleted_at IS NULL", Cart.ProductID).Scan(&count1)
	if count1 == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "product id doesn't exists",
		})
		return
	}
	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM cart_items WHERE user_id=? and product_id=? and deleted_at IS NULL`, userID, Cart.ProductID).Scan(&count)
	if count != 0 {
		var price float64
		database.DB.Model(&models.Product{}).Where("id = ?", Cart.ProductID).Pluck("price", &price)
		var hasoffer bool
		database.DB.Model(&models.Product{}).Where("id = ?", Cart.ProductID).Pluck("has_offer", &hasoffer)
		var finalamount float64
		finalamount = price
		var discount float64
		if hasoffer {
			var discountpercentage uint
			database.DB.Model(&models.Offer{}).Where("product_id = ?", Cart.ProductID).Pluck("discount_percentage", &discountpercentage)
			discount = price * float64(discountpercentage) / 100
			//discount = math.Round(discount*100) / 100
			fmt.Println("disount----", discount)
			database.DB.Raw(`UPDATE cart_items SET dicount = ? WHERE product_id = ?`, discount, Cart.ProductID)
			finalamount = price - discount
			fmt.Println("price---inside", price)
		}
		fmt.Println("price---outside", finalamount)
		var totalamount float64
		database.DB.Model(&models.CartItems{}).Where("user_id = ? and product_id = ?", userID, Cart.ProductID).Pluck("total_amount", &totalamount)
		fmt.Println("total amount:", totalamount)
		totalamount = totalamount + price
		var FinalAmount1 float64
		database.DB.Model(&models.CartItems{}).Where("user_id = ? and product_id = ?", userID, Cart.ProductID).Pluck("final_amount", &FinalAmount1)
		FinalAmount1 = FinalAmount1 + finalamount
		var Discount1 float64
		database.DB.Model(&models.CartItems{}).Where("user_id = ? and product_id = ?", userID, Cart.ProductID).Pluck("discount", &Discount1)
		Discount1 = Discount1 + discount
		fmt.Println("total amount:", totalamount)
		var quantity uint
		database.DB.Model(&models.CartItems{}).Where("user_id = ? and product_id = ?", userID, Cart.ProductID).Pluck("qty", &quantity)
		fmt.Println("quantity:", quantity)
		var stock uint
		database.DB.Model(&models.Product{}).Where("id = ?", Cart.ProductID).Pluck("stock", &stock)
		fmt.Println("--stock", stock)
		if quantity >= 7 {
			c.JSON(http.StatusOK, gin.H{"status": true, "message": "Exceeded maximum quantity for a product"})
			return
		}
		if quantity >= stock {
			c.JSON(http.StatusOK, gin.H{"status": true, "message": "product out of stock"})
			return
		}
		quantity = quantity + 1
		fmt.Println("quantity:", quantity)
		cart := models.CartItems{
			//UserID:      UserID,
			//ProductID:   Cart.ProductID,
			TotalAmount: totalamount,
			Qty:         quantity,
			Price:       price,
			Discount:    Discount1,
			FinalAmount: FinalAmount1,
		}
		database.DB.Model(&models.CartItems{}).Where("user_id = ? and product_id = ?", userID, Cart.ProductID).Updates(&cart)

		// database.DB.Where("user_id = ?", UserID).Order("created_at DESC").First(&Cart)
		// database.DB.Model(&models.CartItems{}).Where("id = ?", UserID).Update("qty", quantity)
		c.JSON(http.StatusOK, gin.H{"status": true, "message": "product added to cart successfully"})
		return
	}
	var stock uint
	database.DB.Model(&models.Product{}).Where("id = ?", Cart.ProductID).Pluck("stock", &stock)
	fmt.Println("--stock", stock)
	if stock == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "product out of stock",
		})
		return
	}
	var price float64
	database.DB.Model(&models.Product{}).Where("id = ?", Cart.ProductID).Pluck("price", &price)
	var hasoffer bool
	database.DB.Model(&models.Product{}).Where("id = ?", Cart.ProductID).Pluck("has_offer", &hasoffer)
	var discount float64
	var finalamount float64
	finalamount = price
	if hasoffer {
		fmt.Println("is it entering in has offer ------")
		var discountpercentage uint
		database.DB.Model(&models.Offer{}).Where("product_id = ?", Cart.ProductID).Pluck("discount_percentage", &discountpercentage)
		discount = price * float64(discountpercentage) / 100
		finalamount = price - (price * float64(discountpercentage) / 100)
		fmt.Println("price---", finalamount)
	}
	fmt.Println("final amount here --", finalamount)
	cart := models.CartItems{
		UserID:      userID,
		ProductID:   Cart.ProductID,
		TotalAmount: price,
		Qty:         1,
		Price:       price,
		Discount:    discount,
		FinalAmount: finalamount,
	}
	database.DB.Create(&cart)
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "product added to cart successfully"})

}

func CartRemove(c *gin.Context) {
	//UserID := c.Param("user_id")
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Claims not found"})
		return
	}

	customClaims, ok := claims.(*middleware.CustomClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
		return
	}

	userID := customClaims.ID
	fmt.Println("print user id : ", userID)
	var Cart models.CartAdd
	err := c.BindJSON(&Cart)
	response := gin.H{
		"status":  false,
		"message": "failed to bind request",
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, response)
		return
	}
	//validate the content of the JSON
	if err := helper.Validate(Cart); err != nil {
		fmt.Println("", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"message":    err.Error(),
			"error_code": http.StatusBadRequest,
		})
		return
	}
	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM cart_items WHERE user_id=? AND product_id=? and deleted_at IS NULL`, userID, Cart.ProductID).Scan(&count)
	if count != 0 {
		var quantity uint
		database.DB.Model(&models.CartItems{}).Where("user_id = ? AND product_id = ?", userID, Cart.ProductID).Pluck("qty", &quantity)
		if quantity == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"messsage": "Product Item doesn't exist in Cart",
			})
			return
		}
		fmt.Println("quantity:", quantity)
		quantity = quantity - 1
		fmt.Println("quantity:", quantity)
		database.DB.Model(&models.CartItems{}).Where("user_id = ? AND product_id = ?", userID, Cart.ProductID).Update("qty", quantity)

		var price float64
		database.DB.Model(&models.CartItems{}).Where("product_id = ?", Cart.ProductID).Pluck("price", &price)
		fmt.Println("product price:", price)
		var totalamount float64
		database.DB.Model(&models.CartItems{}).Where("user_id = ? AND product_id = ?", userID, Cart.ProductID).
			Pluck("total_amount", &totalamount)
		fmt.Println("t a-", totalamount)
		totalamount = totalamount - price
		fmt.Println("t a--", totalamount)
		database.DB.Model(&models.CartItems{}).Where("user_id = ? AND product_id = ?", userID, Cart.ProductID).Order("total_amount DESC").Update("total_amount", totalamount)
		var hasoffer bool
		database.DB.Model(&models.Product{}).Where("id = ?", Cart.ProductID).Pluck("has_offer", &hasoffer)
		fmt.Println("has offer==", hasoffer)
		var discount float64
		var finalamount float64
		if hasoffer {
			fmt.Println("is it entering in has offer ------")
			var discountpercentage uint
			database.DB.Model(&models.Offer{}).Where("product_id = ?", Cart.ProductID).Pluck("discount_percentage", &discountpercentage)
			discount = price * float64(discountpercentage) / 100
			finalamount = price - (price * float64(discountpercentage) / 100)
			fmt.Println("price---", finalamount)
		}
		fmt.Println("price---outside", finalamount)
		var FinalAmount1 float64
		database.DB.Model(&models.CartItems{}).Where("user_id = ? and product_id = ?", userID, Cart.ProductID).Pluck("final_amount", &FinalAmount1)
		FinalAmount1 = FinalAmount1 - finalamount
		database.DB.Model(&models.CartItems{}).Where("user_id = ? AND product_id = ?", userID, Cart.ProductID).Update("final_amount", FinalAmount1)
		var Discount1 float64
		database.DB.Model(&models.CartItems{}).Where("user_id = ? and product_id = ?", userID, Cart.ProductID).Pluck("discount", &Discount1)
		Discount1 = Discount1 - discount
		database.DB.Model(&models.CartItems{}).Where("user_id = ? AND product_id = ?", userID, Cart.ProductID).Update("discount", Discount1)

		c.JSON(http.StatusOK, gin.H{"status": true, "message": "product removed from cart successfully"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "product does not exist in cart"})
}

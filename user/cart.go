package user

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/helper"
	"github.com/Ansalps/GeZOne/middleware"
	"github.com/Ansalps/GeZOne/models"
	requestmodemodels "github.com/Ansalps/GeZOne/requestmodels"
	"github.com/Ansalps/GeZOne/responsemodels"
	"github.com/gin-gonic/gin"
)

func Cart(c *gin.Context) {
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

	// Join carts, cart_items, and products to resolve user_id correctly
	query := `
        SELECT 
            carts.user_id,
            cart_items.product_id,
            products.product_name,
            cart_items.quantity,
            cart_items.unit_price,
            (cart_items.quantity * cart_items.unit_price) AS total_amount
        FROM carts
        JOIN cart_items ON carts.id = cart_items.cart_id
        JOIN products ON cart_items.product_id = products.id
        WHERE carts.user_id = ? 
          AND cart_items.deleted_at IS NULL 
          AND carts.deleted_at IS NULL 
          AND cart_items.quantity > 0
    `
	tx := database.DB.Raw(query, userID).Scan(&cart)

	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "failed to retrieve cart data from the database, or the data doesn't exists",
		})
		return
	}
	fmt.Println("cart_items", cart)
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "successfully retrieved user informations",
		"data": gin.H{
			"cart_items": cart,
		},
	})
}

func CartAdd(c *gin.Context) {
	// 1. Authenticate user from claims
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

	// 2. Bind and validate request body
	var req requestmodemodels.CartAdd
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "failed to bind request",
		})
		return
	}

	if err := helper.Validate(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"message":    err.Error(),
			"error_code": http.StatusBadRequest,
		})
		return
	}

	// 3. Verify product existence
	var product models.Product
	if err := database.DB.Where("id = ?", req.ProductID).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "product id doesn't exist",
		})
		return
	}

	// 4. Check available stock across all variants
	var stock uint
	database.DB.Raw("SELECT COALESCE(SUM(stock), 0) FROM product_variants WHERE product_id = ?", req.ProductID).Scan(&stock)
	if stock == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "product out of stock",
		})
		return
	}

	// 5. Get or Create User's Cart
	var cart models.Cart
	if err := database.DB.FirstOrCreate(&cart, models.Cart{UserID: userID}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "failed to retrieve or create cart",
		})
		return
	}

	// 6. Check for an active offer and calculate unit price
	unitPrice := product.Price
	var offer models.Offer
	now := time.Now()

	err := database.DB.Where("product_id = ? AND start_at <= ? AND end_at >= ?", product.ID, now, now).
		First(&offer).Error

	if err == nil && offer.DiscountPercentage > 0 {
		discount := unitPrice * (offer.DiscountPercentage / 100.0)
		unitPrice -= discount
	}

	// 7. Check if the item is already in the user's cart
	var cartItem models.CartItem
	err = database.DB.Where("cart_id = ? AND product_id = ?", cart.ID, req.ProductID).First(&cartItem).Error

	if err == nil {
		// Item exists -> Check quantity limits
		if cartItem.Quantity >= 7 {
			c.JSON(http.StatusOK, gin.H{
				"status":  false,
				"message": "Exceeded maximum quantity for a product",
			})
			return
		}

		if cartItem.Quantity >= stock {
			c.JSON(http.StatusOK, gin.H{
				"status":  false,
				"message": "product out of stock",
			})
			return
		}

		// Update quantity and unit price
		cartItem.Quantity += 1
		cartItem.UnitPrice = unitPrice

		if err := database.DB.Save(&cartItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "failed to update cart item",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "product quantity updated in cart successfully",
		})
		return
	}

	// 8. Item does not exist -> Insert new record
	newCartItem := models.CartItem{
		CartID:    cart.ID,
		ProductID: product.ID,
		Quantity:  1,
		UnitPrice: unitPrice,
	}

	if err := database.DB.Create(&newCartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "failed to add product to cart",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "product added to cart successfully",
	})
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
	var Cart requestmodemodels.CartAdd
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
		database.DB.Model(&models.CartItem{}).Where("user_id = ? AND product_id = ?", userID, Cart.ProductID).Pluck("qty", &quantity)
		if quantity == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"messsage": "Product Item doesn't exist in Cart",
			})
			return
		}
		fmt.Println("quantity:", quantity)
		quantity = quantity - 1
		fmt.Println("quantity:", quantity)
		database.DB.Model(&models.CartItem{}).Where("user_id = ? AND product_id = ?", userID, Cart.ProductID).Update("qty", quantity)

		var price float64
		database.DB.Model(&models.CartItem{}).Where("product_id = ?", Cart.ProductID).Pluck("price", &price)
		fmt.Println("product price:", price)
		var totalamount float64
		database.DB.Model(&models.CartItem{}).Where("user_id = ? AND product_id = ?", userID, Cart.ProductID).
			Pluck("total_amount", &totalamount)
		fmt.Println("t a-", totalamount)
		totalamount = totalamount - price
		fmt.Println("t a--", totalamount)
		database.DB.Model(&models.CartItem{}).Where("user_id = ? AND product_id = ?", userID, Cart.ProductID).Order("total_amount DESC").Update("total_amount", totalamount)
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
		database.DB.Model(&models.CartItem{}).Where("user_id = ? and product_id = ?", userID, Cart.ProductID).Pluck("final_amount", &FinalAmount1)
		FinalAmount1 = FinalAmount1 - finalamount
		database.DB.Model(&models.CartItem{}).Where("user_id = ? AND product_id = ?", userID, Cart.ProductID).Update("final_amount", FinalAmount1)
		var Discount1 float64
		database.DB.Model(&models.CartItem{}).Where("user_id = ? and product_id = ?", userID, Cart.ProductID).Pluck("discount", &Discount1)
		Discount1 = Discount1 - discount
		database.DB.Model(&models.CartItem{}).Where("user_id = ? AND product_id = ?", userID, Cart.ProductID).Update("discount", Discount1)

		c.JSON(http.StatusOK, gin.H{"status": true, "message": "product removed from cart successfully"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "product does not exist in cart"})
}

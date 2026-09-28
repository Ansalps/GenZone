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
	// Provide original price, unit price, discount (total), and final amount per item
	query := `
		SELECT 
			carts.user_id,
			cart_items.product_id,
			products.product_name,
			cart_items.quantity AS qty,
			cart_items.unit_price AS price,
			(products.price * cart_items.quantity) AS total_amount,
			((products.price - cart_items.unit_price) * cart_items.quantity) AS discount,
			(cart_items.quantity * cart_items.unit_price) AS final_amount
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

	// 4. Determine requested quantity (default 1) and check stock
	reqQty := req.Quantity
	if reqQty == 0 {
		reqQty = 1
	}

	// Check available stock for the requested size if provided, otherwise across all variants
	var stock uint
	if req.Size != "" {
		database.DB.Raw("SELECT COALESCE(stock,0) FROM product_variants WHERE product_id = ? AND size = ?", req.ProductID, req.Size).Scan(&stock)
	} else {
		database.DB.Raw("SELECT COALESCE(SUM(stock), 0) FROM product_variants WHERE product_id = ?", req.ProductID).Scan(&stock)
	}
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

	// 7. Check if the item (same product + size) is already in the user's cart
	var cartItem models.CartItem
	if req.Size != "" {
		err = database.DB.Where("cart_id = ? AND product_id = ? AND size = ?", cart.ID, req.ProductID, req.Size).First(&cartItem).Error
	} else {
		err = database.DB.Where("cart_id = ? AND product_id = ?", cart.ID, req.ProductID).First(&cartItem).Error
	}

	if err == nil {
		// Item exists -> Check quantity limits and stock
		if cartItem.Quantity+reqQty > 7 {
			c.JSON(http.StatusOK, gin.H{
				"status":  false,
				"message": "Exceeded maximum quantity for a product",
			})
			return
		}

		if cartItem.Quantity+reqQty > stock {
			c.JSON(http.StatusOK, gin.H{
				"status":  false,
				"message": "product out of stock",
			})
			return
		}

		// Update quantity and unit price
		cartItem.Quantity += reqQty
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

	// 8. Item does not exist -> Insert new record with requested quantity and size
	newCartItem := models.CartItem{
		CartID:    cart.ID,
		ProductID: product.ID,
		Quantity:  reqQty,
		UnitPrice: unitPrice,
		Size:      req.Size,
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
	var req requestmodemodels.CartAdd
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "failed to bind request"})
		return
	}

	if err := helper.Validate(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error(), "error_code": http.StatusBadRequest})
		return
	}

	// Fetch user's cart
	var cart models.Cart
	if err := database.DB.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": false, "message": "cart not found"})
		return
	}

	// Find the cart item (match size when provided)
	var cartItem models.CartItem
	if req.Size != "" {
		err := database.DB.Where("cart_id = ? AND product_id = ? AND size = ?", cart.ID, req.ProductID, req.Size).First(&cartItem).Error
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"status": true, "message": "product does not exist in cart"})
			return
		}
	} else {
		err := database.DB.Where("cart_id = ? AND product_id = ?", cart.ID, req.ProductID).First(&cartItem).Error
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"status": true, "message": "product does not exist in cart"})
			return
		}
	}

	// Quantity to remove (default 1)
	remQty := req.Quantity
	if remQty == 0 {
		remQty = 1
	}

	if cartItem.Quantity <= remQty {
		// remove the item
		if err := database.DB.Delete(&cartItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "failed to remove cart item"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": true, "message": "product removed from cart"})
		return
	}

	// Decrease quantity
	cartItem.Quantity = cartItem.Quantity - remQty
	if err := database.DB.Save(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "failed to update cart item quantity"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true, "message": "cart item quantity updated"})
}

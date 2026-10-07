package user

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/helper"
	"github.com/Ansalps/GeZOne/middleware"
	"github.com/Ansalps/GeZOne/models"
	requestmodemodels "github.com/Ansalps/GeZOne/requestmodels"
	"github.com/Ansalps/GeZOne/responsemodels"
	"github.com/gin-gonic/gin"
)

const (
	maxQuantityPerProduct = 5
	maxCartTotalQuantity  = 25
)

func validateCartAddition(currentItemQty, requestedQty, currentCartTotal uint) (bool, string) {
	if requestedQty == 0 {
		requestedQty = 1
	}

	if requestedQty > maxQuantityPerProduct {
		return false, fmt.Sprintf("Maximum quantity allowed for each product is %d", maxQuantityPerProduct)
	}

	if currentItemQty+requestedQty > maxQuantityPerProduct {
		return false, fmt.Sprintf("Maximum quantity allowed for this product is %d", maxQuantityPerProduct)
	}

	if currentCartTotal+requestedQty > maxCartTotalQuantity {
		return false, fmt.Sprintf("Cart total quantity cannot exceed %d", maxCartTotalQuantity)
	}

	return true, ""
}

func getUserCartTotalQuantity(userID uint) (uint, error) {
	var total int64
	err := database.DB.Table("carts").
		Joins("JOIN cart_items ON carts.id = cart_items.cart_id").
		Where("carts.user_id = ?", userID).
		Select("COALESCE(SUM(cart_items.quantity), 0)").
		Scan(&total).Error
	if err != nil {
		return 0, err
	}
	return uint(total), nil
}

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
			ci.id,
			ci.cart_id,
			c.user_id,
			ci.product_id,
			p.product_name,
			cat.category_name,
			p.description AS product_description,
			p.image_url AS product_image_url,
			p.popular,
			ci.size,
			COALESCE((SELECT pv2.stock FROM product_variants pv2 WHERE pv2.product_id = ci.product_id AND pv2.size = ci.size LIMIT 1), 0) AS stock,
			p.price AS original_price,
			ci.quantity AS qty,
			p.price AS price,
			CASE
				WHEN COALESCE(MAX(CASE WHEN o.start_at <= now() AND (o.end_at IS NULL OR o.end_at >= now()) THEN o.discount_percentage ELSE 0 END), 0) > 0
				THEN p.price * (1 - (COALESCE(MAX(CASE WHEN o.start_at <= now() AND (o.end_at IS NULL OR o.end_at >= now()) THEN o.discount_percentage ELSE 0 END), 0) / 100.0))
				ELSE p.price
			END AS offer_unit_price,
			(p.price * ci.quantity) AS original_total_price,
			(p.price * ci.quantity) AS total_amount,
			CASE
				WHEN COALESCE(MAX(CASE WHEN o.start_at <= now() AND (o.end_at IS NULL OR o.end_at >= now()) THEN o.discount_percentage ELSE 0 END), 0) > 0
				THEN (p.price * ci.quantity) * (1 - (COALESCE(MAX(CASE WHEN o.start_at <= now() AND (o.end_at IS NULL OR o.end_at >= now()) THEN o.discount_percentage ELSE 0 END), 0) / 100.0))
				ELSE p.price * ci.quantity
			END AS offer_total_price,
			((p.price - (CASE
				WHEN COALESCE(MAX(CASE WHEN o.start_at <= now() AND (o.end_at IS NULL OR o.end_at >= now()) THEN o.discount_percentage ELSE 0 END), 0) > 0
				THEN p.price * (1 - (COALESCE(MAX(CASE WHEN o.start_at <= now() AND (o.end_at IS NULL OR o.end_at >= now()) THEN o.discount_percentage ELSE 0 END), 0) / 100.0))
				ELSE p.price
			END)) * ci.quantity) AS discount,
			(ci.quantity * CASE
				WHEN COALESCE(MAX(CASE WHEN o.start_at <= now() AND (o.end_at IS NULL OR o.end_at >= now()) THEN o.discount_percentage ELSE 0 END), 0) > 0
				THEN p.price * (1 - (COALESCE(MAX(CASE WHEN o.start_at <= now() AND (o.end_at IS NULL OR o.end_at >= now()) THEN o.discount_percentage ELSE 0 END), 0) / 100.0))
				ELSE p.price
			END) AS final_amount,
			(ci.quantity * CASE
				WHEN COALESCE(MAX(CASE WHEN o.start_at <= now() AND (o.end_at IS NULL OR o.end_at >= now()) THEN o.discount_percentage ELSE 0 END), 0) > 0
				THEN p.price * (1 - (COALESCE(MAX(CASE WHEN o.start_at <= now() AND (o.end_at IS NULL OR o.end_at >= now()) THEN o.discount_percentage ELSE 0 END), 0) / 100.0))
				ELSE p.price
			END) AS final_price,
			COALESCE(string_agg(DISTINCT pv.size, ',' ORDER BY pv.size), '') AS size_options,
			COALESCE(json_agg(DISTINCT jsonb_build_object('size', pv.size, 'stock', pv.stock)) FILTER (WHERE pv.id IS NOT NULL), '[]') AS inventory,
			COALESCE(MAX(CASE WHEN o.start_at <= now() AND (o.end_at IS NULL OR o.end_at >= now()) THEN o.discount_percentage ELSE 0 END), 0) AS discount_percentage,
			COALESCE(MAX(to_char(o.start_at, 'YYYY-MM-DD')), '') AS start_date,
			COALESCE(MAX(to_char(o.end_at, 'YYYY-MM-DD')), '') AS end_date
		FROM carts c
		JOIN cart_items ci ON c.id = ci.cart_id
		JOIN products p ON ci.product_id = p.id
		LEFT JOIN categories cat ON cat.id = p.category_id
		LEFT JOIN (
			SELECT product_id, SUM(stock) AS stock_total
			FROM product_variants
			GROUP BY product_id
		) v ON v.product_id = p.id
		LEFT JOIN product_variants pv ON pv.product_id = p.id
		LEFT JOIN offers o ON o.product_id = p.id
		WHERE c.user_id = ? AND ci.quantity > 0
		GROUP BY
			ci.id,
			ci.cart_id,
			c.user_id,
			ci.product_id,
			p.product_name,
			cat.category_name,
			p.description,
			p.image_url,
			p.popular,
			ci.size,
			ci.quantity,
			ci.unit_price,
			p.price,
			v.stock_total
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

	currentCartTotal, err := getUserCartTotalQuantity(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "failed to calculate cart total quantity",
		})
		return
	}

	// 6. Check for an active offer and calculate unit price
	unitPrice := product.Price
	var offer models.Offer
	now := time.Now()

	err = database.DB.Where("product_id = ? AND start_at <= ? AND end_at >= ?", product.ID, now, now).
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
		if ok, msg := validateCartAddition(cartItem.Quantity, reqQty, currentCartTotal); !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": msg,
			})
			return
		}

		// Item exists -> Check quantity limits and stock
		if cartItem.Quantity+reqQty > stock {
			c.JSON(http.StatusBadRequest, gin.H{
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

	if ok, msg := validateCartAddition(0, reqQty, currentCartTotal); !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": msg,
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

func UpdateQuantity(c *gin.Context) {
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

	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "failed to bind request"})
		return
	}

	productIDValue, ok := payload["product_id"]
	if !ok || productIDValue == nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "product_id is required"})
		return
	}

	quantityValue, ok := payload["quantity"]
	if !ok || quantityValue == nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "quantity is required"})
		return
	}

	productID, err := strconv.ParseUint(fmt.Sprintf("%v", productIDValue), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "invalid product_id"})
		return
	}

	quantity, err := strconv.ParseUint(fmt.Sprintf("%v", quantityValue), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "invalid quantity"})
		return
	}

	if quantity == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "quantity must be at least 1"})
		return
	}

	if uint(quantity) > maxQuantityPerProduct {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": fmt.Sprintf("Maximum quantity allowed for this product is %d", maxQuantityPerProduct),
		})
		return
	}

	sizeValue, _ := payload["size"].(string)
	size := ""
	if sizeValue != "" {
		size = sizeValue
	}

	var cart models.Cart
	if err := database.DB.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": false, "message": "cart not found"})
		return
	}

	var cartItem models.CartItem
	if size != "" {
		err := database.DB.Where("cart_id = ? AND product_id = ? AND size = ?", cart.ID, productID, size).First(&cartItem).Error
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"status": false, "message": "item not found in cart"})
			return
		}
	} else {
		err := database.DB.Where("cart_id = ? AND product_id = ?", cart.ID, productID).First(&cartItem).Error
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"status": false, "message": "item not found in cart"})
			return
		}
	}

	var availableStock uint
	if size != "" {
		err := database.DB.Raw("SELECT COALESCE(stock, 0) FROM product_variants WHERE product_id = ? AND size = ?", productID, size).Scan(&availableStock).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "failed to check stock availability"})
			return
		}
	} else {
		err := database.DB.Raw("SELECT COALESCE(SUM(stock), 0) FROM product_variants WHERE product_id = ?", productID).Scan(&availableStock).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "failed to check stock availability"})
			return
		}
	}

	if uint(quantity) > availableStock {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": fmt.Sprintf("Quantity cannot be more than available stock (%d). Please reduce the quantity and try again.", availableStock),
		})
		return
	}

	cartItem.Quantity = uint(quantity)
	if err := database.DB.Save(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "failed to update product quantity"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "product quantity updated successfully",
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
	var req requestmodemodels.CartRemove
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

func CartSummary(c *gin.Context) {
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

	var summary struct {
		ItemsCount     int64   `json:"items_count"`
		TotalQuantity  int64   `json:"total_quantity"`
		OriginalTotal  float64 `json:"original_total"`
		DiscountAmount float64 `json:"discount_amount"`
		FinalTotal     float64 `json:"final_total"`
	}

	query := `
		WITH item_summary AS (
			SELECT
				ci.id,
				ci.product_id,
				ci.quantity,
				p.price AS original_unit_price,
				CASE
					WHEN COALESCE(MAX(CASE WHEN o.start_at <= now() AND (o.end_at IS NULL OR o.end_at >= now()) THEN o.discount_percentage ELSE 0 END), 0) > 0
					THEN p.price * (1 - (COALESCE(MAX(CASE WHEN o.start_at <= now() AND (o.end_at IS NULL OR o.end_at >= now()) THEN o.discount_percentage ELSE 0 END), 0) / 100.0))
					ELSE p.price
				END AS offer_unit_price,
				p.price * ci.quantity AS original_total,
				CASE
					WHEN COALESCE(MAX(CASE WHEN o.start_at <= now() AND (o.end_at IS NULL OR o.end_at >= now()) THEN o.discount_percentage ELSE 0 END), 0) > 0
					THEN (p.price * ci.quantity) * (1 - (COALESCE(MAX(CASE WHEN o.start_at <= now() AND (o.end_at IS NULL OR o.end_at >= now()) THEN o.discount_percentage ELSE 0 END), 0) / 100.0))
					ELSE p.price * ci.quantity
				END AS final_total
			FROM carts c
			JOIN cart_items ci ON c.id = ci.cart_id
			JOIN products p ON ci.product_id = p.id
			LEFT JOIN offers o ON o.product_id = p.id
			WHERE c.user_id = ? AND ci.quantity > 0
			GROUP BY ci.id, ci.product_id, ci.quantity, p.price
		)
		SELECT
			COUNT(id) AS items_count,
			COALESCE(SUM(quantity), 0) AS total_quantity,
			COALESCE(SUM(original_total), 0) AS original_total,
			COALESCE(SUM(original_total - final_total), 0) AS discount_amount,
			COALESCE(SUM(final_total), 0) AS final_total
		FROM item_summary`

	if err := database.DB.Raw(query, userID).Scan(&summary).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "failed to calculate cart summary",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "successfully retrieved cart summary",
		"data": gin.H{
			"items_count":     summary.ItemsCount,
			"total_quantity":  summary.TotalQuantity,
			"original_total":  summary.OriginalTotal,
			"discount_amount": summary.DiscountAmount,
			"final_total":     summary.FinalTotal,
		},
	})
}

func TotalQuantity(c *gin.Context) {
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

	var totalQuantity int64
	err := database.DB.Table("carts").
		Joins("JOIN cart_items ON carts.id = cart_items.cart_id").
		Where("carts.user_id = ?", userID).
		Select("COALESCE(SUM(cart_items.quantity), 0)").
		Scan(&totalQuantity).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "failed to calculate total quantity",
		})
		return
	}
	fmt.Println("Total Quantity: ", totalQuantity)
	c.JSON(http.StatusOK, gin.H{
		"status":         true,
		"message":        "successfully retrieved total quantity",
		"total_quantity": totalQuantity,
	})
}

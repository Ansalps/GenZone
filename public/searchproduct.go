package public

import (
	"net/http"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/responsemodels"
	"github.com/gin-gonic/gin"
)

func SearchProduct(c *gin.Context) {

	search := c.Query("search")
	nameSort := c.Query("name_sort")
	priceSort := c.Query("price_sort")
	newArrivals := c.Query("new_arrivals")
	category := c.Query("category")

	var products []responsemodels.Product

	sql := `
		SELECT
			p.id,
			p.created_at,
			p.updated_at,
			p.category_id,
			c.category_name,
			p.product_name,
			p.description AS product_description,
			p.image_url AS product_image_url,
			p.price,
			p.stock,
			p.popular,
			p.size,
			COALESCE(o.discount_percentage, 0) AS discount_percentage
		FROM products p
		JOIN categories c
			ON p.category_id = c.id
		LEFT JOIN offers o
			ON p.id = o.product_id
			AND o.deleted_at IS NULL
	`

	var args []interface{}

	// --------------------------------
	// Search
	// --------------------------------

	if search != "" {
		sql += ` AND p.product_name ILIKE ?`
		args = append(args, "%"+search+"%")
	}

	// --------------------------------
	// Category
	// --------------------------------

	if category != "" {

		// Check category exists
		var count int64

		err := database.DB.
			Raw(`
				SELECT COUNT(*)
				FROM categories
				WHERE category_name = ?
			`, category).
			Scan(&count).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to check category",
			})
			return
		}

		if count == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Such category does not exist",
			})
			return
		}

		sql += ` AND c.category_name = ?`
		args = append(args, category)
	}

	// --------------------------------
	// Sorting
	// --------------------------------

	var orderBy []string

	if nameSort != "" {

		switch nameSort {

		case "aA-zZ":
			orderBy = append(
				orderBy,
				"p.product_name ASC",
			)

		case "zZ-aA":
			orderBy = append(
				orderBy,
				"p.product_name DESC",
			)

		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Invalid name sort filter",
			})
			return
		}
	}

	if priceSort != "" {

		switch priceSort {

		case "low-high":
			orderBy = append(
				orderBy,
				"p.price ASC",
			)

		case "high-low":
			orderBy = append(
				orderBy,
				"p.price DESC",
			)

		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Invalid price sort filter",
			})
			return
		}
	}

	if newArrivals != "" {

		switch newArrivals {

		case "true":
			orderBy = append(
				orderBy,
				"p.created_at DESC",
			)

		case "false":
			// No sorting required for false

		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "new_arrivals must be true or false",
			})
			return
		}
	}

	// Default sorting
	if len(orderBy) == 0 {
		orderBy = append(
			orderBy,
			"p.created_at DESC",
		)
	}

	sql += ` ORDER BY `

	for i, order := range orderBy {

		if i > 0 {
			sql += `, `
		}

		sql += order
	}

	// --------------------------------
	// Execute query
	// --------------------------------

	if err := database.DB.
		Raw(sql, args...).
		Scan(&products).
		Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to retrieve products",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Products retrieved successfully",
		"data": gin.H{
			"products": products,
		},
	})
}
package user

import (
	"net/http"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/responsemodels"
	"github.com/gin-gonic/gin"
)

func SearchProduct(c *gin.Context) {
	query := c.Query("q")
	nameSort := c.Query("name_sort")
	priceSort := c.Query("price_sort")
	newArrivals := c.Query("new_arrivals")
	category := c.Query("category")
	var products []responsemodels.Product

	sql := `SELECT products.*, categories.category_name AS category_name FROM products
            JOIN categories ON products.category_id = categories.id`

	// Apply search query filter
	if query != "" {
		sql += ` WHERE (products.product_name ILIKE '%` + query + `%' OR products.description ILIKE '%` + query + `%')`
	}

	// Apply category filter
	if category != "" {
		if query != "" {
			sql += ` AND categories.category_name = '` + category + `'`
		} else {
			sql += ` WHERE categories.category_name = '` + category + `'`
		}
	}

	// Apply new arrivals filter
	if newArrivals == "true" {
		sql += ` ORDER BY products.created_at DESC`
	}

	// Apply name sort filter
	if nameSort == "aA-zZ" {
		sql += ` ORDER BY products.product_name ASC`
	} else if nameSort == "zZ-aA" {
		sql += ` ORDER BY products.product_name DESC`
	}

	// Apply price sort filter
	if priceSort == "low-high" {
		sql += ` ORDER BY products.price ASC`
	} else if priceSort == "high-low" {
		sql += ` ORDER BY products.price DESC`
	}

	// Execute the raw SQL query
	database.DB.Raw(sql).Scan(&products)

	c.JSON(http.StatusOK, products)
}

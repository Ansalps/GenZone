package user

import (
	"fmt"
	"net/http"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/responsemodels"
	"github.com/gin-gonic/gin"
)

func SearchProduct(c *gin.Context) {
	query := c.Query("search")
	nameSort := c.Query("name_sort")
	priceSort := c.Query("price_sort")
	newArrivals := c.Query("new_arrivals")
	category := c.Query("category")
	var products []responsemodels.Product

	sql := `SELECT products.*, categories.category_name AS category_name FROM products
            JOIN categories ON products.category_id = categories.id`

	// Apply search query filter
	if query != "" {
		sql += ` WHERE (products.product_name ILIKE '%` + query + `%')`
	}

	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM categories WHERE category_name = ? AND deleted_at IS NULL`, category).Scan(&count)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "such category does not exist",
		})
		return
	}
	// Apply category filter
	fmt.Println("category", category)
	if category != "" {
		if query != "" {
			sql += ` AND categories.category_name = '` + category + `'`
		} else {
			sql += ` WHERE categories.category_name = '` + category + `'`
		}
	}
	if nameSort != "aA-zZ" && nameSort != "zZ-aA" {
		c.JSON(http.StatusBadRequest, gin.H{
			"messsage": "Apply correct names sort filter",
		})
		return
	}
	// Apply name sort filter
	if nameSort == "aA-zZ" {
		sql += ` ORDER BY products.product_name ASC`
	} else if nameSort == "zZ-aA" {
		sql += ` ORDER BY products.product_name DESC`
	}
	if priceSort != "low-high" && priceSort != "high-low" {
		c.JSON(http.StatusBadRequest, gin.H{
			"messsage": "Apply correct price sort filter",
		})
		return
	}
	// Apply price sort filter
	if priceSort == "low-high" {
		if nameSort != "" {
			sql += `, products.price ASC`
		} else {
			sql += ` ORDER BY products.price ASC`
		}

	} else if priceSort == "high-low" {
		if nameSort != "" {
			sql += `, products.price DESC`
		} else {
			sql += ` ORDER BY products.price DESC`
		}
	}
	if newArrivals != "true" && newArrivals != "false" {
		c.JSON(http.StatusBadRequest, gin.H{
			"messsage": "new arrivals can either be true or false",
		})
		return
	}
	// Apply new arrivals filter
	if newArrivals == "true" {
		if nameSort != "" || priceSort != "" {
			sql += `, products.created_at DESC`
		} else {
			sql += ` ORDER BY products.created_at DESC`
		}
	}

	// Execute the raw SQL query
	database.DB.Raw(sql).Scan(&products)

	c.JSON(http.StatusOK, products)
}

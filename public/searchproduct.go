package public

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/models"
	"github.com/Ansalps/GeZOne/responsemodels"
	"github.com/gin-gonic/gin"
)

func GetProduct(c *gin.Context) {
	search := c.Query("search")
	nameSort := c.Query("name_sort")
	priceSort := c.Query("price_sort")
	newArrivals := c.Query("new_arrivals")
	category := c.Query("category")

	var products []responsemodels.Product

	// Added WHERE 1=1 to ensure valid SQL syntax for subsequent AND filters
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
			COALESCE((SELECT COALESCE(SUM(stock),0) FROM product_variants WHERE product_id = p.id), 0) AS stock,
            p.popular,
				COALESCE(string_agg(DISTINCT pv.size, ',' ORDER BY pv.size), '') AS size,
				-- inventory as JSON array of {size, stock}
				COALESCE(json_agg(DISTINCT jsonb_build_object('size', pv.size, 'stock', pv.stock)) FILTER (WHERE pv.id IS NOT NULL), '[]') AS inventory,
				-- Only expose discount_percentage when the offer is currently active
				COALESCE(MAX(CASE WHEN o.start_at <= now() AND (o.end_at IS NULL OR o.end_at >= now()) THEN o.discount_percentage ELSE 0 END), 0) AS discount_percentage,
				COALESCE(MAX(to_char(o.start_at, 'YYYY-MM-DD')), '') AS start_date,
				COALESCE(MAX(to_char(o.end_at, 'YYYY-MM-DD')), '') AS end_date
        FROM products p
        JOIN categories c
            ON p.category_id = c.id
        LEFT JOIN product_variants pv
            ON pv.product_id = p.id
        LEFT JOIN offers o
            ON p.id = o.product_id
        WHERE 1=1
    `

	var args []interface{}

	// --------------------------------
	// Search Filter
	// --------------------------------
	if search != "" {
		sql += ` AND p.product_name ILIKE ?`
		args = append(args, "%"+search+"%")
	}

	// --------------------------------
	// Category Filter
	// --------------------------------
	if category != "" {
		var count int64
		err := database.DB.
			Raw(`SELECT COUNT(*) FROM categories WHERE category_name = ?`, category).
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
	// Group By
	// --------------------------------
	sql += ` GROUP BY p.id, p.created_at, p.updated_at, p.category_id, c.category_name,
		p.product_name, p.description, p.image_url, p.price, p.popular, o.discount_percentage, to_char(o.start_at, 'YYYY-MM-DD'), to_char(o.end_at, 'YYYY-MM-DD') `

	// --------------------------------
	// Sorting
	// --------------------------------
	var orderBy []string

	if nameSort != "" {
		switch nameSort {
		case "aA-zZ":
			orderBy = append(orderBy, "p.product_name ASC")
		case "zZ-aA":
			orderBy = append(orderBy, "p.product_name DESC")
		default:
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid name sort filter"})
			return
		}
	}

	if priceSort != "" {
		switch priceSort {
		case "low-high":
			orderBy = append(orderBy, "p.price ASC")
		case "high-low":
			orderBy = append(orderBy, "p.price DESC")
		default:
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid price sort filter"})
			return
		}
	}

	if newArrivals != "" {
		switch newArrivals {
		case "true":
			orderBy = append(orderBy, "p.created_at DESC")
		case "false":
			// No sorting required
		default:
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "new_arrivals must be true or false"})
			return
		}
	}

	if len(orderBy) == 0 {
		orderBy = append(orderBy, "p.created_at DESC")
	}

	sql += ` ORDER BY ` + strings.Join(orderBy, ", ")

	// --------------------------------
	// Execute query
	// --------------------------------
	if err := database.DB.Raw(sql, args...).Scan(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to retrieve products",
		})
		return
	}

	// Try to unmarshal inventory JSON returned by the SQL into Product.Inventory
	for i := range products {
		if len(products[i].InventoryRaw) > 0 {
			var inv []responsemodels.ProductInventoryItem
			if err := json.Unmarshal(products[i].InventoryRaw, &inv); err == nil {
				products[i].Inventory = inv
				continue
			}
		}
	}

	// Batch fetch variants and attach inventory per product
	if len(products) > 0 {
		ids := make([]uint, 0, len(products))
		for _, p := range products {
			ids = append(ids, p.ID)
		}

		var variants []models.ProductVariant
		if err := database.DB.Where("product_id IN ?", ids).Order("product_id ASC, size ASC").Find(&variants).Error; err == nil {
			variantMap := make(map[uint][]responsemodels.ProductInventoryItem)
			for _, v := range variants {
				variantMap[v.ProductID] = append(variantMap[v.ProductID], responsemodels.ProductInventoryItem{Size: v.Size, Stock: v.Stock})
			}

			for i := range products {
				if inv, ok := variantMap[products[i].ID]; ok {
					products[i].Inventory = inv
				} else {
					products[i].Inventory = []responsemodels.ProductInventoryItem{}
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Products retrieved successfully",
		"data": gin.H{
			"products": products,
		},
	})
}


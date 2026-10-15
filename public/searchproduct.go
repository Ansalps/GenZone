package public

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/helper"
	"github.com/Ansalps/GeZOne/models"
	"github.com/Ansalps/GeZOne/responsemodels"
	"github.com/gin-gonic/gin"
)

func resolveProductImageURL(ctx context.Context, imageURL string) string {
	imageURL = strings.TrimSpace(imageURL)
	if imageURL == "" {
		return ""
	}
	if strings.HasPrefix(imageURL, "http://") || strings.HasPrefix(imageURL, "https://") {
		return imageURL
	}

	if resolvedURL, err := helper.ResolveObjectURL(ctx, "S3_BUCKET_NAME", imageURL); err == nil {
		return resolvedURL
	}

	return imageURL
}

func GetProduct(c *gin.Context) {
	search := c.Query("search")
	nameSort := c.Query("name_sort")
	priceSort := c.Query("price_sort")
	newArrivals := c.Query("new_arrivals")
	category := c.Query("category")

	var products []responsemodels.Product

	// Added WHERE 1=1 to ensure valid SQL syntax for subsequent AND filters
	var sql string
	if database.DB.Dialector.Name() == "sqlite" {
		sql = `
        SELECT
            p.id,
            p.created_at,
            p.updated_at,
            p.category_id,
			p.deleted_at,
            c.category_name,
            p.product_name,
            CASE WHEN p.deleted_at IS NOT NULL THEN 'discontinued' ELSE 'active' END AS product_status,
            p.description AS product_description,
            p.image_url AS product_image_url,
			p.price,
			COALESCE((SELECT COALESCE(SUM(stock),0) FROM product_variants WHERE product_id = p.id), 0) AS stock,
            p.popular,
				COALESCE(GROUP_CONCAT(DISTINCT pv.size), '') AS size,
				COALESCE(MAX(CASE WHEN o.start_at <= datetime('now') AND (o.end_at IS NULL OR o.end_at >= datetime('now')) THEN o.discount_percentage ELSE 0 END), 0) AS discount_percentage,
				COALESCE(MAX(strftime('%Y-%m-%d', o.start_at)), '') AS start_date,
				COALESCE(MAX(strftime('%Y-%m-%d', o.end_at)), '') AS end_date
        FROM products p
        JOIN categories c
            ON p.category_id = c.id AND c.deleted_at IS NULL
        LEFT JOIN product_variants pv
            ON pv.product_id = p.id AND pv.deleted_at IS NULL
        LEFT JOIN offers o
            ON p.id = o.product_id AND o.deleted_at IS NULL
    `
	} else {
		sql = `
        SELECT
            p.id,
            p.created_at,
            p.updated_at,
            p.category_id,
			p.deleted_at,
            c.category_name,
            p.product_name,
            CASE WHEN p.deleted_at IS NOT NULL THEN 'discontinued' ELSE 'active' END AS product_status,
            p.description AS product_description,
            p.image_url AS product_image_url,
			p.price,
			COALESCE((SELECT COALESCE(SUM(stock),0) FROM product_variants WHERE product_id = p.id), 0) AS stock,
            p.popular,
				COALESCE(string_agg(DISTINCT pv.size, ',' ORDER BY pv.size), '') AS size,
				COALESCE(json_agg(DISTINCT jsonb_build_object('size', pv.size, 'stock', pv.stock)) FILTER (WHERE pv.id IS NOT NULL), '[]') AS inventory,
				COALESCE(MAX(CASE WHEN o.start_at <= now() AND (o.end_at IS NULL OR o.end_at >= now()) THEN o.discount_percentage ELSE 0 END), 0) AS discount_percentage,
				COALESCE(MAX(to_char(o.start_at, 'YYYY-MM-DD')), '') AS start_date,
				COALESCE(MAX(to_char(o.end_at, 'YYYY-MM-DD')), '') AS end_date
        FROM products p
        JOIN categories c
            ON p.category_id = c.id AND c.deleted_at IS NULL
        LEFT JOIN product_variants pv
            ON pv.product_id = p.id AND pv.deleted_at IS NULL
        LEFT JOIN offers o
            ON p.id = o.product_id AND o.deleted_at IS NULL
    `
	}

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
			Raw(`SELECT COUNT(*) FROM categories WHERE category_name = ? AND deleted_at IS NULL`, category).
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
	if database.DB.Dialector.Name() == "sqlite" {
		sql += ` GROUP BY p.id, p.created_at, p.updated_at, p.category_id, p.deleted_at, c.category_name,
		p.product_name, p.description, p.image_url, p.price, p.popular, o.discount_percentage,
		strftime('%Y-%m-%d', o.start_at), strftime('%Y-%m-%d', o.end_at) `
	} else {
		sql += ` GROUP BY p.id, p.created_at, p.updated_at, p.category_id, p.deleted_at, c.category_name,
		p.product_name, p.description, p.image_url, p.price, p.popular, o.discount_percentage,
		to_char(o.start_at, 'YYYY-MM-DD'), to_char(o.end_at, 'YYYY-MM-DD') `
	}

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

	for i := range products {
		products[i].ProductImageUrl = resolveProductImageURL(c.Request.Context(), products[i].ProductImageUrl)
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

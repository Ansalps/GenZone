package public

import (
	"encoding/json"
	"net/http"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/helper"
	"github.com/Ansalps/GeZOne/responsemodels"
	"github.com/gin-gonic/gin"
)

func ReadProducts(c *gin.Context) {

	var products []responsemodels.Product
	var sql string

	if database.DB.Dialector.Name() == "sqlite" {
		sql = `
		SELECT
			p.id,
			p.created_at,
			p.updated_at,
			p.deleted_at,
			p.category_id,
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
			ON c.id = p.category_id AND c.deleted_at IS NULL
		LEFT JOIN product_variants pv
			ON pv.product_id = p.id AND pv.deleted_at IS NULL
		LEFT JOIN offers o
			ON o.product_id = p.id AND o.deleted_at IS NULL
		GROUP BY p.id, p.created_at, p.updated_at, p.category_id, c.category_name,
			p.product_name, p.deleted_at, p.description, p.image_url, p.price, p.popular
		ORDER BY p.id DESC
	`
	} else {
		sql = `
		SELECT
			p.id,
			p.created_at,
			p.updated_at,
			p.deleted_at,
			p.category_id,
			c.category_name,
			p.product_name,
			CASE WHEN p.deleted_at IS NOT NULL THEN 'discontinued' ELSE 'active' END AS product_status,
			p.description AS product_description,
			p.image_url AS product_image_url,
			p.price,
			COALESCE((SELECT COALESCE(SUM(stock),0) FROM product_variants WHERE product_id = p.id), 0) AS stock,
			p.popular,
			COALESCE(string_agg(DISTINCT pv.size, ',' ORDER BY pv.size), '') AS size,
			COALESCE(MAX(CASE WHEN o.start_at <= now() AND (o.end_at IS NULL OR o.end_at >= now()) THEN o.discount_percentage ELSE 0 END), 0) AS discount_percentage,
			COALESCE(MAX(to_char(o.start_at, 'YYYY-MM-DD')), '') AS start_date,
			COALESCE(MAX(to_char(o.end_at, 'YYYY-MM-DD')), '') AS end_date
		FROM products p
		JOIN categories c
			ON c.id = p.category_id AND c.deleted_at IS NULL
		LEFT JOIN product_variants pv
			ON pv.product_id = p.id AND pv.deleted_at IS NULL
		LEFT JOIN offers o
			ON o.product_id = p.id AND o.deleted_at IS NULL
		GROUP BY p.id, p.created_at, p.updated_at, p.category_id, c.category_name,
			p.product_name, p.deleted_at, p.description, p.image_url, p.price, p.popular
		ORDER BY p.id DESC
	`
	}

	if err := database.DB.Raw(sql).Scan(&products).Error; err != nil {
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

	// Batch fetch inventory JSON per product using SQL aggregation
	if len(products) > 0 {
		ids := make([]uint, 0, len(products))
		for _, p := range products {
			ids = append(ids, p.ID)
		}

		type invRow struct {
			ProductID uint            `gorm:"column:product_id"`
			Inv       json.RawMessage `gorm:"column:inv"`
		}

		var rowsInv []invRow
		rawSql := `SELECT product_id, json_agg(jsonb_build_object('size', size, 'stock', stock) ORDER BY size) AS inv FROM product_variants WHERE product_id IN ? GROUP BY product_id`
		if database.DB.Dialector.Name() == "sqlite" {
			rawSql = `SELECT product_id, json_group_array(json_object('size', size, 'stock', stock)) AS inv FROM product_variants WHERE product_id IN ? GROUP BY product_id`
		} else {
			rawSql = `SELECT product_id, json_agg(jsonb_build_object('size', size, 'stock', stock) ORDER BY size) AS inv FROM product_variants WHERE product_id IN ? GROUP BY product_id`
		}

		if err := database.DB.Raw(rawSql, ids).Scan(&rowsInv).Error; err == nil {
			invMap := make(map[uint]json.RawMessage)
			for _, r := range rowsInv {
				invMap[r.ProductID] = r.Inv
			}

			for i := range products {
				if raw, ok := invMap[products[i].ID]; ok && len(raw) > 0 {
					var inv []responsemodels.ProductInventoryItem
					if err := json.Unmarshal(raw, &inv); err == nil {
						products[i].Inventory = inv
					}
				}
				if products[i].Inventory == nil {
					products[i].Inventory = []responsemodels.ProductInventoryItem{}
				}
			}
		}
	}

	// Build response objects explicitly to ensure inventory is included
	productsResp := make([]map[string]interface{}, 0, len(products))
	for _, p := range products {
		productImageURL := p.ProductImageUrl
		if productImageURL != "" {
			if resolvedURL, err := helper.ResolveObjectURL(c.Request.Context(), "S3_BUCKET_NAME", productImageURL); err == nil {
				productImageURL = resolvedURL
			}
		}

		prod := map[string]interface{}{
			"id":                  p.ID,
			"created_at":          p.CreatedAt,
			"updated_at":          p.UpdatedAt,
			"deleted_at":          p.DeletedAt,
			"category_id":         p.CategoryID,
			"category_name":       p.CategoryName,
			"product_name":        p.ProductName,
			"product_status":      p.ProductStatus,
			"product_description": p.ProductDescription,
			"product_image_url":   productImageURL,
			"price":               p.Price,
			"stock":               p.Stock,
			"popular":             p.Popular,
			"size":                p.Size,
			"inventory":           p.Inventory,
			"discount_percentage": p.DiscountPercentage,
			"start_date":          p.StartDate,
			"end_date":            p.EndDate,
		}
		productsResp = append(productsResp, prod)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Successfully retrieved products",
		"data": gin.H{
			"products": productsResp,
		},
	})
}

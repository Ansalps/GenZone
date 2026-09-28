package public

import (
	"encoding/json"
	"net/http"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/responsemodels"
	"github.com/gin-gonic/gin"
)

func ReadProducts(c *gin.Context) {

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
			COALESCE((SELECT COALESCE(SUM(stock),0) FROM product_variants WHERE product_id = p.id), 0) AS stock,
			p.popular,
			COALESCE(string_agg(DISTINCT pv.size, ',' ORDER BY pv.size), '') AS size,
			-- inventory as JSON array of {size, stock}
			COALESCE(json_agg(DISTINCT jsonb_build_object('size', pv.size, 'stock', pv.stock)) FILTER (WHERE pv.id IS NOT NULL), '[]') AS inventory,
			COALESCE(MAX(CASE WHEN o.start_at <= now() AND (o.end_at IS NULL OR o.end_at >= now()) THEN o.discount_percentage ELSE 0 END), 0) AS discount_percentage,
			COALESCE(MAX(to_char(o.start_at, 'YYYY-MM-DD')), '') AS start_date,
			COALESCE(MAX(to_char(o.end_at, 'YYYY-MM-DD')), '') AS end_date
		FROM products p
		JOIN categories c
			ON c.id = p.category_id
		LEFT JOIN product_variants pv
			ON pv.product_id = p.id
		LEFT JOIN offers o
			ON o.product_id = p.id
		GROUP BY p.id, p.created_at, p.updated_at, p.category_id, c.category_name,
			p.product_name, p.description, p.image_url, p.price, p.popular
		ORDER BY p.id DESC
	`

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
		prod := map[string]interface{}{
			"id":                  p.ID,
			"created_at":          p.CreatedAt,
			"updated_at":          p.UpdatedAt,
			"category_id":         p.CategoryID,
			"category_name":       p.CategoryName,
			"product_name":        p.ProductName,
			"product_description": p.ProductDescription,
			"product_image_url":   p.ProductImageUrl,
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

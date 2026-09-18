package public

import (
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
			COALESCE(SUM(pv.stock), 0) AS stock,
			p.popular,
			COALESCE(string_agg(DISTINCT pv.size, ',' ORDER BY pv.size), '') AS size,
			COALESCE(o.discount_percentage, 0) AS discount_percentage
		FROM products p
		JOIN categories c
			ON c.id = p.category_id
		LEFT JOIN product_variants pv
			ON pv.product_id = p.id
		LEFT JOIN offers o
			ON o.product_id = p.id
		GROUP BY p.id, p.created_at, p.updated_at, p.category_id, c.category_name,
			p.product_name, p.description, p.image_url, p.price, p.popular, o.discount_percentage
		ORDER BY p.id DESC
	`

	if err := database.DB.Raw(sql).Scan(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to retrieve products",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Successfully retrieved products",
		"data": gin.H{
			"products": products,
		},
	})
}

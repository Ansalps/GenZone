package admin

import (
	"net/http"
	"strings"
	"time"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/models"
	"github.com/Ansalps/GeZOne/responsemodels"
	"github.com/gin-gonic/gin"
)

func BestSelling(c *gin.Context) {
	day := c.Query("day")
	month := c.Query("month")
	year := c.Query("year")
	// A month selection always refers to that month in the current year.
	dateYear := year
	if month != "" {
		dateYear = ""
	}
	// Best-selling reports support a specific day for backwards compatibility,
	// plus a month in the current year or an entire selected year. Period is
	// intentionally not accepted as a filter.
	if _, _, err := resolveReportDateRange("", day, month, dateYear, time.Now()); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var orderitems []models.OrderItems
	var orderitems1 []responsemodels.BestSelling
	sql := `select product_id,count(*) AS count from order_items where order_status='delivered'`
	sqla := ` group by product_id order by count DESC limit 10`
	sql1 := `select count(*),categories.category_name from order_items join products on products.id = order_items.product_id join categories on categories.id = products.category_id where order_items.order_status='delivered'`
	sql1a := ` group by categories.category_name order by count desc`

	dateFilter, dateArgs, err := buildReportDateFilter("order_items.created_at", "", day, month, dateYear, time.Now())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if dateFilter != "" {
		// Both base queries already contain a WHERE clause, so append the
		// date predicates with AND rather than introducing a second WHERE.
		dateCondition := " AND " + strings.TrimPrefix(dateFilter, " WHERE ")
		sql += strings.Replace(dateCondition, "order_items.created_at", "created_at", 1) + sqla
		sql1 += dateCondition + sql1a
		database.DB.Raw(sql, dateArgs...).Scan(&orderitems)
		database.DB.Raw(sql1, dateArgs...).Scan(&orderitems1)
	} else {
		database.DB.Raw(sql + sqla).Scan(&orderitems)
		database.DB.Raw(sql1 + sql1a).Scan(&orderitems1)
	}

	var product1 responsemodels.Product
	var products []responsemodels.Product
	for _, v := range orderitems {
		database.DB.Raw("SELECT * FROM products join categories on products.category_id=categories.id WHERE products.id = ?", v.ProductID).Scan(&product1)
		products = append(products, product1)
	}
	var category []string
	for _, v := range orderitems1 {
		category = append(category, v.CategoryName)
	}
	c.JSON(http.StatusOK, gin.H{
		"best_selling_product":  products,
		"best_selling_category": category,
		"message":               "retrieved best selling products",
	})

}

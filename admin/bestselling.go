package admin

import (
	"fmt"
	"net/http"
	"strconv"
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

	// Attempt to parse the date string
	if day != "" {
		_, err := time.Parse("2006-01-02", day)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid Date Format",
			}) // The date is invalid
			return
		}
	}

	if month != "" {
		int1, err := strconv.Atoi(month)
		if err != nil {
			fmt.Println("Error converting string to int:", err)
		}

		if int1 >= 1 && int1 <= 12 {
			fmt.Println("ok")
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid month format",
			})
			return
		}
	}

	if year != "" {
		// Parse the year string to an integer
		yearint, err := strconv.Atoi(year)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid year format",
			})
			return
		}

		// Get the current year
		currentYear := time.Now().Year()

		if yearint < 1900 || yearint > currentYear {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "year out of range, must be between 1900 and current year",
			})
			return
		}
	}

	var orderitems []models.OrderItems
	var orderitems1 []responsemodels.BestSelling
	//database.DB.Find(&orderitems)
	// sql := `select product_id,count(*) AS count from order_items where order_status='delivered' group by product_id order by count DESC limit 10`
	sql := `select product_id,count(*) AS count from order_items where order_status='delivered'`
	sqla := ` group by product_id order by count DESC limit 10`
	// sql1 := `select count(*),categories.category_name from order_items join products on products.id = order_items.product_id join categories on categories.id = products.category_id where order_items.order_status='delivered' group by categories.category_name order by count desc`
	sql1 := `select count(*),categories.category_name from order_items join products on products.id = order_items.product_id join categories on categories.id = products.category_id where order_items.order_status='delivered'`
	sql1a := ` group by categories.category_name order by count desc`
	if day != "" {
		sql = sql + ` AND DATE(created_at) ='` + day + `'` + sqla
		database.DB.Raw(sql).Scan(&orderitems)
		sql1 = sql1 + ` AND DATE(order_items.created_at) ='` + day + `'` + sql1a
		database.DB.Raw(sql1).Scan(&orderitems1)
	} else if month != "" {
		sql = sql + ` AND EXTRACT(MONTH FROM created_at) = '` + month + `'` + sqla
		database.DB.Raw(sql).Scan(&orderitems)
		sql1 = sql1 + ` AND EXTRACT(MONTH FROM order_items.created_at) = '` + month + `'` + sql1a
		database.DB.Raw(sql1).Scan(&orderitems1)
	} else if year != "" {
		sql = sql + ` AND EXTRACT(YEAR FROM created_at) = '` + year + `'` + sqla
		database.DB.Raw(sql).Scan(&orderitems)
		sql1 = sql1 + ` AND EXTRACT(YEAR FROM order_items.created_at) = '` + year + `'` + sql1a
		database.DB.Raw(sql1).Scan(&orderitems1)
	} else {
		database.DB.Raw(sql).Scan(&orderitems)
		database.DB.Raw(sql1).Scan(&orderitems1)
		//database.DB.Raw(sql + sqla).Scan(&orderitems)
	}
	// database.DB.Raw(sql).Scan(&orderitems)
	// database.DB.Raw(sql1).Scan(&orderitems1)
	//var product []string
	var product1 responsemodels.Product
	var products []responsemodels.Product
	//fmt.Println("---------", orderitems)
	for _, v := range orderitems {
		// var a string
		fmt.Println("product id count", v.ProductID)
		// database.DB.Raw(`select product_name from products where id = ?`, v.ProductID).Scan(&a)

		database.DB.Raw("SELECT * FROM products join categories on products.category_id=categories.id WHERE products.id = ?", v.ProductID).Scan(&product1)

		products = append(products, product1)
		//product = append(product, a)

	}
	var category []string
	for _, v := range orderitems1 {
		fmt.Println("hi---hello", v.Count)
		category = append(category, v.CategoryName)

	}
	c.JSON(http.StatusOK, gin.H{
		"best_selling_product":  products,
		"best_selling_category": category,
		"message":               "retrieved best selling products",
	})

}

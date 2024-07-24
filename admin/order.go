package admin

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/helper"
	"github.com/Ansalps/GeZOne/models"
	"github.com/Ansalps/GeZOne/responsemodels"
	"github.com/gin-gonic/gin"
)

func OrderList(c *gin.Context) {
	var orders []responsemodels.Order
	var address responsemodels.Address
	database.DB.Raw(`SELECT orders.id,orders.created_at,orders.updated_at,orders.deleted_at,orders.user_id,orders.address_id,orders.total_amount,orders.order_status,addresses.created_at,addresses.updated_at,addresses.deleted_at,addresses.user_id,addresses.country,addresses.state,addresses.street_name,addresses.district,addresses.pin_code,addresses.phone,addresses.default
	FROM orders
	JOIN addresses ON orders.address_id = addresses.id
	order by orders.id`).Scan(&orders)
	//database.DB.Find(&Orders)
	for i, v := range orders {
		database.DB.Raw(`SELECT *
	        FROM orders
	        JOIN addresses ON orders.address_id = addresses.id
	        WHERE orders.id = ?`, v.ID).Scan(&address)
		orders[i].Address = address
	}
	c.JSON(http.StatusOK, gin.H{
		"Order": orders,
	})

}

func ChangeOrderStatus(c *gin.Context) {
	orderID := c.Param("id")
	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM orders where id = ?`, orderID).Scan(&count)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "order id does not exist",
		})
		return
	}
	var Order models.CancelOrder
	err := c.BindJSON(&Order)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "false",
			"message": "failed to bind request",
		})
	}
	if err := helper.Validate(Order); err != nil {
		fmt.Println("", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"message":    err.Error(),
			"error_code": http.StatusBadRequest,
		})
		return
	}
	p := Order.OrderStatus
	if p != "delivered" && p != "cancelled" && p != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid order status",
		})
		return
	}
	if Order.OrderStatus == "delivered" {
		fmt.Println("update in payments table")
		//var payment models.Payments
		now := time.Now()
		today := now.Format("2006-01-02")

		payment := models.Payments{
			PaymentDate:   today,
			PaymentStatus: "paid",
		}
		fmt.Println("hi")
		database.DB.Model(&models.Payments{}).Where("order_id = ?", orderID).Updates(&payment)
		fmt.Println("hello")
		var OrderItems []models.OrderItems
		database.DB.Where("order_id = ?", orderID).Find(&OrderItems)

		for _, v := range OrderItems {
			var stock uint
			database.DB.Model(&models.Product{}).Where("id = ?", v.ProductID).Pluck("stock", &stock)
			fmt.Println("stock first", stock)
			stock = stock - v.Qty
			fmt.Println("stock : v.product_id v.qty", stock, v.ProductID, v.Qty)
			database.DB.Model(&models.Product{}).Where("id = ?", v.ProductID).Update("stock", stock)
		}
	}
	order := models.Order{
		OrderStatus: Order.OrderStatus,
	}
	database.DB.Model(&models.Order{}).Where("id = ?", orderID).Updates(&order)
	c.JSON(http.StatusOK, gin.H{
		"meassage": "order status changed successfully",
	})
}

package admin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/helper"
	"github.com/Ansalps/GeZOne/models"
	requestmodels "github.com/Ansalps/GeZOne/requestmodels"
	"github.com/Ansalps/GeZOne/responsemodels"
	"github.com/gin-gonic/gin"
)

func OrderList(c *gin.Context) {
	listorder := c.Query("list_order")
	var orders []responsemodels.Order
	var address responsemodels.Address
	sql := `SELECT orders.id,
		orders.created_at,
		orders.updated_at,
		orders.user_id,
		orders.address_id,
		orders.total_amount,
		orders.payment_method,
		payments.payment_status AS payment_status,
		orders.order_status,
		orders.delivery_status,
		orders.offer_discount,
		orders.coupon_discount,
		orders.coupon_code,
		orders.total_discount_amount,
		orders.final_amount
	FROM orders
	JOIN payments ON orders.id = payments.order_id`

	switch listorder {
	case "":
		sql += ` ORDER BY orders.id DESC`
	case "ASC":
		sql += ` ORDER BY orders.id ASC`
	case "DSC":
		sql += ` ORDER BY orders.id DESC`
	}

	database.DB.Raw(sql).Scan(&orders)
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
func OrderItemsList(c *gin.Context) {
	orderId := c.Param("order_id")
	listorder := c.Query("list_order")
	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM orders where id = ?`, orderId).Scan(&count)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "order id does not exist",
		})
		return
	}
	var orderitems []responsemodels.OrderItems

	sql := `SELECT order_items.id,order_items.created_at,order_items.updated_at,order_items.order_id,order_items.product_id,order_items.product_variant_id,order_items.size,products.product_name,order_items.quantity,order_items.unit_price AS price,order_items.total_amount,order_items.order_status,order_items.delivery_status,order_items.payment_method,order_items.coupon_discount,order_items.offer_discount,order_items.total_discount,order_items.final_amount FROM order_items join products on order_items.product_id=products.id WHERE order_items.order_id = ?`
	switch listorder {
	case "":
		sql += ` ORDER BY order_items.id DESC`
	case "ASC":
		sql += ` ORDER BY order_items.id ASC`
	case "DSC":
		sql += ` ORDER BY order_items.id DESC`
	}

	database.DB.Raw(sql, orderId).Scan(&orderitems)
	c.JSON(http.StatusOK, gin.H{
		"order_items": orderitems,
	})

}

func canUpdateDeliveryStatus(currentStatus, nextStatus string) bool {
	current := strings.ToLower(strings.TrimSpace(currentStatus))
	next := strings.ToLower(strings.TrimSpace(nextStatus))

	if current == "" || next == "" || current == next {
		return false
	}

	switch current {
	case "pending":
		return next == "shipped" || next == "failed"
	case "shipped":
		return next == "delivered" || next == "failed"
	default:
		return false
	}
}

func canUpdatePaymentStatus(currentStatus, nextStatus string) bool {
	current := strings.ToLower(strings.TrimSpace(currentStatus))
	next := strings.ToLower(strings.TrimSpace(nextStatus))

	if current == "" || next == "" || current == next {
		return false
	}

	switch current {
	case "pending":
		return next == "paid" || next == "failed"
	default:
		return false
	}
}

func shouldCompleteOrderOnDelivery(paymentMethod string) bool {
	return !strings.EqualFold(strings.TrimSpace(paymentMethod), "COD")
}

func ChangeDeliveryStatus(c *gin.Context) {
	orderID := c.Param("id")
	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM orders where id = ?`, orderID).Scan(&count)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "order id does not exist",
		})
		return
	}

	var payload requestmodels.ChangeDeliveryStatus
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "failed to bind request",
		})
		return
	}
	if err := helper.Validate(payload); err != nil {
		fmt.Println("", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"message":    err.Error(),
			"error_code": http.StatusBadRequest,
		})
		return
	}

	nextStatus := strings.ToLower(strings.TrimSpace(payload.DeliveryStatus))
	if nextStatus != "pending" && nextStatus != "shipped" && nextStatus != "delivered" && nextStatus != "failed" && nextStatus != "returned" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid delivery status",
		})
		return
	}

	var currentStatus string
	database.DB.Model(&models.Order{}).Where("id = ?", orderID).Pluck("delivery_status", &currentStatus)
	currentStatus = strings.ToLower(strings.TrimSpace(currentStatus))

	if currentStatus == "delivered" || currentStatus == "failed" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "delivery status cannot be updated from the current state",
		})
		return
	}

	if !canUpdateDeliveryStatus(currentStatus, nextStatus) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid delivery status transition",
		})
		return
	}

	orderUpdate := models.Order{DeliveryStatus: nextStatus}
	database.DB.Model(&models.Order{}).Where("id = ?", orderID).Updates(&orderUpdate)
	database.DB.Model(&models.OrderItems{}).Where("order_id = ?", orderID).Update("delivery_status", nextStatus)

	if nextStatus == "delivered" {
		var paymentMethod string
		database.DB.Model(&models.Order{}).Where("id = ?", orderID).Pluck("payment_method", &paymentMethod)
		if shouldCompleteOrderOnDelivery(paymentMethod) {
			database.DB.Model(&models.Order{}).Where("id = ?", orderID).Update("order_status", "completed")
			database.DB.Model(&models.OrderItems{}).Where("order_id = ?", orderID).Update("order_status", "completed")
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "delivery status changed successfully",
	})
}

func ChangePaymentStatus(c *gin.Context) {
	orderID := c.Param("id")
	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM orders where id = ?`, orderID).Scan(&count)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "order id does not exist",
		})
		return
	}

	var payload requestmodels.ChangePaymentStatus
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "failed to bind request",
		})
		return
	}
	if err := helper.Validate(payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"message":    err.Error(),
			"error_code": http.StatusBadRequest,
		})
		return
	}

	nextStatus := strings.ToLower(strings.TrimSpace(payload.PaymentStatus))
	if nextStatus != "pending" && nextStatus != "paid" && nextStatus != "failed" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid payment status",
		})
		return
	}

	var paymentMethod string
	database.DB.Model(&models.Order{}).Where("id = ?", orderID).Pluck("payment_method", &paymentMethod)
	if !strings.EqualFold(paymentMethod, "COD") {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "payment status update is only allowed for COD orders",
		})
		return
	}

	var currentStatus string
	database.DB.Model(&models.Payments{}).Where("order_id = ?", orderID).Order("created_at DESC").Limit(1).Pluck("payment_status", &currentStatus)
	currentStatus = strings.ToLower(strings.TrimSpace(currentStatus))

	if currentStatus == "paid" || currentStatus == "failed" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "payment status cannot be updated from the current state",
		})
		return
	}

	if !canUpdatePaymentStatus(currentStatus, nextStatus) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid payment status transition",
		})
		return
	}

	database.DB.Model(&models.Payments{}).Where("order_id = ?", orderID).Order("created_at DESC").Limit(1).Update("payment_status", nextStatus)
	if nextStatus == "paid" {
		database.DB.Model(&models.Order{}).Where("id = ?", orderID).Update("order_status", "completed")
		database.DB.Model(&models.OrderItems{}).Where("order_id = ?", orderID).Update("order_status", "completed")
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "payment status changed successfully",
	})
}

func ChangeOrderStatus(c *gin.Context) {
	ChangeDeliveryStatus(c)
}

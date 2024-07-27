package user

import (
	"fmt"
	"net/http"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/helper"
	"github.com/Ansalps/GeZOne/middleware"
	"github.com/Ansalps/GeZOne/models"
	"github.com/Ansalps/GeZOne/responsemodels"
	"github.com/gin-gonic/gin"
)

func Order(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Claims not found"})
		return
	}

	customClaims, ok := claims.(*middleware.CustomClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
		return
	}

	userID := customClaims.ID
	//addressid verifying
	var OrderAdd models.OrderAdd
	err := c.BindJSON(&OrderAdd)
	response := gin.H{
		"status":  false,
		"message": "failed to bind request",
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// Validate the content of the JSON
	if err := helper.Validate(OrderAdd); err != nil {
		fmt.Println("", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"message":    err.Error(),
			"error_code": http.StatusBadRequest,
		})
		return
	}
	var count1 int64
	database.DB.Raw(`SELECT COUNT(*) FROM addresses where id = ? AND user_id = ? AND deleted_at IS NULL`, OrderAdd.AddressID, userID).Scan(&count1)
	if count1 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "address_id does not exist for this particular user",
		})
		return
	}
	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM cart_items WHERE user_id=? and deleted_at IS NULL`, userID).Scan(&count)
	fmt.Println("count ", count)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "false",
			"message": "cart empty, order can't be placed",
		})
		return
	}
	var totalquantity uint
	database.DB.Raw(`SELECT SUM(qty) FROM cart_items WHERE user_id=? and deleted_at IS NULL`, userID).Scan(&totalquantity)
	fmt.Println("total quantity", totalquantity)
	if totalquantity == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "false",
			"message": "cart empty, order can't be placed",
		})
		return
	}

	var totalamount float64
	database.DB.Raw("SELECT SUM(total_amount) from cart_items where user_id = ? and deleted_at IS NULL", userID).Scan(&totalamount)
	order := models.Order{
		UserID:      userID,
		AddressID:   OrderAdd.AddressID,
		TotalAmount: totalamount,
	}
	database.DB.Create(&order)
	var CartItems []models.CartItems
	database.DB.Where("user_id = ?", userID).Find(&CartItems)

	var ID uint
	//database.DB.Model(&models.Order{}).Where("user_id = ?", userID).Pluck("id", &ID)
	database.DB.Raw(`SELECT id FROM orders where user_id = ? ORDER BY created_at DESC LIMIT 1`, userID).Scan(&ID)
	fmt.Println("latest order id ", ID)
	//var orderItem models.OrderItems
	for _, v := range CartItems {
		//var Product models.Product
		//database.DB.Where("id = ?", v.ProductID).First(&Product)
		//database.DB.Where("price=?",v.)
		fmt.Println("qty", v.Qty)
		if v.Qty == 0 {
			continue
		}
		for i := 0; i < int(v.Qty); i++ {
			var price float64
			database.DB.Model(&models.Product{}).Where("id = ?", v.ProductID).Pluck("price", &price)
			fmt.Println("id", ID)
			orderItem := models.OrderItems{
				OrderID:   ID,
				ProductID: v.ProductID,
				//Qty:         v.Qty,
				Price: price,
				//TotalAmount: float64(v.Qty) * price,
			}
			fmt.Println("order id", orderItem.OrderID)
			fmt.Println("order item create hi")
			database.DB.Create(&orderItem)
			fmt.Println("order item create hello")
		}

	}
	//clearing cart
	//var cart models.CartItems
	//database.DB.Exec("DELETE FROM cart_items where user_id=?", userID).Scan(&cart)

	//database.DB.Create(&orderItem)
	//var Payment models.Payments
	Payment := models.Payments{
		UserID:      userID,
		OrderID:     order.ID,
		TotalAmount: totalamount,
	}
	database.DB.Create(&Payment)
	database.DB.Where("user_id = ?", userID).Delete(&models.CartItems{})
	var order1 responsemodels.Order
	var address responsemodels.Address
	var orderitems1 []responsemodels.OrderItems
	database.DB.Raw(`SELECT orders.id,orders.created_at,orders.updated_at,orders.deleted_at,orders.user_id,orders.address_id,orders.total_amount,orders.order_status FROM orders join addresses on orders.address_id=addresses.id WHERE orders.user_id = ? ORDER BY orders.created_at desc LIMIT 1`, userID).Scan(&order1)
	fmt.Println("-----------------")
	fmt.Println("user id ", userID)
	var orderid uint
	database.DB.Raw(`SELECT id FROM orders WHERE user_id = ? ORDER BY created_at desc limit 1`, userID).Scan(&orderid)
	fmt.Println("order id ", orderid)
	var addressid uint
	database.DB.Raw(`SELECT address_id FROM orders WHERE user_id = ? ORDER BY created_at desc limit 1`, userID).Scan(&addressid)
	fmt.Println("address id", addressid)
	database.DB.Raw(`SELECT * FROM addresses WHERE id = ?`, addressid).Scan(&address)
	order1.Address = address
	database.DB.Raw(`SELECT order_items.id,order_items.created_at,order_items.updated_at,order_items.deleted_at,order_items.order_id,order_items.product_id,products.product_name,order_items.price,order_items.order_status FROM order_items join products on order_items.product_id=products.id WHERE order_items.order_id = ? ORDER BY order_items.id`, orderid).Scan(&orderitems1)
	c.JSON(http.StatusOK, gin.H{"message": "Order added successfully",
		"Order":       order1,
		"Order items": orderitems1})
}

//Address selection, clearing cart, updating payments table

package user

import (
	"fmt"
	"net/http"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/helper"
	"github.com/Ansalps/GeZOne/middleware"
	"github.com/Ansalps/GeZOne/models"
	"github.com/gin-gonic/gin"
)

func Order(c *gin.Context) {
	//userID := c.Param("user_id")
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
	fmt.Println("print user id : ", userID)
	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM cart_items WHERE user_id=? and deleted_at IS NULL`, userID).Scan(&count)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "false",
			"message": "cart empty, order can't be placed",
		})
		return
	}
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
	fmt.Println("address id ", OrderAdd.AddressID)
	fmt.Println("user id ", userID)
	var count1 int64
	database.DB.Raw(`SELECT COUNT(*) FROM addresses where id = ? AND user_id = ? AND deleted_at IS NULL`, OrderAdd.AddressID, userID).Scan(&count1)
	if count1 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "address_id does not exist for this particular user",
		})
		return
	}
	// var address models.Address
	// tx := database.DB.Where("id = ?", OrderAdd.AddressID).Find(&address)
	// if tx.Error != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"message": "Address id does not exist",
	// 	})
	// }
	var totalamount float64
	//database.DB.Model(&models.CartItems{}).Where("user_id = ?", userID).Pluck("total_amount", &totalamount)
	fmt.Println("hi order")
	database.DB.Raw("SELECT SUM(total_amount) from cart_items where user_id = ? and deleted_at IS NULL", userID).Scan(&totalamount)
	fmt.Println("hello order")
	order := models.Order{
		UserID:      userID,
		AddressID:   OrderAdd.AddressID,
		TotalAmount: totalamount,
	}
	database.DB.Create(&order)
	var CartItems []models.CartItems
	database.DB.Where("user_id = ?", userID).Find(&CartItems)

	var ID string
	database.DB.Model(&models.Order{}).Where("user_id = ?", userID).Pluck("id", &ID)
	var orderItem models.OrderItems
	for _, v := range CartItems {
		//var Product models.Product
		//database.DB.Where("id = ?", v.ProductID).First(&Product)
		//database.DB.Where("price=?",v.)
		if v.Qty == 0 {
			continue
		}
		var price float64
		database.DB.Model(&models.Product{}).Where("id = ?", v.ProductID).Pluck("price", &price)
		orderItem = models.OrderItems{
			OrderID:     ID,
			ProductID:   v.ProductID,
			Qty:         v.Qty,
			Price:       price,
			TotalAmount: float64(v.Qty) * price,
		}
		database.DB.Create(&orderItem)
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
	c.JSON(http.StatusOK, gin.H{"message": "Order added successfully"})
}

//Address selection, clearing cart, updating payments table

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

func Profile(c *gin.Context) {
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
	var User []responsemodels.User
	database.DB.Where("id = ?", userID).First(&User)
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "successfully retrieved user informations",
		"data": gin.H{
			"User": User,
		},
	})
}

func ProfileEdit(c *gin.Context) {
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
	var Profile models.ProfileEdit
	err := c.BindJSON(&Profile)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "false",
			"message": "failed to bind request",
		})
		return
	}
	// Validate the content of the JSON
	if err := helper.Validate(Profile); err != nil {
		fmt.Println("", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"message":    err.Error(),
			"error_code": http.StatusBadRequest,
		})
		return
	}
	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM users where phone = ?`, Profile.Phone).Scan(&count)
	if count != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "already registered mobile",
		})
		return
	}
	user := models.User{
		FirstName: Profile.FirstName,
		LastName:  Profile.LastName,
		Phone:     Profile.Phone,
	}
	//database.DB.Where("id = ?", userID).Updates(&user)
	database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(&user)
	c.JSON(http.StatusOK, gin.H{"message": "updated user profile"})
}

func PasswordChange(c *gin.Context) {
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
	var PasswordChange models.PasswordChange
	err := c.BindJSON(&PasswordChange)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "false",
			"message": "failed to bind request",
		})
		return
	}
	// Validate the content of the JSON
	if err := helper.Validate(PasswordChange); err != nil {
		fmt.Println("", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"message":    err.Error(),
			"error_code": http.StatusBadRequest,
		})
		return
	}
	if PasswordChange.Password != PasswordChange.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"message":    "password should match",
			"error_code": http.StatusBadRequest,
		})
		return
	}
	passwordchange := models.User{
		Password: PasswordChange.Password,
	}
	database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(&passwordchange)
	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

func AddressList(c *gin.Context) {
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
	var Address []responsemodels.Address
	database.DB.Where("user_id = ? and deleted_at IS NULL", userID).Find(&Address)
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "successfully retrieved user informations",
		"data": gin.H{
			"Address": Address,
		},
	})
}

func OrderList(c *gin.Context) {
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
	var orders []responsemodels.Order
	var address responsemodels.Address
	//var Order responsemodels.Order
	//database.DB.Where("user_id = ?", userID).Find(&Order)
	// database.DB.Raw(`SELECT * FROM orders join addresses on orders.address_id = addresses.id where orders.user_id = ?`, userID).Scan(&Order)

	database.DB.Raw(`SELECT orders.id,orders.created_at,orders.updated_at,orders.deleted_at,orders.user_id,orders.address_id,orders.total_amount,orders.order_status,addresses.created_at,addresses.updated_at,addresses.deleted_at,addresses.user_id,addresses.country,addresses.state,addresses.street_name,addresses.district,addresses.pin_code,addresses.phone,addresses.default
	FROM orders
	JOIN addresses ON orders.address_id = addresses.id
	WHERE orders.user_id = ? order by orders.id`, userID).Scan(&orders)
	// database.DB.Raw(`SELECT *
	//         FROM orders
	//         JOIN addresses ON orders.address_id = addresses.id
	//         WHERE orders.user_id = ?`, userID).Scan(&Order)

	for i, v := range orders {
		database.DB.Raw(`SELECT *
	        FROM orders
	        JOIN addresses ON orders.address_id = addresses.id
	        WHERE orders.user_id = ? AND orders.id = ?`, userID, v.ID).Scan(&address)
		orders[i].Address = address
	}

	// query.Find(&Address)
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "successfully retrieved user informations",
		"data": gin.H{
			//"Address": Address,
			"Order": orders,
		},
	})
}

// func OrderList(c *gin.Context) {
// 	// userID := c.Param("userID")
// 	claims, exists := c.Get("claims")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Claims not found"})
// 		return
// 	}

// 	customClaims, ok := claims.(*helper.CustomClaims)
// 	if !ok {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
// 		return
// 	}

// 	userID := customClaims.ID
// 	fmt.Println("print user id : ", userID)

// 	var orders []responsemodels.Order

// 	query := `
//         SELECT *
//         FROM orders
//         JOIN addresses ON orders.address_id = addresses.id
//         WHERE orders.user_id = ?`

// 	rows, err := database.DB.Raw(query, userID).Rows()
// 	if err != nil {
// 		fmt.Println("Error fetching data:", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to retrieve data"})
// 		return
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var order responsemodels.Order
// 		var address responsemodels.Address

// 		err = rows.Scan(
// 			&order.ID, &order.CreatedAt, &order.UpdatedAt, &order.DeletedAt, &order.UserID, &order.AddressID,
// 			&order.TotalAmount, &order.OrderStatus,
// 			&address.ID, &address.CreatedAt, &address.UpdatedAt, &address.DeletedAt, &address.UserID, &address.Country,
// 			&address.State, &address.District, &address.StreetName, &address.PinCode, &address.Phone, &address.Default,
// 		)
// 		if err != nil {
// 			fmt.Println("Error scanning data:", err)
// 			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to scan data"})
// 			return
// 		}
// 		order.Address = address
// 		orders = append(orders, order)
// 	}

//		c.JSON(http.StatusOK, gin.H{
//			"status":  true,
//			"message": "successfully retrieved user informations",
//			"data": gin.H{
//				"Order": orders,
//			},
//		})
//	}
func OrderItemsList(c *gin.Context) {
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
	orderId := c.Param("order_id")
	//var order models.Order
	// t := database.DB.Where("id = ?", orderId).Find(&order)
	// if t.Error != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"message": "order id does not exist",
	// 	})
	// }
	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM orders where id = ? AND user_id = ?`, orderId, userID).Scan(&count)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "order id does not exist for this particular user",
		})
		return
	}
	var orderitems []responsemodels.OrderItems
	//tx := database.DB.Where("order_id = ?", orderId).Find(&orderitems)
	tx := database.DB.Raw(`SELECT order_items.id,order_items.created_at,order_items.updated_at,order_items.deleted_at,order_items.order_id,order_items.product_id,products.product_name,order_items.price,order_items.order_status FROM order_items join products on order_items.product_id=products.id WHERE order_items.order_id = ? ORDER BY order_items.id`, orderId).Scan(&orderitems)
	if tx.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "order_id does not exist",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"order items": orderitems,
	})
}

func CancelOrder(c *gin.Context) {
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
	//userID:=c.Param("user_id")
	orderID := c.Param("order_id")
	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM orders where id = ? AND user_id = ?`, orderID, userID).Scan(&count)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "order id does not exist for this particular user",
		})
		return
	}
	//var Order models.Order
	// tx := database.DB.Where("id = ?", orderID).First(&Order)
	// if tx.Error != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"message": "Order_id does not exist in database",
	// 	})
	// 	return
	// }
	// err := c.BindJSON(&CancelOrder)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"status":  "false",
	// 		"message": "failed to bind request",
	// 	})
	// 	return
	// }
	// Validate the content of the JSON
	// if err := helper.Validate(CancelOrder); err != nil {
	// 	fmt.Println("", err)
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"status":     false,
	// 		"message":    err.Error(),
	// 		"error_code": http.StatusBadRequest,
	// 	})
	// 	return
	// }
	// cancelorder := models.Order{
	// 	OrderStatus: CancelOrder.OrderStatus,
	// }
	var orderstatus string
	database.DB.Model(&models.Order{}).Where("id = ?", orderID).Pluck("order_status", &orderstatus)
	if orderstatus == "cancelled" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cannot cancel the order item, item is already cancelled bu user or admin",
		})
		return
	}
	if orderstatus == "shipped" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cannot cancel the order, items are already shipped",
		})
		return
	}
	if orderstatus == "delivered" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cannot cancel the order, items are delivered",
		})
		return
	}
	if orderstatus == "failed" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "your order is already failed due to some inventory issues",
		})
		return
	}
	database.DB.Model(&models.Order{}).Where("id = ?", orderID).Update("order_status", "cancelled")
	database.DB.Model(&models.OrderItems{}).Where("order_id = ?", orderID).Update("order_status", "cancelled")
	//database.DB.Raw(`UPDATE order_items SET order_status = 'cancelled' WHERE `)
	c.JSON(http.StatusOK, gin.H{"message": "Order cancelled successfully"})
}

func CancelSingleOrderItem(c *gin.Context) {
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
	ItemID := c.Param("orderitem_id")
	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM order_items where id = ?`, ItemID).Scan(&count)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "order item id does not exist for this particular user",
		})
		return
	}
	var orderstatus string
	database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Pluck("order_status", &orderstatus)
	fmt.Println("order status", orderstatus)
	if orderstatus == "cancelled" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cannot cancel the order item, item is already cancelled",
		})
		return
	}
	if orderstatus == "shipped" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cannot cancel the order item, item is already shipped",
		})
		return
	}
	if orderstatus == "delivered" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cannot cancel the order item, item is delivered",
		})
		return
	}
	if orderstatus == "failed" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "your order is already failed due to some inventory issues",
		})
		return
	}

	var price float64
	database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Pluck("price", &price)
	var orderid uint
	database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Pluck("order_id", &orderid)
	fmt.Println("price :", price)
	var totalamount float64
	database.DB.Model(&models.Order{}).Where("id = ?", orderid).Pluck("total_amount", &totalamount)
	totalamount = totalamount - price
	database.DB.Model(&models.Order{}).Where("id = ?", orderid).Update("total_amount", totalamount)
	database.DB.Model(&models.Payments{}).Where("order_id = ?", orderid).Update("total_amount", totalamount)
	var count1 int
	database.DB.Raw(`SELECT COUNT(*) FROM order_items WHERE order_id = ? AND order_status != 'cancelled'`, orderid).Scan(&count1)
	fmt.Println("count1", count1)
	if count1 == 1 {
		database.DB.Model(&models.Order{}).Where("id = ?", orderid).Update("order_status", "cancelled")
	}
	database.DB.Model(&models.OrderItems{}).Where("id = ?", ItemID).Update("order_status", "cancelled")
	c.JSON(http.StatusOK, gin.H{"message": "Order item cancelled successfully"})
}

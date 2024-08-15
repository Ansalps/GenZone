package admin

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/models"
	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf/v2"
)

func GenerateInvoice(c *gin.Context) {
	OrderID := c.Param("order_id")
	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM orders WHERE id = ? AND order_status!='cancelled'`, OrderID).Scan(&count)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "order is cancelled or order id does not exist",
		})
		return
	}

	database.DB.Exec(`TRUNCATE invoices`)

	var userid uint
	database.DB.Model(&models.Order{}).Where("id = ?", OrderID).Pluck("user_id", &userid)

	var firstname string
	database.DB.Model(&models.User{}).Where("id = ?", userid).Pluck("first_name", &firstname)

	var lastname string
	database.DB.Model(&models.User{}).Where("id = ?", userid).Pluck("last_name", &lastname)

	fullname := firstname + " " + lastname

	var orderdate time.Time
	database.DB.Raw(`SELECT created_at from orders where id = ?`, OrderID).Scan(&orderdate)

	var OrderStatus string
	database.DB.Raw(`SELECT order_status from orders where id = ?`, OrderID).Scan(&OrderStatus)

	var paymentmethod string
	database.DB.Raw(`SELECT payment_method from orders where id = ?`, OrderID).Scan(&paymentmethod)

	var paymentstatus string

	if paymentmethod == "COD" {
		if OrderStatus == "delivered" {
			paymentstatus = "paid"
		} else {
			paymentstatus = "pending"
		}
	} else {
		paymentstatus = "paid"
	}

	// var orderaddress []models.Order
	// database.DB.Raw(`SELECT * from orders where id = ?`, OrderID).Scan(&orderaddress)
	// fmt.Println("kkkkkkk", orderaddress)
	// var Country string
	// var State string
	// var District string
	// var StreetName string
	// var PinCode string
	// var Phone string
	// database.DB.Preload(Addre)
	// for _, v := range orderaddress {
	// 	Country = v.Address.Country
	// 	State = v.Address.State
	// 	District = v.Address.District
	// 	StreetName = v.Address.StreetName
	// 	PinCode = v.Address.PinCode
	// 	Phone = v.Address.Phone
	// }

	// var order models.Order
	// database.DB.Where("id =?",OrderID).First(&order)

	var order models.Order
	database.DB.Preload("Address").First(&order, OrderID)

	Country := order.Address.Country
	State := order.Address.State
	District := order.Address.District
	StreetName := order.Address.StreetName
	PinCode := order.Address.PinCode
	Phone := order.Address.Phone
	// var OrderStatus string
	// database.DB
	var orderitems []models.OrderItems
	database.DB.Raw(`SELECT * FROM order_items WHERE order_status!='return' AND order_status!='cancelled' AND order_id = ? order by product_id`, OrderID).Scan(&orderitems)

	var productid string
	fmt.Println("highest=====")
	for i, v := range orderitems {
		fmt.Println("is it enteing in rangee")
		fmt.Println("first productid ", productid)
		if productid == v.ProductID {
			fmt.Println("is it entering here?")
			var qty uint
			database.DB.Model(&models.Invoice{}).Where("product_id = ?", v.ProductID).Pluck("quantity", &qty)
			qty = qty + 1
			var mrp float64
			database.DB.Model(&models.Invoice{}).Where("product_id = ?", v.ProductID).Pluck("mrp", &mrp)
			mrp = mrp + v.Price
			var coupondiscount float64
			database.DB.Model(&models.Invoice{}).Where("product_id = ?", v.ProductID).Pluck("coupon_discount", &coupondiscount)
			coupondiscount = coupondiscount + v.CouponDiscount
			coupondiscount = math.Round(coupondiscount*100) / 100
			var offerdiscount float64
			database.DB.Model(&models.Invoice{}).Where("product_id = ?", v.ProductID).Pluck("offer_discount", &offerdiscount)
			offerdiscount = offerdiscount + v.OfferDiscount
			var totaldiscount float64
			database.DB.Model(&models.Invoice{}).Where("product_id = ?", v.ProductID).Pluck("total_discount", &totaldiscount)
			totaldiscount = totaldiscount + v.TotalDiscount
			var finalprce float64
			database.DB.Model(&models.Invoice{}).Where("product_id = ?", v.ProductID).Pluck("final_price", &finalprce)
			finalprce = finalprce + finalprce

			invoice := models.Invoice{
				Quantity:       qty,
				MRP:            mrp,
				CouponDiscount: coupondiscount,
				OfferDiscount:  offerdiscount,
				TotalDiscount:  totaldiscount,
				FinalPrice:     finalprce,
			}
			database.DB.Model(&models.Invoice{}).Where("product_id = ?", productid).Updates(&invoice)
			continue
		}
		productid = v.ProductID
		fmt.Println("second productid ", productid)
		CouponDiscount1 := math.Round(v.CouponDiscount*100) / 100
		fmt.Println("======", CouponDiscount1)
		var orderitem1 models.OrderItems
		database.DB.Preload("Product").First(&orderitem1, v.ID)
		fmt.Println("what is printed here ?", orderitem1.Product.ProductName)
		invoice := models.Invoice{
			No:             i + 1,
			ProductID:      v.ProductID,
			ProductName:    orderitem1.Product.ProductName,
			Quantity:       1,
			MRP:            v.Price,
			CouponDiscount: CouponDiscount1,
			OfferDiscount:  v.OfferDiscount,
			TotalDiscount:  v.TotalDiscount,
			FinalPrice:     v.Price - v.TotalDiscount,
		}
		database.DB.Create(&invoice)
	}
	fmt.Println("loweat===")
	var invoiceitem []models.Invoice
	database.DB.Find(&invoiceitem)
	fmt.Println("jjjjj", invoiceitem)
	var grandtotal float64
	database.DB.Raw(`SELECT SUM(final_price) from invoices`).Scan(&grandtotal)
	fmt.Println("from here===")

	pdf := gofpdf.New("P", "mm", "A2", "")

	// Add a new page
	pdf.AddPage()

	// Set the title
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Tax Invoice")
	pdf.Ln(12)

	// Set the title
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "GenZone")
	pdf.Ln(12)

	// Set the title
	// pdf.SetFont("Arial", "B", 16)
	// pdf.Cell(0, 10, OrderID)
	// pdf.Ln(12)
	// // Set the title
	// pdf.SetFont("Arial", "B", 16)
	orderDateString := orderdate.Format("2006-01-02 15:04:05")
	// pdf.Cell(0, 10, orderDateString)
	// pdf.Ln(12)
	// // Set the title
	// pdf.SetFont("Arial", "B", 16)
	// pdf.Cell(0, 10, paymentstatus)
	// pdf.Ln(12)

	// Add summary information
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, fmt.Sprintf("Order ID: %s", OrderID))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Order Date: %s", orderDateString))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Patment Status: %s", paymentstatus))
	pdf.Ln(8)

	// Set the title
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Address")
	pdf.Ln(12)
	// Add summary information
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, fmt.Sprintf("Name: %s", fullname))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Address: %s,%s,%s,%s-%s", StreetName, District, State, Country, PinCode))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Phone: %s", Phone))
	pdf.Ln(8)

	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(240, 240, 240) // Light grey background
	headers := []string{"No", "Product ID", "Product Name", "Quantity", "MRP", "Coupon Discount", "Offer Discount", "Total Discount", "Final Price"}

	widths := []float64{15, 20, 30, 20, 25, 35, 35, 35, 35} // Adjust the widths based on content
	for i, header := range headers {
		pdf.CellFormat(widths[i], 10, header, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// Set data with alternating row colors
	pdf.SetFont("Arial", "", 10)
	fill := false
	for i, item := range invoiceitem {
		if fill {
			pdf.SetFillColor(230, 230, 230) // Slightly darker grey for alternating rows
		} else {
			pdf.SetFillColor(255, 255, 255) // White for other rows
		}
		fill = !fill
		fmt.Println("hirrrrrhello ", i)
		pdf.CellFormat(widths[0], 10, strconv.Itoa(int(item.No)), "1", 0, "C", true, 0, "")
		pdf.CellFormat(widths[1], 10, item.ProductID, "1", 0, "C", true, 0, "")
		pdf.CellFormat(widths[2], 10, item.ProductName, "1", 0, "C", true, 0, "")
		pdf.CellFormat(widths[3], 10, strconv.Itoa(int(item.Quantity)), "1", 0, "C", true, 0, "")
		pdf.CellFormat(widths[4], 10, fmt.Sprintf("%.2f", item.MRP), "1", 0, "C", true, 0, "")
		pdf.CellFormat(widths[5], 10, fmt.Sprintf("%.2f", item.CouponDiscount), "1", 0, "C", true, 0, "")
		pdf.CellFormat(widths[6], 10, fmt.Sprintf("%.2f", item.OfferDiscount), "1", 0, "C", true, 0, "")
		pdf.CellFormat(widths[7], 10, fmt.Sprintf("%.2f", item.TotalDiscount), "1", 0, "C", true, 0, "")
		pdf.CellFormat(widths[8], 10, fmt.Sprintf("%.2f", item.FinalPrice), "1", 0, "C", true, 0, "")
		pdf.Ln(-1)
		fmt.Println("end")
	}

	grandtotalString := strconv.FormatFloat(grandtotal, 'f', 2, 64)
	pdf.Cell(0, 10, fmt.Sprintf("Grand Total: %s", grandtotalString))
	pdf.Ln(8)

	// Footer with the current date
	pdf.SetY(-15)
	pdf.SetFont("Arial", "I", 8)
	pdf.Cell(0, 10, fmt.Sprintf("Generated on %s", time.Now().Format("2006-01-02")))

	pdf.OutputFileAndClose("invoice.pdf")

	c.File("invoice.pdf")

	// c.JSON(http.StatusOK, gin.H{
	// 	"title":          "Tax Invoice",
	// 	"buisness_name":  "GenZOne",
	// 	"order_id":       OrderID,
	// 	"order_date":     orderdate,
	// 	"payment_status": paymentstatus,
	// 	"address":        "Address",
	// 	"name":           fullname,
	// 	"country":        Country,
	// 	"state":          State,
	// 	"district":       District,
	// 	"street_name":    StreetName,
	// 	"pin_code":       PinCode,
	// 	"phone":          Phone,
	// 	"invoice":        invoiceitem,
	// 	"grand_total":    grandtotal,
	// 	"message":        "invoice generated succesfully",
	// 	//"pdf_file":       fileName,
	// })
}

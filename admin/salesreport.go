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
	"github.com/xuri/excelize/v2"
)

func GenerateSalesReport(c *gin.Context) {
	database.DB.Exec(`TRUNCATE sales_report_items`)
	var OrderItems []models.OrderItems
	//database.DB.Find(&OrderItems)
	database.DB.Raw(`SELECT * FROM order_items WHERE order_status='delivered' order by product_id,order_id`).Scan(&OrderItems)
	var orderid uint
	var productid string
	for _, v := range OrderItems {

		fmt.Println("first orderid ", orderid)
		fmt.Println("first productid ", productid)
		if orderid == v.OrderID && productid == v.ProductID {
			fmt.Println("is it entering here?")
			var qty uint
			database.DB.Model(&models.SalesReportItem{}).Where("order_id = ? AND product_id = ?", v.OrderID, v.ProductID).Pluck("qty", &qty)
			qty = qty + 1
			var coupondiscount float64
			database.DB.Model(&models.SalesReportItem{}).Where("order_id = ? AND product_id = ?", v.OrderID, v.ProductID).Pluck("coupon_discount", &coupondiscount)
			coupondiscount = coupondiscount + v.CouponDiscount
			coupondiscount = math.Round(coupondiscount*100) / 100
			var offerdiscount float64
			database.DB.Model(&models.SalesReportItem{}).Where("order_id = ? AND product_id = ?", v.OrderID, v.ProductID).Pluck("offer_discount", &offerdiscount)
			offerdiscount = offerdiscount + v.OfferDiscount
			var totaldiscount float64
			database.DB.Model(&models.SalesReportItem{}).Where("order_id = ? AND product_id = ?", v.OrderID, v.ProductID).Pluck("total_discount", &totaldiscount)
			totaldiscount = totaldiscount + v.TotalDiscount
			var paidamount float64
			database.DB.Model(&models.SalesReportItem{}).Where("order_id = ? AND product_id = ?", v.OrderID, v.ProductID).Pluck("paid_amount", &paidamount)
			paidamount = paidamount + v.PaidAmount
			//CouponDiscount1 := math.Round(v.CouponDiscount*100) / 100
			salesreport := models.SalesReportItem{
				Qty:            qty,
				CouponDiscount: coupondiscount,
				OfferDiscount:  offerdiscount,
				TotalDiscount:  totaldiscount,
				PaidAmount:     paidamount,
			}
			database.DB.Model(&models.SalesReportItem{}).Where("order_id = ? AND product_id = ?", orderid, productid).Updates(&salesreport)
			continue
		}

		orderid = v.OrderID
		productid = v.ProductID
		fmt.Println("second orderid ", orderid)
		fmt.Println("second productid ", productid)
		var productname string
		database.DB.Raw(`SELECT product_name from products where id = ?`, v.ProductID).Scan(&productname)
		fmt.Println("hi hello ", productname)
		CouponDiscount1 := math.Round(v.CouponDiscount*100) / 100
		fmt.Println("======", CouponDiscount1)
		salesreport := models.SalesReportItem{

			//coupondiscount = math.Round(coupondiscount*100) / 100
			OrderID:        orderid,
			ProductID:      productid,
			ProductName:    productname,
			Qty:            1,
			Price:          v.Price,
			OrderStatus:    v.OrderStatus,
			PaymentMethod:  v.PaymentMethod,
			CouponDiscount: CouponDiscount1,
			OfferDiscount:  v.OfferDiscount,
			TotalDiscount:  v.TotalDiscount,
			PaidAmount:     v.PaidAmount,
			OrderDate:      v.CreatedAt,
			DeliveredDate:  v.DeliveredDate,
		}
		database.DB.Create(&salesreport)

	}
	var salesreportitem []models.SalesReportItem
	database.DB.Find(&salesreportitem)
	c.JSON(http.StatusOK, gin.H{
		"sales report": salesreportitem,
		"message":      "sales report generated succesfully",
	})
}
func FilterSalesReport(c *gin.Context) {
	day := c.Query("day")
	month := c.Query("month")
	year := c.Query("year")
	var salesreportitem []models.SalesReportItem
	type SalesReportSummary struct {
		TotalQuantity   uint    `json:"total_quantity"`
		TotalPaidAmount float64 `json:"total_paid_amount"`
		TotalDiscount   float64 `json:"total_discount"`
	}
	var summary SalesReportSummary
	//database.DB.Find(&salesreportitem)
	sql := `SELECT * FROM sales_report_items `
	sql1 := `SELECT SUM(qty)AS total_quantity, SUM(paid_amount) AS total_paid_amount, SUM(total_discount) AS total_discount FROM sales_report_items `
	if day != "" {
		sql = sql + `WHERE DATE(order_date) ='` + day + `'`
		database.DB.Raw(sql).Scan(&salesreportitem)
		sql1 = sql1 + `WHERE DATE(order_date) ='` + day + `'`
		database.DB.Raw(sql1).Scan(&summary)
	} else if month != "" {
		sql = sql + `WHERE EXTRACT(MONTH FROM order_date) = '` + month + `'`
		database.DB.Raw(sql).Scan(&salesreportitem)
		sql1 = sql1 + `WHERE EXTRACT(MONTH FROM order_date) = '` + month + `'`
		database.DB.Raw(sql1).Scan(&summary)
	} else if year != "" {
		sql = sql + `WHERE EXTRACT(YEAR FROM order_date) = '` + year + `'`
		database.DB.Raw(sql).Scan(&salesreportitem)
		sql1 = sql1 + `WHERE EXTRACT(YEAR FROM order_date) = '` + year + `'`
		database.DB.Raw(sql1).Scan(&summary)
	} else {
		database.DB.Raw(sql).Scan(&salesreportitem)
		database.DB.Raw(sql1).Scan(&summary)
	}
	fmt.Println("ooo", summary)
	fmt.Println("---", summary.TotalQuantity)

	c.JSON(http.StatusOK, gin.H{
		"sales report":         salesreportitem,
		"message":              "sales report generated succesfully",
		"Overall Sales Count":  summary.TotalQuantity,
		"Overall Order Amount": summary.TotalPaidAmount,
		"Overall Discount":     summary.TotalDiscount,
	})
}

type SalesReportSummary struct {
	TotalQuantity   uint    `json:"total_quantity"`
	TotalPaidAmount float64 `json:"total_paid_amount"`
	TotalDiscount   float64 `json:"total_discount"`
}

func FilterSalesReportPdfExcel(c *gin.Context) {
	day := c.Query("day")
	month := c.Query("month")
	year := c.Query("year")
	var salesreportitem []models.SalesReportItem

	var summary SalesReportSummary

	sql := `SELECT * FROM sales_report_items `
	sql1 := `SELECT SUM(qty)AS total_quantity, SUM(paid_amount) AS total_paid_amount, SUM(total_discount) AS total_discount FROM sales_report_items `
	if day != "" {
		sql += `WHERE DATE(order_date) ='` + day + `'`
		sql1 += `WHERE DATE(order_date) ='` + day + `'`
	} else if month != "" {
		sql += `WHERE EXTRACT(MONTH FROM order_date) = '` + month + `'`
		sql1 += `WHERE EXTRACT(MONTH FROM order_date) = '` + month + `'`
	} else if year != "" {
		sql += `WHERE EXTRACT(YEAR FROM order_date) = '` + year + `'`
		sql1 += `WHERE EXTRACT(YEAR FROM order_date) = '` + year + `'`
	}

	database.DB.Raw(sql).Scan(&salesreportitem)
	database.DB.Raw(sql1).Scan(&summary)

	format := c.Query("format")
	fileName := "sales_report." + format

	// Use the current directory or a specific path
	// filePath := "path/to/your/directory/" + fileName
	// Or use the temp directory
	filePath := fileName // Using the current directory for simplicity

	if format == "xlsx" {
		err := GenerateExcelReport(salesreportitem, filePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate Excel report"})
			return
		}
	} else if format == "pdf" {
		err := GeneratePDFReport(salesreportitem, summary, filePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF report"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format"})
		return
	}

	// Serve the file to the client
	c.File(filePath)

	// Optionally, you can remove the file after serving it
	// os.Remove(filePath)

	c.JSON(http.StatusOK, gin.H{
		"sales report":         salesreportitem,
		"message":              "sales report generated successfully",
		"Overall Sales Count":  summary.TotalQuantity,
		"Overall Order Amount": summary.TotalPaidAmount,
		"Overall Discount":     summary.TotalDiscount,
	})
}

func GenerateExcelReport(salesReportItems []models.SalesReportItem, filePath string) error {

	fmt.Println("Hello 0")
	f := excelize.NewFile()

	fmt.Println("Hello 1")
	// Create a new sheet
	index, _ := f.NewSheet("SalesReport")

	fmt.Println("Hello 2")
	// Set headers
	headers := []string{"Order ID", "Product ID", "Product Name", "Quantity", "Price", "Order Status", "Payment Method", "Coupon Discount", "Offer Discount", "Total Discount", "Paid Amount", "Order Date"}
	fmt.Println("hello 3")
	for i, header := range headers {
		fmt.Println("hello 4")
		col := string('A' + i)
		f.SetCellValue("SalesReport", col+"1", header)
	}

	// Fill in data
	for i, item := range salesReportItems {
		fmt.Println("hello 5")
		row := strconv.Itoa(i + 2)
		f.SetCellValue("SalesReport", "A"+row, item.OrderID)
		f.SetCellValue("SalesReport", "B"+row, item.ProductID)
		f.SetCellValue("SalesReport", "C"+row, item.ProductName)
		f.SetCellValue("SalesReport", "D"+row, item.Qty)
		f.SetCellValue("SalesReport", "E"+row, item.Price)
		f.SetCellValue("SalesReport", "F"+row, item.OrderStatus)
		f.SetCellValue("SalesReport", "G"+row, item.PaymentMethod)
		f.SetCellValue("SalesReport", "H"+row, item.CouponDiscount)
		f.SetCellValue("SalesReport", "I"+row, item.OfferDiscount)
		f.SetCellValue("SalesReport", "J"+row, item.TotalDiscount)
		f.SetCellValue("SalesReport", "K"+row, item.PaidAmount)
		f.SetCellValue("SalesReport", "L"+row, item.OrderDate.Format("2006-01-02"))
		f.SetCellValue("SalesReport", "M"+row, item.DeliveredDate)
	}

	fmt.Println("hello 6")
	// Set the active sheet
	f.SetActiveSheet(index)

	fmt.Println("hello 7")
	// Save the spreadsheet
	if err := f.SaveAs(filePath); err != nil {
		fmt.Println("Error saving file:", err)
		return err
	}

	fmt.Println("hello 8")
	return nil
}
func GeneratePDFReport(salesReportItems []models.SalesReportItem, summary SalesReportSummary, filePath string) error {
	pdf := gofpdf.New("P", "mm", "A2", "")

	// Add a new page
	pdf.AddPage()

	// Set the title
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Sales Report")
	pdf.Ln(12)

	// Add summary information
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, fmt.Sprintf("Overall Sales Count: %d", summary.TotalQuantity))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Overall Order Amount: Rs.%.2f", summary.TotalPaidAmount))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Overall Discount: Rs.%.2f", summary.TotalDiscount))
	pdf.Ln(15)

	// Set headers with a background color
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(240, 240, 240) // Light grey background
	headers := []string{"Order ID", "Product ID", "Product Name", "Quantity", "Price", "Order Status", "Payment Method", "Coupon Discount", "Offer Discount", "Total Discount", "Paid Amount", "Order Date", "Delivered Date"}
	widths := []float64{15, 20, 30, 20, 25, 35, 35, 35, 35, 35, 30, 25, 35} // Adjust the widths based on content
	for i, header := range headers {
		pdf.CellFormat(widths[i], 10, header, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// Set data with alternating row colors
	pdf.SetFont("Arial", "", 10)
	fill := false
	for _, item := range salesReportItems {
		if fill {
			pdf.SetFillColor(230, 230, 230) // Slightly darker grey for alternating rows
		} else {
			pdf.SetFillColor(255, 255, 255) // White for other rows
		}
		fill = !fill

		pdf.CellFormat(widths[0], 10, strconv.Itoa(int(item.OrderID)), "1", 0, "C", true, 0, "")
		pdf.CellFormat(widths[1], 10, item.ProductID, "1", 0, "C", true, 0, "")
		pdf.CellFormat(widths[2], 10, item.ProductName, "1", 0, "C", true, 0, "")
		pdf.CellFormat(widths[3], 10, strconv.Itoa(int(item.Qty)), "1", 0, "C", true, 0, "")
		pdf.CellFormat(widths[4], 10, fmt.Sprintf("%.2f", item.Price), "1", 0, "C", true, 0, "")
		pdf.CellFormat(widths[5], 10, item.OrderStatus, "1", 0, "C", true, 0, "")
		pdf.CellFormat(widths[6], 10, item.PaymentMethod, "1", 0, "C", true, 0, "")
		pdf.CellFormat(widths[7], 10, fmt.Sprintf("%.2f", item.CouponDiscount), "1", 0, "C", true, 0, "")
		pdf.CellFormat(widths[8], 10, fmt.Sprintf("%.2f", item.OfferDiscount), "1", 0, "C", true, 0, "")
		pdf.CellFormat(widths[9], 10, fmt.Sprintf("%.2f", item.TotalDiscount), "1", 0, "C", true, 0, "")
		pdf.CellFormat(widths[10], 10, fmt.Sprintf("%.2f", item.PaidAmount), "1", 0, "C", true, 0, "")
		pdf.CellFormat(widths[11], 10, item.OrderDate.Format("2006-01-02"), "1", 0, "C", true, 0, "")
		pdf.CellFormat(widths[12], 10, item.DeliveredDate, "1", 0, "C", true, 0, "")
		pdf.Ln(-1)
	}

	// Footer with the current date
	pdf.SetY(-15)
	pdf.SetFont("Arial", "I", 8)
	pdf.Cell(0, 10, fmt.Sprintf("Generated on %s", time.Now().Format("2006-01-02")))

	// Save the PDF to file
	return pdf.OutputFileAndClose(filePath)
}

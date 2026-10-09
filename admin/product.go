package admin

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/helper"
	"github.com/Ansalps/GeZOne/models"
	requestmodemodels "github.com/Ansalps/GeZOne/requestmodels"
	"github.com/Ansalps/GeZOne/responsemodels"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ReadProducts(c *gin.Context) {
	listOrder := c.Query("list_order")
	category := c.Query("category")

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
			COALESCE((SELECT COALESCE(SUM(stock),0) FROM product_variants WHERE product_id = p.id), 0) AS stock,
			p.popular,
			COALESCE(string_agg(DISTINCT pv.size, ',' ORDER BY pv.size), '') AS size,
			-- Only expose discount_percentage when the offer is currently active
			COALESCE(MAX(CASE WHEN o.start_at <= now() AND (o.end_at IS NULL OR o.end_at >= now()) THEN o.discount_percentage ELSE 0 END), 0) AS discount_percentage,
			COALESCE(MAX(to_char(o.start_at, 'YYYY-MM-DD')), '') AS start_date,
			COALESCE(MAX(to_char(o.end_at, 'YYYY-MM-DD')), '') AS end_date
		FROM products p
		JOIN categories c 
			ON c.id = p.category_id
		LEFT JOIN product_variants pv
			ON pv.product_id = p.id
		LEFT JOIN offers o 
			ON p.id = o.product_id
	`

	// Parameters for the SQL query
	var args []interface{}

	// Filter by category if selected
	if category != "" {
		sql += ` WHERE LOWER(TRIM(c.category_name)) = LOWER(TRIM(?))`
		args = append(args, category)
	}

	// Sorting
	sql += ` GROUP BY p.id, p.created_at, p.updated_at, p.category_id, c.category_name,
		p.product_name, p.description, p.image_url, p.price, p.popular `

	switch listOrder {
	case "DSC":
		sql += ` ORDER BY p.created_at DESC`
	default:
		sql += ` ORDER BY p.created_at ASC`
	}

	var products []responsemodels.Product

	if err := database.DB.Raw(sql, args...).Scan(&products).Error; err != nil {
		log.Println("ReadProducts error:", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to retrieve products from database",
		})
		return
	}

	// Batch fetch variants for all products to avoid per-product queries and ensure accurate stock
	if len(products) > 0 {
		ids := make([]uint, 0, len(products))
		for _, p := range products {
			ids = append(ids, p.ID)
		}

		var allVariants []models.ProductVariant
		if err := database.DB.Where("product_id IN ?", ids).Order("product_id ASC, size ASC").Find(&allVariants).Error; err == nil {
			// map productID -> []ProductInventoryItem
			variantMap := make(map[uint][]responsemodels.ProductInventoryItem)
			for _, v := range allVariants {
				variantMap[v.ProductID] = append(variantMap[v.ProductID], responsemodels.ProductInventoryItem{Size: v.Size, Stock: v.Stock})
			}

			for i := range products {
				if inv, ok := variantMap[products[i].ID]; ok {
					products[i].Inventory = inv
				} else {
					products[i].Inventory = []responsemodels.ProductInventoryItem{}
				}
			}
		}
	}

	// If product has offer dates returned as strings from SQL, they're already set in product.StartDate / EndDate.

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Successfully retrieved products",
		"data": gin.H{
			"products": products,
		},
	})
}
func ReadProductById(c *gin.Context) {
	ProductID := c.Param("id")
	var product models.Product

	err := database.DB.Preload("Variants").First(&product, "id = ?", ProductID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  false,
				"message": "category id does not exist",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "database error while reading category by id",
		})
		return
	}

	var category models.Category

	err = database.DB.First(&category, "id = ?", product.CategoryID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  false,
				"message": "category id does not exist",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "database error while reading category by id",
		})
		return
	}

	var totalStock uint
	var sizes []string
	inventory := make([]responsemodels.ProductInventoryItem, 0, len(product.Variants))
	for _, variant := range product.Variants {
		totalStock += variant.Stock
		sizes = append(sizes, variant.Size)
		inventory = append(inventory, responsemodels.ProductInventoryItem{
			Size:  variant.Size,
			Stock: variant.Stock,
		})
	}

	// Attempt to fetch any existing offer for this product to include discount and dates
	var offer models.Offer
	discountPercent := int64(0)
	startDateStr := ""
	endDateStr := ""
	hasOffer := false
	if err := database.DB.Where("product_id = ?", product.ID).First(&offer).Error; err == nil {
		hasOffer = true
		// Only treat the offer as active for discount if current time is within the offer window
		now := time.Now()
		active := !offer.StartAt.IsZero() && offer.StartAt.After(time.Time{}) && offer.StartAt.Before(now.Add(time.Second))
		// If start is zero treat as already started
		if offer.StartAt.IsZero() {
			active = true
		}
		if !offer.EndAt.IsZero() && offer.EndAt.Before(now) {
			active = false
		}

		if active {
			discountPercent = int64(offer.DiscountPercentage)
		}

		if !offer.StartAt.IsZero() {
			startDateStr = offer.StartAt.Format("2006-01-02")
		}
		if !offer.EndAt.IsZero() {
			endDateStr = offer.EndAt.Format("2006-01-02")
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data": responsemodels.Product{
			ID:                 product.ID,
			CreatedAt:          product.CreatedAt,
			CategoryID:         product.CategoryID,
			CategoryName:       category.CategoryName,
			ProductName:        product.ProductName,
			ProductDescription: product.Description,
			ProductImageUrl:    product.ImageURL,
			Price:              product.Price,
			Stock:              int64(totalStock),
			Popular:            product.Popular,
			Size:               strings.Join(sizes, ","),
			Inventory:          inventory,
			DiscountPercentage: discountPercent,
			StartDate:          startDateStr,
			EndDate:            endDateStr,
			OfferID:            offer.ID,
			HasOffer:           hasOffer,
		},
	})
}
func AddProduct(c *gin.Context) {

	// 1. Bind multipart/form-data fields
	var req requestmodemodels.Product

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Failed to bind request",
		})
		return
	}

	inventoryJSON := strings.TrimSpace(c.PostForm("inventory"))
	inventoryProvided := inventoryJSON != ""
	var inventory []models.ProductVariant

	if inventoryProvided {
		if err := json.Unmarshal([]byte(inventoryJSON), &inventory); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Invalid inventory format",
			})
			return
		}

		if len(inventory) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Inventory cannot be empty",
			})
			return
		}

		for _, item := range inventory {
			if item.Stock < 0 {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  false,
					"message": "Inventory stock must not be negative",
				})
				return
			}
			if item.Size != "Small" && item.Size != "Medium" && item.Size != "Large" {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  false,
					"message": "Invalid size selection",
				})
				return
			}
		}

		req.Size = inventory[0].Size
		req.Stock = inventory[0].Stock
	}

	if inventoryProvided {
		if req.CategoryName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "category_name is required"})
			return
		}
		if req.ProductName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "product_name is required"})
			return
		}
		if req.Description == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "product_description is required"})
			return
		}
		if req.Price <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "price must be greater than zero"})
			return
		}
	} else {
		if err := helper.Validate(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": err.Error(),
			})
			return
		}
	}

	if !inventoryProvided && (req.Size != "Small" && req.Size != "Medium" && req.Size != "Large") {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid size selection",
		})
		return
	}

	if req.DiscountPercentage < 0 || req.DiscountPercentage > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Discount percentage must be between 0 and 100",
		})
		return
	}

	fileHeader, err := c.FormFile("product_image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Product image is required",
		})
		return
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/webp" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Only JPG, PNG, and WEBP images are allowed",
		})
		return
	}

	imageURL, err := helper.UploadToS3(c.Request.Context(), fileHeader, "products", "S3_BUCKET_NAME")
	if err != nil {
		log.Println("S3 Upload Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to upload product image to S3",
		})
		return
	}

	var categoryID uint
	err = database.DB.
		Model(&models.Category{}).
		Select("id").
		Where("category_name = ?", req.CategoryName).
		Scan(&categoryID).
		Error

	if err != nil {
		log.Println("Category lookup error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to find category",
		})
		return
	}

	if categoryID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Category does not exist",
		})
		return
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to start transaction",
		})
		return
	}

	product := models.Product{
		CategoryID:  categoryID,
		ProductName: req.ProductName,
		Description: req.Description,
		ImageURL:    imageURL,
		Price:       req.Price,
		Popular:     req.Popular,
	}

	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback()
		log.Println("Product creation error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to create product",
		})
		return
	}

	if inventoryProvided {
		for _, item := range inventory {
			variant := models.ProductVariant{
				ProductID: product.ID,
				Size:      item.Size,
				Stock:     item.Stock,
			}
			if err := tx.Create(&variant).Error; err != nil {
				tx.Rollback()
				log.Println("Product variant creation error:", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  false,
					"message": "Failed to create product variant",
				})
				return
			}
		}
	} else {
		variant := models.ProductVariant{
			ProductID: product.ID,
			Size:      req.Size,
			Stock:     req.Stock,
		}
		if err := tx.Create(&variant).Error; err != nil {
			tx.Rollback()
			log.Println("Product variant creation error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to create product variant",
			})
			return
		}
	}

	if req.DiscountPercentage > 0 {
		startAt, endAt, perr := ParseProductOfferDates(req.StartDate, req.EndDate)
		if perr != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": perr.Error()})
			return
		}

		offer := models.Offer{
			ProductID:          product.ID,
			DiscountPercentage: req.DiscountPercentage,
		}

		if !startAt.IsZero() {
			offer.StartAt = startAt
		}
		if !endAt.IsZero() {
			offer.EndAt = endAt
		}

		if err := tx.Create(&offer).Error; err != nil {
			tx.Rollback()
			log.Println("Offer creation error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to create product offer",
			})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Transaction commit error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to save product",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  true,
		"message": "Product Added Successfully",
		"data":    product,
	})
}

func ParseProductOfferDates(startStr, endStr string) (time.Time, time.Time, error) {
	var startAt, endAt time.Time
	var err error

	layouts := []string{
		"2006-01-02T15:04",
		"2006-01-02 15:04:05",
		time.RFC3339,
		"2006-01-02",
	}

	parseDate := func(dateStr string) (time.Time, error) {
		dateStr = strings.TrimSpace(dateStr)
		if dateStr == "" {
			return time.Time{}, nil
		}
		for _, layout := range layouts {
			if t, e := time.Parse(layout, dateStr); e == nil {
				return t, nil
			}
		}
		return time.Time{}, fmt.Errorf("invalid date format: %s", dateStr)
	}

	if startStr != "" {
		startAt, err = parseDate(startStr)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	}

	if endStr != "" {
		endAt, err = parseDate(endStr)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	}

	if !startAt.IsZero() && !endAt.IsZero() && endAt.Before(startAt) {
		return time.Time{}, time.Time{}, errors.New("end date must be after start date")
	}

	return startAt, endAt, nil
}

func EditProduct(c *gin.Context) {
	productID := c.Param("id")

	// 1. Bind multipart/form-data
	var reqProduct requestmodemodels.Product
	if err := c.ShouldBind(&reqProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Failed to bind form fields",
		})
		return
	}

	// 2. Parse inventory JSON
	inventoryJSON := strings.TrimSpace(c.PostForm("inventory"))
	type InventoryInput struct {
		Size  string `json:"size"`
		Stock uint   `json:"stock"`
	}
	var inventoryInputs []InventoryInput

	if inventoryJSON != "" {
		if err := json.Unmarshal([]byte(inventoryJSON), &inventoryInputs); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Invalid inventory format",
			})
			return
		}

		for _, item := range inventoryInputs {
			if item.Size != "Small" && item.Size != "Medium" && item.Size != "Large" {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  false,
					"message": "Invalid size choice in inventory (must be Small, Medium, or Large)",
				})
				return
			}
		}
	}

	// 3. Validate main request fields
	if reqProduct.CategoryName == "" || reqProduct.ProductName == "" || reqProduct.Description == "" || reqProduct.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Missing required fields (category_name, product_name, description, price)",
		})
		return
	}

	if reqProduct.DiscountPercentage < 0 || reqProduct.DiscountPercentage > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Discount percentage must be between 0 and 100",
		})
		return
	}

	// 4. Resolve Category ID
	var category models.Category
	if err := database.DB.Where("category_name = ?", reqProduct.CategoryName).First(&category).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Category does not exist",
		})
		return
	}

	// 5. Check if product exists
	var product models.Product
	if err := database.DB.Where("id = ?", productID).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"status": false, "message": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Database query failed"})
		return
	}

	// 6. Inspect & Upload Image to S3 if a new file is present
	fileHeader, fileErr := c.FormFile("product_image")
	var newImageURL string

	if fileErr == nil {
		f, openErr := fileHeader.Open()
		if openErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to read uploaded image"})
			return
		}
		defer f.Close()

		buffer := make([]byte, 512)
		_, _ = f.Read(buffer)
		contentType := http.DetectContentType(buffer)

		if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/webp" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Only JPG, PNG, and WEBP images are allowed",
			})
			return
		}

		uploadedURL, err := helper.UploadToS3(c.Request.Context(), fileHeader, "products", "S3_BUCKET_NAME")
		if err != nil {
			log.Println("S3 Upload Error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to upload product image to S3",
			})
			return
		}
		newImageURL = uploadedURL
	}

	// 7. Begin Database Transaction
	tx := database.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to start database transaction"})
		return
	}

	// 8. Update Product Record
	updateData := map[string]interface{}{
		"category_id":  category.ID,
		"product_name": reqProduct.ProductName,
		"description":  reqProduct.Description,
		"price":        reqProduct.Price,
		"popular":      reqProduct.Popular,
	}

	if newImageURL != "" {
		updateData["image_url"] = newImageURL
	}

	if err := tx.Model(&product).Updates(updateData).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to update product core information"})
		return
	}

	// 9. Process Product Variants / Stock
	if len(inventoryInputs) > 0 {
		for _, item := range inventoryInputs {
			var variant models.ProductVariant
			err := tx.Where("product_id = ? AND size = ?", product.ID, item.Size).First(&variant).Error

			if errors.Is(err, gorm.ErrRecordNotFound) {
				newVariant := models.ProductVariant{
					ProductID: product.ID,
					Size:      item.Size,
					Stock:     item.Stock,
				}
				if err := tx.Create(&newVariant).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to create variant"})
					return
				}
			} else if err == nil {
				if err := tx.Model(&variant).Update("stock", item.Stock).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to update variant stock"})
					return
				}
			} else {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to query variant record"})
				return
			}
		}
	}

	// 10. Process Offers
	var existingOffer models.Offer
	offerErr := tx.Where("product_id = ?", product.ID).First(&existingOffer).Error

	if reqProduct.DiscountPercentage > 0 {
		startAt, endAt, perr := ParseProductOfferDates(reqProduct.StartDate, reqProduct.EndDate)
		if perr != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": perr.Error()})
			return
		}

		if offerErr == nil {
			// Update Existing Offer
			offerUpdates := map[string]interface{}{
				"discount_percentage": reqProduct.DiscountPercentage,
				"start_at":            startAt,
				"end_at":              endAt,
			}
			if err := tx.Model(&existingOffer).Updates(offerUpdates).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to update offer"})
				return
			}
		} else if errors.Is(offerErr, gorm.ErrRecordNotFound) {
			// Create New Offer
			newOffer := models.Offer{
				ProductID:          product.ID,
				DiscountPercentage: reqProduct.DiscountPercentage,
				StartAt:            startAt,
				EndAt:              endAt,
			}
			if err := tx.Create(&newOffer).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to create offer"})
				return
			}
		}
	} else {
		// Delete active offer if discount percentage is set to 0
		if offerErr == nil {
			if err := tx.Delete(&existingOffer).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to delete existing offer"})
				return
			}
		}
	}

	// 11. Commit Transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to commit database updates"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Product updated successfully",
	})
}

func ProductDelete(c *gin.Context) {
	ProductID := c.Param("id")

	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM products WHERE id = ?`, ProductID).Scan(&count)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "product id does not exist",
		})
		return
	}

	database.DB.Where("product_id = ?", ProductID).Delete(&models.ProductVariant{})
	database.DB.Where("id = ?", ProductID).Delete(&models.Product{})
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "product deleted successfully"})
}

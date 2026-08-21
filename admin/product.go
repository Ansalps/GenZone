package admin

import (
	"errors"
	"fmt"
	"log"
	"net/http"

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
			p.stock,
			p.popular,
			p.size,
			COALESCE(o.discount_percentage, 0) AS discount_percentage
		FROM products p
		JOIN categories c 
			ON c.id = p.category_id
		LEFT JOIN offers o 
			ON p.id = o.product_id
	`

	// Parameters for the SQL query
	var args []interface{}

	// Filter by category if selected
	if category != "" {
		sql += ` WHERE c.category_name = ?`
		args = append(args, category)
	}

	// Sorting
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

	fmt.Println("category:", category)
	fmt.Println("products:", products)

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

	// Pass the pointer &category as the target for First()
	err := database.DB.First(&product, "id = ?", ProductID).Error
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

	// Pass the pointer &category as the target for First()
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

	// Map your database model to your response model
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
			Stock:              int64(product.Stock),
			Popular:            product.Popular,
			Size:               product.Size,
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

	// 2. Validate request fields
	if err := helper.Validate(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	// 3. Validate size
	if req.Size != "Small" &&
		req.Size != "Medium" &&
		req.Size != "Large" {

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid size selection",
		})
		return
	}

	// 4. Validate discount percentage
	if req.DiscountPercentage < 0 ||
		req.DiscountPercentage > 100 {

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Discount percentage must be between 0 and 100",
		})
		return
	}

	// 5. Get uploaded image
	fileHeader, err := c.FormFile("product_image")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Product image is required",
		})
		return
	}

	// Optional: validate image content type
	contentType := fileHeader.Header.Get("Content-Type")

	if contentType != "image/jpeg" &&
		contentType != "image/png" &&
		contentType != "image/webp" {

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Only JPG, PNG, and WEBP images are allowed",
		})
		return
	}

	// 6. Upload image to S3
	imageURL, err := helper.UploadToS3(fileHeader)

	if err != nil {
		log.Println("S3 Upload Error:", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to upload product image to S3",
		})
		return
	}

	// 7. Find category ID
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

	// 8. Begin transaction
	tx := database.DB.Begin()

	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to start transaction",
		})
		return
	}

	// 9. Create product
	product := models.Product{
		CategoryID:  categoryID,
		ProductName: req.ProductName,
		Description: req.Description,
		ImageURL:    imageURL,
		Price:       req.Price,
		Stock:       req.Stock,
		Size:        req.Size,
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

	// 10. Create offer if discount exists
	if req.DiscountPercentage > 0 {

		offer := models.Offer{
			ProductID:          product.ID,
			DiscountPercentage: req.DiscountPercentage,
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

	// 11. Commit transaction
	if err := tx.Commit().Error; err != nil {

		log.Println("Transaction commit error:", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to save product",
		})
		return
	}

	// 12. Return success
	c.JSON(http.StatusCreated, gin.H{
		"status":  true,
		"message": "Product Added Successfully",
		"data":    product,
	})
}

func EditProduct(c *gin.Context) {
	productID := c.Param("id")

	// 1. Bind multipart/form-data fields
	var reqProduct requestmodemodels.Product

	if err := c.ShouldBind(&reqProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Failed to bind request",
		})
		return
	}

	// 2. Validate request
	if err := helper.Validate(reqProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"message":    err.Error(),
			"error_code": http.StatusBadRequest,
		})
		return
	}

	// 3. Validate size
	if reqProduct.Size != "Small" &&
		reqProduct.Size != "Medium" &&
		reqProduct.Size != "Large" {

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid size selection",
		})
		return
	}

	// 4. Validate discount
	if reqProduct.DiscountPercentage < 0 ||
		reqProduct.DiscountPercentage > 100 {

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Discount percentage must be between 0 and 100",
		})
		return
	}

	// 5. Find category
	var categoryID uint

	err := database.DB.
		Model(&models.Category{}).
		Select("id").
		Where(
			"category_name = ?",
			reqProduct.CategoryName,
		).
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

	// 6. Find existing product
	var product models.Product

	err = database.DB.
		Where(
			"id = ?",
			productID,
		).
		First(&product).
		Error

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  false,
				"message": "Product not found",
			})
			return
		}

		log.Println("Product lookup error:", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to find product",
		})
		return
	}

	// 7. Check whether a new image was uploaded
	fileHeader, fileErr := c.FormFile("product_image")

	var newImageURL string

	if fileErr == nil {

		// Validate image type
		contentType := fileHeader.Header.Get("Content-Type")

		if contentType != "image/jpeg" &&
			contentType != "image/png" &&
			contentType != "image/webp" {

			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Only JPG, PNG, and WEBP images are allowed",
			})
			return
		}

		// Upload new image to S3
		newImageURL, err = helper.UploadToS3(fileHeader)

		if err != nil {
			log.Println("S3 Upload Error:", err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to upload product image to S3",
			})
			return
		}
	}

	// 8. Begin transaction
	tx := database.DB.Begin()

	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to start transaction",
		})
		return
	}

	// 9. Prepare product update
	updateData := map[string]interface{}{
		"category_id":  categoryID,
		"product_name": reqProduct.ProductName,
		"description":  reqProduct.Description,
		"price":        reqProduct.Price,
		"stock":        reqProduct.Stock,
		"size":         reqProduct.Size,
		"popular":      reqProduct.Popular,
	}

	// Only update image if a new image was selected
	if fileErr == nil {
		updateData["image_url"] = newImageURL
	}

	// 10. Update product
	result := tx.
		Model(&product).
		Updates(updateData)

	if result.Error != nil {
		tx.Rollback()

		log.Println("Product update error:", result.Error)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to update product",
		})
		return
	}

	// 11. Find existing offer
	var existingOffer models.Offer

	offerErr := tx.
		Where(
			"product_id = ? AND deleted_at IS NULL",
			product.ID,
		).
		First(&existingOffer).
		Error

	// 12. Handle offer
	if reqProduct.DiscountPercentage > 0 {

		if offerErr == nil {

			// Existing offer → update
			err = tx.
				Model(&existingOffer).
				Update(
					"discount_percentage",
					reqProduct.DiscountPercentage,
				).
				Error

			if err != nil {
				tx.Rollback()

				log.Println("Offer update error:", err)

				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  false,
					"message": "Failed to update offer",
				})
				return
			}

		} else if errors.Is(offerErr, gorm.ErrRecordNotFound) {

			// No offer → create new offer
			newOffer := models.Offer{
				ProductID:          product.ID,
				DiscountPercentage: reqProduct.DiscountPercentage,
			}

			if err := tx.Create(&newOffer).Error; err != nil {
				tx.Rollback()

				log.Println("Offer creation error:", err)

				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  false,
					"message": "Failed to create offer",
				})
				return
			}

		} else {

			// Unexpected database error
			tx.Rollback()

			log.Println("Offer lookup error:", offerErr)

			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to check existing offer",
			})
			return
		}

	} else {

		// Discount = 0
		// Remove existing offer
		if offerErr == nil {

			if err := tx.Delete(&existingOffer).Error; err != nil {
				tx.Rollback()

				log.Println("Offer deletion error:", err)

				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  false,
					"message": "Failed to remove offer",
				})
				return
			}

		} else if !errors.Is(offerErr, gorm.ErrRecordNotFound) {

			tx.Rollback()

			log.Println("Offer lookup error:", offerErr)

			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to check existing offer",
			})
			return
		}
	}

	// 13. Commit transaction
	if err := tx.Commit().Error; err != nil {

		log.Println("Transaction commit error:", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to save product",
		})
		return
	}

	// 14. Success
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Product updated successfully",
	})
}

func ProductDelete(c *gin.Context) {
	ProductID := c.Param("id")
	fmt.Println(ProductID)
	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM products WHERE id = ?`, ProductID).Scan(&count)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "product id does not exist",
		})
		return
	}

	database.DB.Where("id = ?", ProductID).Delete(&models.Product{})
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "product deleted successfully"})
}

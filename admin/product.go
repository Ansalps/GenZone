package admin

import (
	"errors"
	"fmt"
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
    listorder := c.Query("list_order")

    // Explicitly select columns to avoid ID collisions between joined tables
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
        JOIN categories c ON c.id = p.category_id AND c.deleted_at IS NULL
        LEFT JOIN offers o ON p.id = o.product_id AND o.deleted_at IS NULL
        WHERE p.deleted_at IS NULL
    `

    // Sorting logic
    switch listorder {
    case "DSC":
        sql += ` ORDER BY p.created_at DESC`
    default: // Handles "ASC" and empty string default
        sql += ` ORDER BY p.created_at ASC`
    }

    var products []responsemodels.Product
    if err := database.DB.Raw(sql).Scan(&products).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status":  false,
            "message": "Failed to retrieve products from database",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status":  true,
        "message": "Successfully retrieved products",
        "data": gin.H{
            "products": products,
        },
    })
}
func ReadProductById(c *gin.Context){
	ProductID:=c.Param("id")
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
			ID:   product.ID,
			CreatedAt: product.CreatedAt,
			CategoryID: product.CategoryID, 
			CategoryName: category.CategoryName,
			ProductName: product.ProductName,
			ProductDescription: product.Description,
			ProductImageUrl: product.ImageURL,
			Price: product.Price,
			Stock: int64(product.Stock),
			Popular: product.Popular,
			Size: product.Size,
			
		},
	})
}
func AddProduct(c *gin.Context) {
    var req requestmodemodels.Product // Ensure JSON tags match front-end payload

    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "status":  false,
            "message": "failed to bind request",
        })
        return
    }

    if err := helper.Validate(req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "status":  false,
            "message": err.Error(),
        })
        return
    }

    // Validate size (if not validated by helper)
    if req.Size != "Small" && req.Size != "Medium" && req.Size != "Large" {
        c.JSON(http.StatusBadRequest, gin.H{
            "status":  false,
            "message": "Invalid size selection",
        })
        return
    }

    // 1. Fetch Category ID in a single query (checking for soft delete)
    var categoryID uint
    err := database.DB.Model(&models.Category{}).
        Select("id").
        Where("category_name = ? AND deleted_at IS NULL", req.CategoryName).
        Scan(&categoryID).Error

    if err != nil || categoryID == 0 {
        c.JSON(http.StatusBadRequest, gin.H{
            "status":  false,
            "message": "Category does not exist",
        })
        return
    }

    // 2. Begin Transaction (Ensures both product and offer insert or both fail)
    tx := database.DB.Begin()

    product := models.Product{
        CategoryID:  categoryID,
        ProductName: req.ProductName,
        Description: req.Description,
        ImageURL:    req.ImageUrl,
        Price:       req.Price,
        Stock:       req.Stock,
        Size:        req.Size,
        Popular:     req.Popular,
    }

    if err := tx.Create(&product).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{
            "status":  false,
            "message": "Failed to create product",
        })
        return
    }

    // 3. Insert into offers table if discount_percentage > 0
    if req.DiscountPercentage > 0 {
        offer := models.Offer{
            ProductID:          product.ID, // Uses created product's ID
            DiscountPercentage: req.DiscountPercentage,
        }

        if err := tx.Create(&offer).Error; err != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{
                "status":  false,
                "message": "Failed to create product offer",
            })
            return
        }
    }

    // Commit transaction
    tx.Commit()

    c.JSON(http.StatusOK, gin.H{
        "status":  true,
        "message": "Product Added Successfully",
    })
}

func EditProduct(c *gin.Context) {
    productID := c.Param("id")

    var reqProduct requestmodemodels.Product
    if err := c.ShouldBindJSON(&reqProduct); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "status":  false,
            "message": "Failed to bind request",
        })
        return
    }

    if err := helper.Validate(reqProduct); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "status":     false,
            "message":    err.Error(),
            "error_code": http.StatusBadRequest,
        })
        return
    }

    if reqProduct.Size != "Medium" && reqProduct.Size != "Small" && reqProduct.Size != "Large" {
        c.JSON(http.StatusBadRequest, gin.H{
            "status":  false,
            "message": "Unknown size",
        })
        return
    }

    // Check category existence
    var categoryID uint
    err := database.DB.Model(&models.Category{}).
        Select("id").
        Where("category_name = ? AND deleted_at IS NULL", reqProduct.CategoryName).
        Scan(&categoryID).Error

    if err != nil || categoryID == 0 {
        c.JSON(http.StatusBadRequest, gin.H{
            "status":  false,
            "message": "Category does not exist",
        })
        return
    }

    // Begin Database Transaction
    tx := database.DB.Begin()

    // 1. Update Product attributes
    updateData := map[string]interface{}{
        "category_id":  categoryID,
        "product_name": reqProduct.ProductName,
        "description":  reqProduct.Description,
        "image_url":    reqProduct.ImageUrl,
        "price":        reqProduct.Price,
        "stock":        reqProduct.Stock,
        "size":         reqProduct.Size,
        "popular":      reqProduct.Popular,
    }

    result := tx.Model(&models.Product{}).
        Where("id = ? AND deleted_at IS NULL", productID).
        Updates(updateData)

    if result.Error != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to update product"})
        return
    }

    if result.RowsAffected == 0 {
        tx.Rollback()
        c.JSON(http.StatusNotFound, gin.H{"status": false, "message": "Product not found"})
        return
    }

    // 2. Handle Offers table updates
    var existingOffer models.Offer
    offerErr := tx.Where("product_id = ? AND deleted_at IS NULL", productID).First(&existingOffer).Error

    if reqProduct.DiscountPercentage > 0 {
        if offerErr == nil {
            // Offer exists -> Update it
            if err := tx.Model(&existingOffer).Update("discount_percentage", reqProduct.DiscountPercentage).Error; err != nil {
                tx.Rollback()
                c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to update offer"})
                return
            }
        } else {
            // No offer exists -> Create new offer
            newOffer := models.Offer{
                ProductID:          existingOffer.ProductID, // Or parse productID string to uint
                DiscountPercentage: reqProduct.DiscountPercentage,
            }
            // Parse productID parameter to uint for new offer
            var pid uint
            fmt.Sscanf(productID, "%d", &pid)
            newOffer.ProductID = pid

            if err := tx.Create(&newOffer).Error; err != nil {
                tx.Rollback()
                c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to create offer"})
                return
            }
        }
    } else {
        // Discount set to 0 -> Soft delete existing offer if it exists
        if offerErr == nil {
            if err := tx.Delete(&existingOffer).Error; err != nil {
                tx.Rollback()
                c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to remove offer"})
                return
            }
        }
    }

    // Commit Transaction
    tx.Commit()

    c.JSON(http.StatusOK, gin.H{"status": true, "message": "Product updated successfully"})
}

func ProductDelete(c *gin.Context) {
	ProductID := c.Param("id")
	fmt.Println(ProductID)
	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM products WHERE id = ? AND deleted_at IS NULL`, ProductID).Scan(&count)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "product id does not exist",
		})
		return
	}

	database.DB.Where("id = ?", ProductID).Delete(&models.Product{})
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "product deleted successfully"})
}

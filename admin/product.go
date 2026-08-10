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
	var product []responsemodels.Product
	//tx := database.DB.Find(&product)
	// tx := database.DB.Raw(`SELECT * FROM categories join products on categories.id=products.category_id and products.deleted_at IS NULL AND categories.deleted_at IS NULL`).Scan(&product)
	sql := `SELECT * FROM categories join products on categories.id=products.category_id and products.deleted_at IS NULL AND categories.deleted_at IS NULL`

	switch listorder{
	case "":
		sql += ` ORDER BY products.created_at ASC`
	case "ASC":
		sql += ` ORDER BY products.created_at ASC`
	case "DSC":
		sql += ` ORDER BY products.created_at DESC`
	}
	
	tx := database.DB.Raw(sql).Scan(&product)
	if tx.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "failed to retrieve data from the database, or the data doesn't exists",
		})
		return
	}
	fmt.Println("response",product)
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "successfully retrieved user informations",
		"data": gin.H{
			"products": product,
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
			Description: product.Description,
			ImageUrl: product.ImageURL,
			Price: product.Price,
			Stock: int(product.Stock),
			Popular: product.Popular,
			Size: product.Size,
			HasOffer: product.HasOffer,
			OfferDiscountPercent: product.OfferDiscountPercent,
			DiscountAmount: uint(product.DiscountAmount),
			TotalDiscountedAmount: uint(product.TotalDiscountedAmount),
		},
	})
}
func AddProduct(c *gin.Context) {
	
	var Product requestmodemodels.Product
	err := c.BindJSON(&Product)
	response := gin.H{
		"status":  false,
		"message": "failed to bind request",
	}
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	fmt.Println("product",Product)
	if err := helper.Validate(Product); err != nil {
		fmt.Println("", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"error_code": http.StatusBadRequest,
		})
		return
	}
	if Product.Price != float64(int(Product.Price)) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "price should not contain decimal places",
		})
		return
	}

	p := Product.Size
	
	if p != "Medium" && p != "Small" && p != "Large" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unknown size",
		})
		return
	}

	fmt.Println("product.categoryname ", Product.CategoryName)
	var count int64
	err = database.DB.Raw(`SELECT COUNT(*) FROM categories WHERE categories.category_name=? and categories.deleted_at is NULL`, Product.CategoryName).Scan(&count).Error
	if err != nil {
		fmt.Println("failed to execute query", err)
	}
	fmt.Println("count", count)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category does not exist"})
		return
	}
	var categoryid uint
	database.DB.Raw(`SELECT id from categories where category_name = ?`, Product.CategoryName).Scan(&categoryid)
	if categoryid == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Category does not exist"})
		return
	}
	//var product models.Product
	product := models.Product{
		CategoryID:  categoryid,
		ProductName: Product.ProductName,
		Description: Product.Description,
		ImageURL:    Product.ImageUrl,
		Price:       Product.Price,
		Stock:       Product.Stock,
		Size:        Product.Size,
		Popular:     Product.Popular,
		HasOffer: Product.HasOffer,
		OfferDiscountPercent: uint(Product.OfferDiscountPercent),
		DiscountAmount: Product.DiscountAmount,
		TotalDiscountedAmount: Product.TotalDiscountedAmount,
	}
	database.DB.Create(&product)
	

	c.JSON(http.StatusOK, gin.H{"status": true, "message": "Product Added Successfully"})

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
    err := database.DB.Raw(`SELECT id FROM categories WHERE category_name = ? AND deleted_at IS NULL`, reqProduct.CategoryName).Scan(&categoryID).Error
    if err != nil || categoryID == 0 {
        c.JSON(http.StatusBadRequest, gin.H{
            "status":  false,
            "message": "Category does not exist",
        })
        return
    }

    // Perform update using a map so boolean/zero values update correctly
    updateData := map[string]interface{}{
        "category_id":             categoryID,
        "product_name":            reqProduct.ProductName,
        "description":             reqProduct.Description,
        "image_url":               reqProduct.ImageUrl,
        "price":                   reqProduct.Price,
        "stock":                   reqProduct.Stock,
        "size":                    reqProduct.Size,
        "popular":                 reqProduct.Popular,
        "has_offer":               reqProduct.HasOffer,
        "offer_discount_percent": reqProduct.OfferDiscountPercent,
        "discount_amount":         reqProduct.DiscountAmount,
        "total_discounted_amount": reqProduct.TotalDiscountedAmount,
    }

    result := database.DB.Model(&models.Product{}).
        Where("id = ? AND deleted_at IS NULL", productID).
        Updates(updateData)

    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to update product"})
        return
    }

    if result.RowsAffected == 0 {
        c.JSON(http.StatusNotFound, gin.H{"status": false, "message": "Product not found"})
        return
    }

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

package admin

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/helper"
	"github.com/Ansalps/GeZOne/models"
	requestmodemodels "github.com/Ansalps/GeZOne/requestmodels"
	"github.com/Ansalps/GeZOne/responsemodels"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ReadCategory(c *gin.Context) {
	
	listorder := c.Query("list_order")
	var category []responsemodels.Category
	//tx := database.DB.Find(&category)
	sql := `SELECT * FROM categories WHERE deleted_at IS NULL`

	switch listorder {
	case "":
		sql += ` ORDER BY categories.created_at ASC`
	case "ASC":
		sql += ` ORDER BY categories.created_at ASC`
	case "DSC":
		sql += ` ORDER BY categories.created_at DESC`
	}
	tx := database.DB.Raw(sql).Scan(&category)
	if tx.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "failed to retrieve data from the database, or the data doesn't exists",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "successfully retrieved user informations",
		"data": gin.H{
			"categories": category,
		},
	})
}
func ReadCategoryById(c *gin.Context) {
	categoryID := c.Param("id")

	var category models.Category

	// Pass the pointer &category as the target for First()
	err := database.DB.First(&category, "id = ?", categoryID).Error
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
		"data": responsemodels.Category{
			ID:   category.ID,
			CategoryName: category.CategoryName,
			Description: category.Description, 
			ImageUrl: category.ImageURL,
		},
	})
}
func AddCategory(c *gin.Context) {
    // 1. Get textual form fields instead of c.BindJSON
    categoryName := c.PostForm("category_name")
    description := c.PostForm("description")
	fmt.Println("hello")
    categoryReq := requestmodemodels.Category{
        CategoryName: categoryName,
        Description:  description,
    }

    // 2. Validate struct fields
    if err := helper.Validate(categoryReq); err != nil {
		fmt.Println("err",err)
        c.JSON(http.StatusBadRequest, gin.H{
            "status":     false,
            "message":    err.Error(),
            "error_code": http.StatusBadRequest,
        })
        return
    }

    // 3. Check for existing category name
    var count int64
    err:=database.DB.Raw(`SELECT COUNT(*) FROM categories WHERE category_name = ? AND deleted_at IS NULL`, categoryName).Scan(&count).Error
	if err!=nil{
		fmt.Println("err",err)
		c.JSON(http.StatusInternalServerError,gin.H{})
	}
    if count != 0 {
		fmt.Println("hi")
        c.JSON(http.StatusBadRequest, gin.H{
            "status":  false,
            "message": "category name already exists",
        })
        return
    }

    // 4. Retrieve the uploaded image file header
    fileHeader, err := c.FormFile("image")
    if err != nil {
		fmt.Println("err",err)
        c.JSON(http.StatusBadRequest, gin.H{
            "status":  false,
            "message": "image file is required",
        })
        return
    }

    // 5. Upload the file to S3 using your helper
    imageURL, err := helper.UploadToS3(fileHeader)
    if err != nil {
        log.Println("S3 Upload Error:", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "status":  false,
            "message": "failed to upload image to S3",
        })
        return
    }

    // 6. Save the category record with the generated S3 image URL
    category := models.Category{
        CategoryName: categoryName,
        Description:  description,
        ImageURL:     imageURL,
    }

    if err := database.DB.Create(&category).Error; err != nil {
        log.Println("DB Create Error:", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "status":  false,
            "message": "database error while adding category",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status":  true,
        "message": "Category Added",
        "data":    category,
    })
}

func EditCategory(c *gin.Context) {
    categoryID := c.Param("id")

    // 1. Find existing category
    var category models.Category

    if err := database.DB.
        Where("id = ? AND deleted_at IS NULL", categoryID).
        First(&category).Error; err != nil {

        c.JSON(http.StatusBadRequest, gin.H{
            "status":  false,
            "message": "category id does not exist",
        })
        return
    }

    // 2. Get form fields
    categoryName := c.PostForm("category_name")
    description := c.PostForm("description")

    // 3. Validate request
    categoryRequest := requestmodemodels.Category{
        CategoryName: categoryName,
        Description:  description,
    }

    if err := helper.Validate(categoryRequest); err != nil {
		fmt.Println("err",err)
        c.JSON(http.StatusBadRequest, gin.H{
            "status":     false,
            "message":    err.Error(),
            "error_code": http.StatusBadRequest,
        })
        return
    }

    // 4. Check if another category already has this name
    var count int64

    err := database.DB.Raw(`
        SELECT COUNT(*)
        FROM categories
        WHERE category_name = ?
        AND id != ?
        AND deleted_at IS NULL
    `, categoryName, categoryID).Scan(&count).Error

    if err != nil {
        log.Println("Category name check error:", err)

        c.JSON(http.StatusInternalServerError, gin.H{
            "status":  false,
            "message": "failed to check category name",
        })
        return
    }

    if count != 0 {
        c.JSON(http.StatusBadRequest, gin.H{
            "status":  false,
            "message": "category name already exists",
        })
        return
    }

    // 5. Update text fields
    category.CategoryName = categoryName
    category.Description = description

    // 6. Check whether a new image was uploaded
    fileHeader, err := c.FormFile("image")

    if err == nil {
        // New image uploaded

        imageURL, err := helper.UploadToS3(fileHeader)

        if err != nil {
            log.Println("S3 Upload Error:", err)

            c.JSON(http.StatusInternalServerError, gin.H{
                "status":  false,
                "message": "failed to upload image to S3",
            })
            return
        }

        // Replace old image URL with new one
        category.ImageURL = imageURL
    }

    // If err != nil here, it simply means no new image was uploaded.
    // Therefore the existing ImageURL remains unchanged.

    // 7. Save changes
    if err := database.DB.Save(&category).Error; err != nil {
        log.Println("DB Update Error:", err)

        c.JSON(http.StatusInternalServerError, gin.H{
            "status":  false,
            "message": "database error while updating category",
        })
        return
    }

    // 8. Return updated category
    c.JSON(http.StatusOK, gin.H{
        "status":  true,
        "message": "Category Updated Successfully",
        "data":    category,
    })
}

func CategoryDelete(c *gin.Context) {
	fmt.Println("hello")
	CategoryID := c.Param("id")
	fmt.Println(CategoryID)
	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM categories WHERE id = ? AND deleted_at IS NULL`, CategoryID).Scan(&count)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "category id does not exist",
		})
		return
	}
	// var Category models.Category
	// err := c.BindJSON(&Category)
	// response := gin.H{
	// 	"status":  false,
	// 	"message": "failed to bind request",
	// }
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, response)
	// 	return
	// }
	database.DB.Where("id = ?", CategoryID).Delete(&models.Category{})
	database.DB.Model(&models.Product{}).Where("category_id = ?", CategoryID).Update("deleted_at", gorm.DeletedAt{Time: time.Now(), Valid: true})
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "category deleted succesfully"})
}

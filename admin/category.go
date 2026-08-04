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
	fmt.Println("is it here in REad Category")
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

	var Category models.Category
	err := c.BindJSON(&Category)
	response := gin.H{
		"status":  false,
		"message": "failed to bind request",
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// Validate the content of the JSON
	if err := helper.Validate(Category); err != nil {
		fmt.Println("", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"message":    err.Error(),
			"error_code": http.StatusBadRequest,
		})
		return
	}

	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM categories where category_name = ? AND deleted_at IS NULL`, Category.CategoryName).Scan(&count)
	if count != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "category name already exists",
		})
		return
	}
	category := models.Category{
		CategoryName: Category.CategoryName,
		Description:  Category.Description,
		ImageURL:     Category.ImageURL,
	}

	err = database.DB.Create(&category).Error
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "database error while adding category"})
	}
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "Category Added"})

}

func EditCategory(c *gin.Context) {

	CategoryID := c.Param("id")
	var category requestmodemodels.Category
	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM categories WHERE id = ? AND deleted_at IS NULL`, CategoryID).Scan(&count)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "category id does not exist",
		})
		return
	}
	var Category models.Category
	err := c.BindJSON(&Category)
	response := gin.H{
		"status":  false,
		"message": "failed to bind request",
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, response)
		return
	}
	if err := helper.Validate(Category); err != nil {
		fmt.Println("", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"message":    err.Error(),
			"error_code": http.StatusBadRequest,
		})
		return
	}
	category = requestmodemodels.Category{
		CategoryName: Category.CategoryName,
		Description:  Category.Description,
		ImageUrl:     Category.ImageURL,
	}
	database.DB.Model(&models.Category{}).Where("id = ?", CategoryID).Updates(&category)
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "Category Updated Successfully"})
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

package admin

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/helper"
	"github.com/Ansalps/GeZOne/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Category(c *gin.Context) {
	listorder := c.Query("list_order")
	var category []models.Category
	//tx := database.DB.Find(&category)
	sql := `SELECT * FROM categories WHERE deleted_at IS NULL`
	if listorder == "" || listorder == "ASC" {
		sql += ` ORDER BY categories.id ASC`
	} else if listorder == "DSC" {
		sql += ` ORDER BY categories.id DESC`
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

func CategoryAdd(c *gin.Context) {
	fmt.Println("hello")
	var Category models.CategoryEdit
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
	// var category models.Category
	// tx := database.DB.Where("category_name = ?", Category.CategoryName).Find(&category)
	// fmt.Println("", category)
	// fmt.Println("--", Category.CategoryName)
	// if tx.Error == nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"message": "category name already exists",
	// 	})
	// 	return
	// }
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
		ImageUrl:     Category.ImageUrl,
	}

	database.DB.Create(&category)
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "Category Added"})

}

func CategoryEdit(c *gin.Context) {
	fmt.Println("hello")
	CategoryID := c.Param("id")
	var category models.Category
	fmt.Println(CategoryID)
	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM categories WHERE id = ? AND deleted_at IS NULL`, CategoryID).Scan(&count)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "category id does not exist",
		})
		return
	}
	var Category models.CategoryEdit
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
	category = models.Category{
		CategoryName: Category.CategoryName,
		Description:  Category.Description,
		ImageUrl:     Category.ImageUrl,
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

package user

import (
	"fmt"
	"net/http"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/helper"
	"github.com/Ansalps/GeZOne/middleware"
	"github.com/Ansalps/GeZOne/models"
	"github.com/Ansalps/GeZOne/requestmodels"
	"github.com/Ansalps/GeZOne/responsemodels"
	"github.com/gin-gonic/gin"
)

func Profile(c *gin.Context) {
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
	var User responsemodels.User
	if err := database.DB.Where("id = ?", userID).First(&User).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": false, "message": "user not found"})
		return
	}
	if User.ProfilePicture != "" {
		if resolvedURL, err := helper.ResolveObjectURL(c.Request.Context(), "PRIVATE_S3_BUCKET_NAME", User.ProfilePicture); err == nil {
			User.ProfilePicture = resolvedURL
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "successfully retrieved user informations",
		"data": gin.H{
			"User": User,
		},
	})
}

func ProfilePicture(c *gin.Context) {
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
	var user models.User
	if err := database.DB.Select("profile_picture").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "profile picture not found",
		})
		return
	}

	profilePicture := user.ProfilePicture
	if profilePicture != "" {
		if resolvedURL, err := helper.ResolveObjectURL(c.Request.Context(), "PRIVATE_S3_BUCKET_NAME", profilePicture); err == nil {
			profilePicture = resolvedURL
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "profile picture retrieved successfully",
		"data": gin.H{
			"profile_picture": profilePicture,
		},
	})
}

func ProfilePictureUpload(c *gin.Context) {
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
	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "image file is required",
		})
		return
	}

	imageKey, err := helper.UploadToS3(c.Request.Context(), fileHeader, "profile-pictures", "PRIVATE_S3_BUCKET_NAME")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "failed to upload profile picture",
			"error":   err.Error(),
		})
		return
	}

	updateData := models.User{ProfilePicture: imageKey}
	if err := database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(&updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "failed to save profile picture",
			"error":   err.Error(),
		})
		return
	}

	var updatedUser models.User
	if err := database.DB.First(&updatedUser, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "failed to fetch updated profile",
		})
		return
	}

	resolvedProfilePicture := updatedUser.ProfilePicture
	if resolvedProfilePicture != "" {
		if resolvedURL, err := helper.ResolveObjectURL(c.Request.Context(), "PRIVATE_S3_BUCKET_NAME", resolvedProfilePicture); err == nil {
			resolvedProfilePicture = resolvedURL
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "profile picture updated successfully",
		"data": gin.H{
			"profile_picture": resolvedProfilePicture,
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

	var Profile requestmodels.ProfileEdit
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

	var PasswordChange requestmodels.PasswordChange
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

func Wishlist(c *gin.Context) {
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

	var wishlist []responsemodels.Wishlist
	database.DB.Raw(`SELECT wishlists.id,wishlists.created_at,wishlists.updated_at,wishlists.user_id,wishlists.product_id,products.product_name,categories.category_name,products.description,products.image_url,products.price,COALESCE(SUM(product_variants.stock), 0) AS stock,products.popular,COALESCE(string_agg(DISTINCT product_variants.size, ',' ORDER BY product_variants.size), '') AS size FROM wishlists JOIN products ON wishlists.product_id = products.id JOIN categories ON categories.id = products.category_id LEFT JOIN product_variants ON product_variants.product_id = products.id WHERE wishlists.user_id = ? GROUP BY wishlists.id, wishlists.created_at, wishlists.updated_at, wishlists.user_id, wishlists.product_id, products.product_name, categories.category_name, products.description, products.image_url, products.price, products.popular`, userID).Scan(&wishlist)
	c.JSON(http.StatusOK, gin.H{
		"data":    wishlist,
		"message": "succesfully shown wishlist",
	})
}

func WishlistAdd(c *gin.Context) {
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

	var wishlistadd requestmodels.WishlistAdd
	err := c.BindJSON(&wishlistadd)
	response := gin.H{
		"status":  false,
		"message": "failed to bind request",
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// Validate the content of the JSON
	if err := helper.Validate(wishlistadd); err != nil {
		fmt.Println("", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"message":    err.Error(),
			"error_code": http.StatusBadRequest,
		})
		return
	}
	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM wishlists WHERE product_id = ?`, wishlistadd.ProductID).Scan(&count)
	if count != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "product already added in wishlist",
		})
		return
	}
	wishlist := models.Wishlist{
		UserID:    userID,
		ProductID: wishlistadd.ProductID,
	}
	database.DB.Create(&wishlist)
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "Wishlist Added"})
}

func WishlistRemove(c *gin.Context) {
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

	var wishlistadd requestmodels.WishlistAdd
	err := c.BindJSON(&wishlistadd)
	response := gin.H{
		"status":  false,
		"message": "failed to bind request",
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// Validate the content of the JSON
	if err := helper.Validate(wishlistadd); err != nil {
		fmt.Println("", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"message":    err.Error(),
			"error_code": http.StatusBadRequest,
		})
		return
	}
	var count int64
	database.DB.Raw(`SELECT COUNT(*) FROM wishlists WHERE product_id = ? AND user_id = ?`, wishlistadd.ProductID, userID).Scan(&count)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "product does not exist wishlist for this particular user",
		})
		return
	}
	database.DB.Where("product_id = ? AND user_id = ?", wishlistadd.ProductID, userID).Delete(&models.Wishlist{})
	c.JSON(http.StatusOK, gin.H{
		"message": "product removed from wishlist successfully",
	})
}

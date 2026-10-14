package user

import (
	"net/http"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/middleware"
	"github.com/Ansalps/GeZOne/models"
	"github.com/Ansalps/GeZOne/responsemodels"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func WalletListing(c *gin.Context) {
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
	var wallet responsemodels.Wallet
	database.DB.Raw(`SELECT * from wallets where user_id = ?`, userID).Scan(&wallet)
	c.JSON(http.StatusOK, gin.H{
		"data":    wallet,
		"message": "listing wallet",
	})
}

func WalletTransactionListing(c *gin.Context) {
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
	var wallettransaction []responsemodels.WalletTransaction
	database.DB.Raw(`SELECT * from wallet_transactions where user_id = ?`, userID).Scan(&wallettransaction)
	c.JSON(http.StatusOK, gin.H{
		"data":    wallettransaction,
		"message": "listing wallettransactions",
	})

}

func AddMoney(c *gin.Context) {
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
	var payload struct {
		Amount float64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "invalid request payload",
		})
		return
	}

	if payload.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "amount must be greater than zero",
		})
		return
	}

	var wallet models.Wallet
	if err := database.DB.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			wallet = models.Wallet{UserID: userID, Balance: 0}
			if createErr := database.DB.Create(&wallet).Error; createErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  false,
					"message": "failed to initialize wallet",
				})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "failed to load wallet",
			})
			return
		}
	}

	newBalance := wallet.Balance + payload.Amount
	if err := database.DB.Model(&wallet).Update("balance", newBalance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "failed to update wallet balance",
		})
		return
	}

	transaction := models.WalletTransaction{
		UserID:          userID,
		Amount:          payload.Amount,
		TransactionType: "Credit",
		Description:     "Wallet top-up",
	}
	if err := database.DB.Create(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "wallet updated but transaction log failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "money added to wallet successfully",
		"data": gin.H{
			"balance": newBalance,
		},
	})
}

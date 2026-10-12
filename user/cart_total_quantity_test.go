package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/middleware"
	"github.com/Ansalps/GeZOne/models"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestTotalQuantity_EmptyCartReturnsZero(t *testing.T) {
	gin.SetMode(gin.TestMode)

	if database.DB == nil {
		database.Initialize()
		database.AutoMigrate()
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/cart-total-quantity", nil)
	c.Set("claims", &middleware.CustomClaims{ID: 999999})

	TotalQuantity(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d with body: %s", w.Code, w.Body.String())
	}

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}

	if body["status"] != true {
		t.Fatalf("expected status true, got %#v", body["status"])
	}

	if total, ok := body["total_quantity"]; !ok || total != float64(0) {
		t.Fatalf("expected total_quantity to be 0, got %#v", body["total_quantity"])
	}
}

func TestGetUserCartSummary_UsesCartJoinAndCurrentSchema(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.Cart{}, &models.CartItem{}, &models.Product{}, &models.Category{}); err != nil {
		t.Fatalf("failed to migrate sqlite schema: %v", err)
	}

	category := models.Category{CategoryName: "Shoes", Description: "desc", ImageURL: "img"}
	if err := db.Create(&category).Error; err != nil {
		t.Fatalf("failed to create category: %v", err)
	}
	product := models.Product{CategoryID: category.ID, ProductName: "Sneaker", Description: "desc", ImageURL: "img", Price: 120}
	if err := db.Create(&product).Error; err != nil {
		t.Fatalf("failed to create product: %v", err)
	}
	user := models.User{FirstName: "A", LastName: "B", Email: "a@example.com", Password: "pw", Phone: "1234567890", Status: "Active"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	cart := models.Cart{UserID: user.ID}
	if err := db.Create(&cart).Error; err != nil {
		t.Fatalf("failed to create cart: %v", err)
	}
	if err := db.Create(&models.CartItem{CartID: cart.ID, ProductID: product.ID, Size: "M", Quantity: 2, UnitPrice: 80}).Error; err != nil {
		t.Fatalf("failed to create cart item 1: %v", err)
	}
	if err := db.Create(&models.CartItem{CartID: cart.ID, ProductID: product.ID, Size: "L", Quantity: 1, UnitPrice: 50}).Error; err != nil {
		t.Fatalf("failed to create cart item 2: %v", err)
	}

	count, totalQty, totalAmount, err := getUserCartSummary(db, user.ID)
	if err != nil {
		t.Fatalf("getUserCartSummary returned error: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected 2 cart items, got %d", count)
	}
	if totalQty != 3 {
		t.Fatalf("expected total quantity 3, got %d", totalQty)
	}
	if totalAmount != 210 {
		t.Fatalf("expected total amount 210, got %v", totalAmount)
	}
}

package main

import (
	"time"

	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/route"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	database.Initialize()
	database.AutoMigrate()
	godotenv.Load(".env")
}

func main() {

	router := gin.Default()
	// Use default CORS middleware (allows all origins)
	// Configure CORS explicitly for cookies
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Must be explicit, CANNOT be "*"
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // 👈 CRITICAL: Allows browser to store cookies!
		MaxAge:           24 * time.Hour,
	}))
	route.RegisterUrls(router)
	router.LoadHTMLGlob("templates/*")
	router.Run(":8080")

}

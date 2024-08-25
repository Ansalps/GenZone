package main

import (
	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/route"
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
	route.RegisterUrls(router)
	router.LoadHTMLGlob("templates/*")
	router.Run(":8080")

}

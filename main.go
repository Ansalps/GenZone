package main

import (
	"github.com/Ansalps/GeZOne/database"
	"github.com/Ansalps/GeZOne/route"
	"github.com/gin-gonic/gin"
)

func init() {
	database.Initialize()
	database.AutoMigrate()
}

func main() {
	router := gin.Default()
	route.RegisterUrls(router)
	router.Run(":8080")
}

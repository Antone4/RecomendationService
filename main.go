package main

import (
	"recommendation-service/internal/database"
	"recommendation-service/internal/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDB()

	r := gin.Default()
	r.GET("/recommend", handler.RecommendHandler)
	r.Run(":8080")
}

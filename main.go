package main

import (
	"fmt"
	"main/global"
	"main/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	if er := godotenv.Load(); er != nil {
		fmt.Println("Missing .env file")
	}
}
func main() {
	router := gin.Default()
	router.Use(global.RateLimitMiddleware())
	router.POST("/shorten", services.GenerateShortUrlEndpoint)
	router.GET("/shorten/:short_url", services.GetOriginalUrlEndpoint)
	router.DELETE("/shorten/:short_url", services.DeleteShortUrlEndpoint)
	router.Run(":8080")
}

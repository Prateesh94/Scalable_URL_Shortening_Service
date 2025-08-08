package main

import (
	"fmt"
	"main/global"
	"main/services"
	"os"

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
	run := os.Getenv("ENV")
	if run == "cloud" {
		router.Use(global.RateLimitMiddleware())
	}

	router.POST("/shorten", services.GenerateShortUrlEndpoint)
	router.GET("/shorten/:short_url", services.GetOriginalUrlEndpoint)
	router.GET("/shorten/:short_url/count", services.GetCountEndpoint)
	router.DELETE("/shorten/:short_url", services.DeleteShortUrlEndpoint)
	router.Run(":8080")
}

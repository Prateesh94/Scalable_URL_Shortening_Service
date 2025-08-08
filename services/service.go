package services

import (
	"main/global"
	"main/internal"
	"main/internal/database"

	"github.com/gin-gonic/gin"
)

func GenerateShortUrlEndpoint(c *gin.Context) {
	var data global.Data
	var userdata global.Data
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	short_url := internal.GenerateShortUrl(data.OriginalUrl)
	if short_url == "" {
		c.JSON(500, gin.H{"error": "Failed to generate short URL, URL may already exist"})
		return
	}
	userdata, er := database.AddData(short_url, data.OriginalUrl)
	if er != nil {
		c.JSON(500, gin.H{"error": "Failed to add data to the database"})
		return
	}
	data.ShortUrl = short_url
	c.JSON(200, userdata)

}

func GetOriginalUrlEndpoint(c *gin.Context) {
	short_url := c.Param("short_url")
	if short_url == "" {
		c.JSON(400, gin.H{"error": "Short URL is required"})
		return
	}
	original_url, err := database.RetrieveOriginalURL(short_url)
	if err != nil {
		c.JSON(500, gin.H{"error": "The short URL does not exist"})
		return
	}
	c.JSON(200, gin.H{"original_url": original_url})

}

func DeleteShortUrlEndpoint(c *gin.Context) {
	short_url := c.Param("short_url")
	if short_url == "" {
		c.JSON(400, gin.H{"error": "Short URL is required"})
		return
	}
	err := database.DeleteCacheUrl(short_url)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete short URL, it may not exist"})
		return
	}
	c.JSON(200, gin.H{"message": "Short URL deleted successfully"})
}

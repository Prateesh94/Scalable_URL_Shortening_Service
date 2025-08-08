package services

import (
	"os"
	"url-shortener/global"
	"url-shortener/internal"
	"url-shortener/internal/database"

	"github.com/gin-gonic/gin"
)

var runstate = os.Getenv("ENV")

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
	var original_url string
	var err error
	if short_url == "" {
		c.JSON(400, gin.H{"error": "Short URL is required"})
		return
	}
	if runstate == "cloud" {
		original_url, err = database.RetrieveOriginalURL(short_url)
	} else {
		original_url, err = database.GetOriginalUrl(short_url)
	}
	if err != nil {
		c.JSON(500, gin.H{"error": "The short URL does not exist"})
		return
	}
	c.JSON(200, gin.H{"original_url": original_url})

}

func DeleteShortUrlEndpoint(c *gin.Context) {
	short_url := c.Param("short_url")
	var err error
	if short_url == "" {
		c.JSON(400, gin.H{"error": "Short URL is required"})
		return
	}
	if runstate == "cloud" {
		err = database.DeleteCacheUrl(short_url)
	} else {
		err = database.DeleteShortUrl(short_url)
	}

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete short URL, it may not exist"})
		return
	}
	c.JSON(200, gin.H{"message": "Short URL deleted successfully"})
}

func GetCountEndpoint(c *gin.Context) {
	short_url := c.Param("short_url")
	if short_url == "" {
		c.JSON(400, gin.H{"error": "Short URL is required"})
		return
	}
	count, err := database.GetCount(short_url)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve count for the short URL"})
		return
	}
	c.JSON(200, gin.H{"url": short_url, "count": count})
}

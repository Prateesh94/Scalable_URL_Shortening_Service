package internal

import (
	"crypto/sha256"
	"encoding/base64"
	"main/internal/database"
)

func GenerateShortUrl(original_url string) string {

	hash := sha256.Sum256([]byte(original_url))
	random := base64.URLEncoding.EncodeToString(hash[:])
	//fmt.Println("Generated hash:", random)
	shorturl := random[:8] // Take the first 8 characters of the base64 encoded string
	//fmt.Println("Short URL:", shorturl)
	check := database.CheckUrl(shorturl)
	if check {
		return ""
	}
	//
	return shorturl

}

package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

// Cache for Redis
var ctx = context.Background()
var rdb *redis.Client

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"), // default to localhost if not set
		Password: os.Getenv("REDIS_PWD"), // no password set for local dev
		DB:       0,                      // use default DB
	})
}

func RetrieveOriginalURL(shortURL string) (string, error) {
	// 1. Try to get the original URL from the cache
	originalURL, err := rdb.Get(ctx, shortURL).Result()
	if err == nil {
		// Cache hit! 🎉
		UpdateCount(shortURL) // Update the count in the database
		return originalURL, nil
	}

	// If err is redis.Nil, the key wasn't in the cache.
	// If err is any other error, log it but continue to the database.
	if err != redis.Nil {
		//log.Printf("Redis error: %v", err)
	}

	// 2. Cache miss, so query the database
	dbOriginalURL, err := GetOriginalUrl(shortURL)
	if err != nil {
		return "", err
	}

	// 3. Store the result in the cache for future requests
	// We set a TTL of 1 hour to keep the data fresh.
	ttl := 30 * time.Minute
	err = rdb.Set(ctx, shortURL, dbOriginalURL, ttl).Err()
	if err != nil {
		log.Printf("Failed to set cache for %s: %v", shortURL, err)
	}

	return dbOriginalURL, nil
}

func DeleteCacheUrl(shortURL string) error {
	// 1. Delete from the cache
	err := rdb.Del(ctx, shortURL).Err()
	if err != nil {
		log.Printf("Failed to delete cache for %s: %v", shortURL, err)
	}

	// 2. Delete from the database
	return DeleteShortUrl(shortURL)
}

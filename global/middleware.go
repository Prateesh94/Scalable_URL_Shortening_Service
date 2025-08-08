package global

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var rdb *redis.Client

func init() {
	// Initialize Redis client once at application startup.
	// In a real app, this would be configured from an env variable.
	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"), // default to localhost if not set
		Password: os.Getenv("REDIS_PWD"), // no password set for local dev
		DB:       0,                      // use default DB
	})
}
func isRateLimited(key string, limit int, window time.Duration) (bool, error) {
	// Use the user's IP or a user ID as the Redis key.
	redisKey := fmt.Sprintf("rate_limit:%s", key)
	now := time.Now().UnixMilli()

	pipe := rdb.Pipeline()
	// Add the current timestamp to the sorted set.
	pipe.ZAdd(ctx, redisKey, &redis.Z{Score: float64(now), Member: strconv.FormatInt(now, 10)})
	// Remove all timestamps older than the time window.
	pipe.ZRemRangeByScore(ctx, redisKey, "0", strconv.FormatInt(now-window.Milliseconds(), 10))
	// Count the number of remaining elements (requests in the window).
	pipe.ZCard(ctx, redisKey)
	// Set a short expiration on the key.
	pipe.Expire(ctx, redisKey, window)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	count, err := pipe.ZCard(ctx, redisKey).Result()
	if err != nil {
		return false, err
	}

	return count > int64(limit), nil
}

// rateLimitMiddleware is a reusable middleware function that checks the rate limit.
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use the client's IP address as the key for rate limiting.
		ip := c.ClientIP()
		limit := 5
		window := 30 * time.Second

		isLimited, err := isRateLimited(ip, limit, window)
		if err != nil {
			log.Printf("Rate limiter check failed: %v", err)
			// Fail open: allow the request if the rate limiter is unavailable.
			c.Next()
			return
		}

		if isLimited {
			// If limited, abort the request and return a 429 status code.
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			return
		}

		// If not limited, proceed to the next middleware or handler.
		c.Next()
	}
}

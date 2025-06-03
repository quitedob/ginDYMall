package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"douyin/pkg/utils/response" // Adjust if your response package is different
	"douyin/mylog"             // Adjust if your log package is different
)

// RateLimitMiddleware provides rate limiting functionality based on a key (user ID or IP) and path.
func RateLimitMiddleware(rdb *redis.Client, prefix string, limit int64, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if rdb == nil {
			// Assuming mylog is initialized and available. If not, use standard log.
			// For example: log.Printf("Error: Redis client (rdb) is nil...")
			mylog.Error("Redis client (rdb) is nil in RateLimitMiddleware. Rate limiting disabled for this request.")
			c.Next() // Allow request if Redis is not configured/available for rate limiting
			return
		}

		var key string
		// Try to get userID from context (assuming it's set by AuthMiddleware as uint)
		userID, userExists := c.GetUint("userID") // gin.Context has GetUint

		if userExists {
			key = fmt.Sprintf("%s:user:%d:path:%s", prefix, userID, c.FullPath())
		} else {
			ip := c.ClientIP()
			key = fmt.Sprintf("%s:ip:%s:path:%s", prefix, ip, c.FullPath())
		}

		ctx := context.Background() // Or c.Request.Context() if appropriate for Redis operations

		// INCR command to increment the counter for the key
		count, err := rdb.Incr(ctx, key).Result()
		if err != nil {
			mylog.Errorf("Redis INCR failed for rate limiting key %s: %v", key, err)
			// Allowing request if Redis fails, to prevent blocking users due to Redis issues.
			c.Next()
			return
		}

		// If it's the first request for this key in the current window, set expiration
		if count == 1 {
			if err := rdb.Expire(ctx, key, window).Err(); err != nil {
				mylog.Errorf("Redis EXPIRE failed for rate limiting key %s: %v", key, err)
				// If EXPIRE fails, the key might persist indefinitely. Consider deleting or logging.
			}
		}

		// Check if the count exceeds the limit
		if count > limit {
			// Try to get TTL to set Retry-After header more accurately
			ttl, ttlErr := rdb.TTL(ctx, key).Result()
			if ttlErr == nil && ttl > 0 {
				c.Header("Retry-After", strconv.Itoa(int(ttl.Seconds())))
			} else {
				// Fallback to window if TTL fails or key has no expiration (should not happen if EXPIRE worked)
				c.Header("Retry-After", strconv.Itoa(int(window.Seconds())))
				if ttlErr != nil && ttlErr != redis.Nil { // Log if TTL error is not just "key has no TTL"
					mylog.Warnf("Redis TTL failed for rate limiting key %s: %v. Using default window for Retry-After.", key, ttlErr)
				}
			}
			
			response.Fail(c, http.StatusTooManyRequests, "请求过于频繁，请稍后再试 (Too Many Requests, please try again later)")
			c.Abort()
			return
		}

		c.Next()
	}
}

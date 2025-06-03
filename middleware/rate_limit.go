package middleware

import (
	"context"
	"douyin/pkg/utils/response" // Assuming this is your standardized response package
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8" // Using go-redis/redis
)

// RateLimitMiddleware creates a middleware for rate limiting requests.
func RateLimitMiddleware(rdb *redis.Client, prefix string, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context() // Use request context for Redis operations

		// Attempt to get user ID from context (set by AuthMiddleware, if applicable)
		userIDVal, userExists := c.Get("user_id")
		var keyIdentifier string

		if userExists {
			// If user is authenticated, use user ID for rate limiting
			userID, ok := userIDVal.(uint) // Assuming user_id is uint
			if !ok {
				// Fallback or error if type assertion fails, though this shouldn't happen if set correctly
				keyIdentifier = c.ClientIP() // Fallback to IP
			} else {
				keyIdentifier = strconv.FormatUint(uint64(userID), 10)
			}
		} else {
			// If user is not authenticated, use client IP
			keyIdentifier = c.ClientIP()
		}

		// Construct the Redis key
		// Include path to make the rate limit path-specific
		// You might want to normalize the path or use c.FullPath() if parameters are involved
		rateLimitKey := fmt.Sprintf("rate_limit:%s:%s:%s", prefix, keyIdentifier, c.FullPath())

		// Use Redis pipeline for atomic INCR and EXPIRE
		pipe := rdb.Pipeline()
		countCmd := pipe.Incr(ctx, rateLimitKey)
		pipe.Expire(ctx, rateLimitKey, window) // Set/update expiration window
		_, err := pipe.Exec(ctx)

		if err != nil {
			// If Redis fails, decide whether to allow or deny the request.
			// Allowing might be safer to not block users due to Redis issues.
			// Log the error.
			// customlog.LogrusObj.Errorf("Rate limiter Redis error: %v", err)
			fmt.Printf("Rate limiter Redis error: %v for key %s\n", err, rateLimitKey) // Placeholder log
			c.Next()
			return
		}

		count := countCmd.Val()

		if count > int64(limit) {
			// Limit exceeded
			// Set a Retry-After header if appropriate (value in seconds)
			// c.Header("Retry-After", strconv.FormatInt(int64(window.Seconds()), 10))

			// Using your standardized response if available
			// Assuming response.Fail is defined as: func Fail(code int, msg string) APIResponse
			// And APIResponse struct has Code, Message, Data fields.
			// The business code for rate limiting could be specific, e.g., 429 or a custom one.
			apiResp := response.Fail(http.StatusTooManyRequests, "Too many requests. Please try again later.")
			c.JSON(http.StatusTooManyRequests, apiResp)
			c.Abort()
			return
		}

		c.Next()
	}
}

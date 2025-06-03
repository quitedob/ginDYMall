package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

// ContextTimeout returns a middleware that sets a timeout on the request context.
func ContextTimeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a new context with a timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel() // Ensure cancel is called to release resources

		// Replace the request's context with the new timed context
		c.Request = c.Request.WithContext(ctx)

		// Call the next handler
		c.Next()
	}
}

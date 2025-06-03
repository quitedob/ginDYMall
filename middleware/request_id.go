package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestID is a middleware that injects a unique request ID into the context and response header.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate a new UUID
		rid := uuid.New().String()

		// Set it in the Gin context
		c.Set("RequestID", rid)

		// Set it in the response header
		c.Writer.Header().Set("X-Request-ID", rid)

		// Continue to the next handler
		c.Next()
	}
}

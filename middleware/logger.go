package middleware

import (
	"time"

	"douyin/pkg/utils/log" // Assuming your logger package path

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// LoggerMiddleware is a custom Gin middleware for logging requests using Logrus.
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now() // Start timer

		// Process request
		c.Next()

		// After request
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		requestID := c.GetString("RequestID") // Get RequestID from context (set by RequestID middleware)
		// Get error message if any
		errMsg := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// Prepare log fields
		fields := logrus.Fields{
			"statusCode": statusCode,
			"latency":    latency.String(), // Log latency as string for readability
			"clientIP":   clientIP,
			"method":     method,
			"path":       path,
			"requestID":  requestID,
		}

		// Create a log entry with these fields
		entry := log.LogrusObj.WithFields(fields) // Use your global LogrusObj

		if errMsg != "" {
			entry.Error(errMsg) // Log as error if c.Errors has content
		} else {
			if statusCode >= 500 {
				entry.Error("Server error")
			} else if statusCode >= 400 {
				entry.Warn("Client error")
			} else {
				entry.Info("Request completed")
			}
		}
	}
}

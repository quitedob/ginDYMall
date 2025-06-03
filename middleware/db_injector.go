package middleware

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"douyin/global" // Assuming global.DB is where the DB instance is stored
)

// DBInjectorMiddleware injects the GORM DB instance into the Gin context.
func DBInjectorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if global.DB == nil {
			// Log error or handle, but for now, we proceed. RBAC will fail if DB is nil.
			// Consider aborting if DB is critical and not available.
			// For example, you could log and abort:
			// log.Error("DB instance is nil in DBInjectorMiddleware")
			// c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database not available"})
			// return
		}
		c.Set("db", global.DB) // Set the DB instance from global variable
		c.Next()
	}
}

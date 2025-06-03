package middleware

import (
	"github.com/gin-gonic/gin"
	i18nUtils "douyin/pkg/utils/i18n" // Your i18n utility package
)

// I18nMiddleware sets the i18n.Localizer in the Gin context.
func I18nMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		loc := i18nUtils.GetLocalizer(c)
		c.Set("localizer", loc) // Key used in i18nUtils.MustLocalize
		c.Next()
	}
}

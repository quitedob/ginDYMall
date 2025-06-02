// middleware/jwt.go
package middleware

import (
	"douyin/consts" // Use our own consts
	"douyin/pkg/utils/jwt"
	"douyin/pkg/utils/log"
	"errors" // For creating error instances
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	// "time" // time.Now().Unix() is not needed here if ParseAccessToken handles expiry checks
)

// AuthMiddleware 返回一个 Gin 中间件，用于统一拦截并校验 JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(consts.HeaderAuthorization) // Use const for "Authorization"
		log.Infof("AuthMiddleware 收到 Header: %s", authHeader)

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少 Authorization 头部"})
			c.Abort()
			return
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization 格式无效，需 'Bearer <token>'"})
			c.Abort()
			return
		}
		accessTokenString := strings.TrimPrefix(authHeader, bearerPrefix)

		claims, err := jwt.ParseAccessToken(accessTokenString)
		if err != nil {
			log.Warnf("解析访问令牌失败: %v", err)
			// Check if the error indicates an expired token
			// The error message "令牌已过期或未生效" is one of the custom messages from ParseAccessToken
			if err.Error() == "令牌已过期或未生效" { // This condition might need to be more robust
				refreshTokenString := c.GetHeader(consts.HeaderRefreshToken)
				if refreshTokenString == "" {
					log.Warn("访问令牌已过期，但未找到刷新令牌")
					c.JSON(http.StatusUnauthorized, gin.H{"error": "访问令牌已过期，请重新登录"})
					c.Abort()
					return
				}

				log.Info("尝试使用刷新令牌续期")
				refreshClaims, refreshErr := jwt.ParseRefreshToken(refreshTokenString)
				if refreshErr != nil {
					log.Warnf("解析刷新令牌失败: %v", refreshErr)
					c.JSON(http.StatusUnauthorized, gin.H{"error": "刷新令牌无效或已过期，请重新登录"})
					c.Abort()
					return
				}

				// Refresh token is valid, generate new tokens
				newAccessToken, newRefreshToken, genErr := jwt.GenerateToken(refreshClaims.UserId, refreshClaims.Username)
				if genErr != nil {
					log.Errorf("生成新令牌失败: %v", genErr)
					_ = c.Error(errors.New("无法续期会话，请稍后重试")) // Pass to global error handler
					c.Abort()
					return
				}

				// Set new tokens in response headers
				c.Header(consts.HeaderNewAccessToken, newAccessToken)
				c.Header(consts.HeaderNewRefreshToken, newRefreshToken)
				log.Info("令牌已刷新，新令牌已在响应头中设置")

				// Set user info in context for the current request
				c.Set("user_id", refreshClaims.UserId)
				c.Set("user_name", refreshClaims.Username)
				c.Next() // Continue processing the request with new context
				return
			}

			// Access token is invalid for reasons other than expiry
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
			c.Abort()
			return
		}

		// Access token is valid and not expired
		c.Set("user_id", claims.UserId)
		c.Set("user_name", claims.Username)
		c.Next()
	}
}

// SetToken 将访问令牌和刷新令牌设置到响应头和Cookie中
// This function might be used by login/register handlers.
// Note: consts.AccessTokenHeader and consts.RefreshTokenHeader from the original code
// might be different from our consts.HeaderAuthorization or consts.HeaderRefreshToken.
// For now, this function is left as is, but its usage should be reviewed
// in context of where it's called. The new headers for refreshed tokens are
// consts.HeaderNewAccessToken and consts.HeaderNewRefreshToken.
func SetToken(c *gin.Context, accessToken, refreshToken string) {
	secure := IsHttps(c)
	// Assuming the original consts were for setting initial tokens during login.
	// If these are header names, they should align with how tokens are expected by clients.
	// For example, Authorization for access token, X-Refresh-Token for refresh.
	c.Header(consts.HeaderAuthorization, "Bearer "+accessToken) // Standard way to send access token
	c.Header(consts.HeaderRefreshToken, refreshToken)

	// Cookie setting might be an alternative or complementary way to handle tokens.
	// MaxAge here seems to be from a different consts package.
	// Defaulting to AccessTokenExpireDuration for access token cookie.
	c.SetCookie(consts.HeaderAuthorization, accessToken, int(consts.AccessTokenExpireDuration.Seconds()), "/", "", secure, true)
	c.SetCookie(consts.HeaderRefreshToken, refreshToken, int(consts.RefreshTokenExpireDuration.Seconds()), "/", "", secure, true)
}

// IsHttps 判断请求是否使用HTTPS协议
func IsHttps(c *gin.Context) bool {
	// HeaderForwardedProto might also need to be a const if used commonly.
	return c.GetHeader("X-Forwarded-Proto") == "https" || c.Request.TLS != nil
}

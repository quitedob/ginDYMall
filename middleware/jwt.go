// middleware/jwt.go
package middleware

import (
	"douyin/pkg/utils/jwt"
	"douyin/pkg/utils/log"
	"github.com/CocaineCong/gin-mall/consts"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

// AuthMiddleware 返回一个 Gin 中间件，用于统一拦截并校验 JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1) 从请求头中读取 Authorization
		//    要求格式为：Authorization: Bearer xxxxxxx
		authHeader := c.GetHeader("Authorization")
		log.Infof("AuthMiddleware 收到 Header: %s", authHeader)

		if authHeader == "" {

			// 这里为了演示，直接用 401
			// 你也可以用 200 + 自定义json

			c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少 Authorization 头部"})
			c.Abort()
			return
		}

		// 2) 去掉 "Bearer " 前缀，得到真正的 token
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization 格式无效，需 'Bearer <token>'"})
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, bearerPrefix)

		// 3) 调用你写好的 ParseToken 校验签名
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			// 如果签名无效、token 已被吊销等，就会返回错误
			log.Errorf("Token:%s", tokenString)
			log.Errorf("JWT111 校验失败: %s", err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效或过期的令牌"})
			c.Abort()
			return
		}

		// 4) 检查过期时间
		if claims.ExpiresAt < time.Now().Unix() {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "令牌已过期"})
			c.Abort()
			return
		}

		// 5) 如果一切正常，就放行
		//    你可以把用户信息塞进 c，用于后续 Handler 获取
		//    例如：c.Set("user_id", claims.UserId)
		c.Set("user_id", claims.UserId)
		c.Set("user_name", claims.Username)

		// 放行，让后续 Handler 正常执行
		c.Next()
	}
}

// SetToken 将访问令牌和刷新令牌设置到响应头和Cookie中
func SetToken(c *gin.Context, accessToken, refreshToken string) {
	secure := IsHttps(c)
	c.Header(consts.AccessTokenHeader, accessToken)
	c.Header(consts.RefreshTokenHeader, refreshToken)
	c.SetCookie(consts.AccessTokenHeader, accessToken, consts.MaxAge, "/", "", secure, true)
	c.SetCookie(consts.RefreshTokenHeader, refreshToken, consts.MaxAge, "/", "", secure, true)
}

// IsHttps 判断请求是否使用HTTPS协议
func IsHttps(c *gin.Context) bool {
	return c.GetHeader(consts.HeaderForwardedProto) == "https" || c.Request.TLS != nil
}

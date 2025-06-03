package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Cors 跨域中间件，用于处理跨域请求
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求方法和请求头中的 Origin 字段
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")

		// 如果存在 Origin，则允许跨域访问
		if origin != "" {
			// TODO: Make this list of allowed origins configurable
			allowedOrigins := []string{"https://production.example.com", "http://localhost:3000", "http://127.0.0.1:3000"} // Example origins
			isAllowed := false
			for _, allowed := range allowedOrigins {
				if origin == allowed {
					isAllowed = true
					break
				}
			}

			if isAllowed {
				c.Header("Access-Control-Allow-Origin", origin)
			} else {
				// Optionally, log the blocked origin or handle it specifically.
				// For now, if not allowed, the ACAO header will not be set for this origin,
				// and the browser will block the cross-origin request.
				// Alternatively, to be more explicit, you could set a default allowed origin if one exists,
				// or simply not set the header, which is the current behavior if !isAllowed.
				// If you want to explicitly deny, you might respond with an error, but that's not typical for ACAO.
			}

			// 允许访问的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			// 允许携带的请求头信息
			c.Header("Access-Control-Allow-Headers", strings.Join([]string{
				"Authorization", "Content-Length", "X-CSRF-Token", "Token", "session",
				"X_Requested_With", "Accept", "Origin", "Host", "Connection",
				"Accept-Encoding", "Accept-Language", "DNT", "X-CustomHeader", "Keep-Alive",
				"User-Agent", "If-Modified-Since", "Cache-Control", "Content-Type", "Pragma",
			}, ", "))
			// 允许浏览器解析的响应头信息
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type, Expires, Last-Modified, Pragma")
			// 设置预检请求的缓存时间，单位为秒
			c.Header("Access-Control-Max-Age", "172800")
			// 跨域请求是否允许携带 Cookie，根据实际情况设置（此处设置为 true 表示允许携带 Cookie）
			c.Header("Access-Control-Allow-Credentials", "true")
			// 设置响应内容类型为 JSON
			c.Set("content-type", "application/json")
		}

		// 对于 OPTIONS 请求，直接返回 200 状态，结束请求
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		// 继续处理请求
		c.Next()
	}
}

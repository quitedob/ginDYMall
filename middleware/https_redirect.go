package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HTTPSRedirectMiddleware redirects HTTP requests to HTTPS.
// It checks for TLS connection and X-Forwarded-Proto header (if behind a proxy).
// Note: In many production environments, TLS termination and HTTPS redirection
// are handled at the load balancer or reverse proxy level.
// This middleware is useful if the application server itself needs to handle redirection.
func HTTPSRedirectMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isHTTPS := false

		// Check if the connection is already secure (TLS is set)
		if c.Request.TLS != nil {
			isHTTPS = true
		}

		// Check X-Forwarded-Proto header, commonly set by reverse proxies / load balancers
		// when they terminate TLS and forward the request as HTTP to the app server.
		if !isHTTPS {
			forwardedProto := c.GetHeader("X-Forwarded-Proto")
			if forwardedProto == "https" {
				isHTTPS = true
			}
		}
		
		// Scheme header can also be checked, though X-Forwarded-Proto is more standard for this.
		// if !isHTTPS {
		// 	scheme := c.Request.Header.Get("Scheme") // Less common for this specific purpose
		// 	if scheme == "https" {
		// 		isHTTPS = true
		// 	}
		// }


		if !isHTTPS {
			// Construct the target URL with HTTPS scheme
			// c.Request.URL.Host is not available, use c.Request.Host which includes hostname:port
			targetURL := "https://" + c.Request.Host + c.Request.URL.String()
			
			// Perform a permanent redirect
			c.Redirect(http.StatusPermanentRedirect, targetURL)
			c.Abort() // Stop further processing for this request
			return
		}

		// If already HTTPS, proceed to the next handler
		c.Next()
	}
}

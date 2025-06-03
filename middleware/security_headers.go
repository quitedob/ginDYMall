package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeadersMiddleware adds common security-related HTTP headers to responses.
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevents clickjacking attacks
		c.Header("X-Frame-Options", "DENY")

		// Prevents browsers from MIME-sniffing a response away from the declared Content-Type
		c.Header("X-Content-Type-Options", "nosniff")

		// Enables the XSS filter built into most recent web browsers.
		// "1; mode=block" prevents rendering of the page if an attack is detected.
		c.Header("X-XSS-Protection", "1; mode=block")

		// A basic Content Security Policy (CSP) that allows resources (scripts, styles, images, etc.)
		// to be loaded only from the same origin ('self').
		// This is a starting point; a more complex application might need a more detailed CSP.
		c.Header("Content-Security-Policy", "default-src 'self'")

		// HTTP Strict Transport Security (HSTS)
		// Tells browsers to always connect to this site using HTTPS for the next N seconds.
		// Only uncomment and use this if your site is fully HTTPS capable and you intend to enforce it.
		// Ensure you understand the implications, especially `includeSubDomains` if you have subdomains
		// that might not be HTTPS ready.
		// c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains") // 1 year

		c.Next()
	}
}

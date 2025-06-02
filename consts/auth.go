package consts

import "time"

// Token durations
const (
	AccessTokenExpireDuration  = 15 * time.Minute
	RefreshTokenExpireDuration = 7 * 24 * time.Hour
)

// Header names for tokens
const (
	HeaderAuthorization   = "Authorization" // Standard header for access token (e.g., "Bearer <token>")
	HeaderRefreshToken    = "X-Refresh-Token"
	HeaderNewAccessToken  = "New-Access-Token"
	HeaderNewRefreshToken = "New-Refresh-Token"
)

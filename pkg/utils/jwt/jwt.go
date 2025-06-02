// 文件：pkg/utils/jwt/jwt.go
// 作用：实现 JWT 令牌的生成、解析、吊销，同时与 Redis 交互存储令牌

package jwt

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	"douyin/consts"
	"douyin/pkg/utils/log"
	"douyin/repository/cache"
)

// jwtSecret is the secret key for signing JWT tokens. It's initialized by Init().
var jwtSecret []byte

// Init sets the JWT secret key. This must be called before generating or parsing tokens.
func Init(secret string) {
	if secret == "" {
		log.Panic("JWT secret cannot be empty")
	}
	jwtSecret = []byte(secret)
	log.Info("JWT Secret initialized successfully.")
}

// Claims defines the JWT payload for access tokens.
type Claims struct {
	UserId   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// RefreshClaims defines the JWT payload for refresh tokens.
type RefreshClaims struct {
	UserId         uint   `json:"user_id"`
	Username       string `json:"username"`
	IsRefreshToken bool   `json:"is_refresh_token"`
	jwt.StandardClaims
}

// GenerateToken creates access_token, refresh_token.
// Access token is stored in Redis for potential revocation.
func GenerateToken(userId uint, username string) (string, string, error) {
	if len(jwtSecret) == 0 {
		log.Error("JWT secret is not initialized. Call jwt.Init() first.")
		return "", "", errors.New("JWT secret not initialized")
	}
	nowTime := time.Now()
	// Use consts for durations
	accessTokenExpireTime := nowTime.Add(consts.AccessTokenExpireDuration)
	refreshTokenExpireTime := nowTime.Add(consts.RefreshTokenExpireDuration)

	// Access Token Claims
	claims := Claims{
		UserId:   userId,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: accessTokenExpireTime.Unix(),
			Issuer:    "douyin", // Consistent issuer name
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)
	if err != nil {
		log.Errorf("生成访问令牌失败：%s", err)
		return "", "", err
	}

	// Refresh Token Claims
	refreshClaims := RefreshClaims{
		UserId:         userId,
		Username:       username, // Include username for convenience if needed upon refresh
		IsRefreshToken: true,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshTokenExpireTime.Unix(),
			Issuer:    "douyin", // Consistent issuer name
		},
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(jwtSecret)
	if err != nil {
		log.Errorf("生成刷新令牌失败：%s", err)
		return "", "", err
	}

	// Store access token in Redis for revocation check (optional, but present in original code)
	// The key "jwt:%d" might be too generic if other JWTs are stored for the user.
	// Consider a more specific key like "accesstoken_active:%d"
	redisKey := fmt.Sprintf("jwt_access_token:%d", userId) // More specific key
	if err := cache.RedisClient.Set(context.Background(), redisKey, accessToken, consts.AccessTokenExpireDuration).Err(); err != nil {
		log.Errorf("存储访问令牌到 Redis 失败：%s", err)
		// Depending on requirements, this might or might not be a fatal error for token generation
	}

	return accessToken, refreshToken, nil
}

// ParseAccessToken validates and parses an access token string.
// It also checks against Redis if the token is still considered active.
func ParseAccessToken(tokenString string) (*Claims, error) {
	if len(jwtSecret) == 0 {
		log.Error("JWT secret is not initialized during ParseAccessToken. Call jwt.Init() first.")
		return nil, errors.New("JWT secret not initialized")
	}
	tokenClaims, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			// Check if token is still active in Redis (if this strategy is kept)
			redisKey := fmt.Sprintf("jwt_access_token:%d", claims.UserId)
			storedToken, redisErr := cache.RedisClient.Get(context.Background(), redisKey).Result()
			if redisErr != nil {
				log.Errorf("从 Redis 获取访问令牌失败（可能已过期或吊销）：%s", redisErr)
				return nil, errors.New("令牌无效或已过期") // More generic message to client
			}
			if storedToken != tokenString {
				log.Errorf("访问令牌不匹配 Redis 中存储的令牌，可能已被吊销或更新")
				return nil, errors.New("令牌已更新或吊销")
			}
			return claims, nil
		}
	}
	// Consolidate error reporting
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("令牌格式错误")
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return nil, errors.New("令牌已过期或未生效")
			} else {
				return nil, errors.New("无法处理的令牌")
			}
		}
		return nil, fmt.Errorf("解析令牌失败: %w", err)
	}
	return nil, errors.New("未知令牌解析错误") // Should not happen if err is nil and claims not valid
}

// ParseRefreshToken validates and parses a refresh token string.
func ParseRefreshToken(tokenString string) (*RefreshClaims, error) {
	if len(jwtSecret) == 0 {
		log.Error("JWT secret is not initialized during ParseRefreshToken. Call jwt.Init() first.")
		return nil, errors.New("JWT secret not initialized")
	}
	tokenClaims, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*RefreshClaims); ok && tokenClaims.Valid {
			if !claims.IsRefreshToken {
				return nil, errors.New("提供的令牌不是有效的刷新令牌")
			}
			return claims, nil
		}
	}
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("刷新令牌格式错误")
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return nil, errors.New("刷新令牌已过期或未生效")
			} else {
				return nil, errors.New("无法处理的刷新令牌")
			}
		}
		return nil, fmt.Errorf("解析刷新令牌失败: %w", err)
	}
	return nil, errors.New("未知刷新令牌解析错误")
}


// RevokeToken 吊销用户 Token，将 Redis 中存储的令牌删除
// This now specifically targets the access token stored in Redis.
func RevokeToken(userId uint) error {
	redisKey := fmt.Sprintf("jwt_access_token:%d", userId) // Match the key used in GenerateToken
	err := cache.RedisClient.Del(context.Background(), redisKey).Err()
	if err != nil {
		log.Errorf("吊销令牌失败：%s", err.Error())
		return err
	}
	log.Infof("成功吊销令牌，用户ID=%d", userId)
	return nil
}

// EmailClaims 定义邮箱验证 Token 的负载信息
type EmailClaims struct {
	UserId        uint   `json:"user_id"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	OperationType uint   `json:"operation_type"`
	jwt.StandardClaims
}

// GenerateEmailToken 签发用于邮箱验证的 Token，15 分钟后过期
func GenerateEmailToken(userId, Operation uint, email, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(15 * time.Minute)

	claims := EmailClaims{
		UserId:        userId,
		Email:         email,
		Password:      password,
		OperationType: Operation,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "douyin商城邮箱验证",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	if err != nil {
		log.Errorf("生成邮箱验证Token失败：%s", err.Error())
		return "", err
	}
	log.Infof("生成邮箱验证Token成功，用户ID=%d, 邮箱=%s", userId, email)
	return token, nil
}

// ParseEmailToken 验证并解析邮箱验证 Token，返回 EmailClaims 信息
func ParseEmailToken(token string) (*EmailClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &EmailClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*EmailClaims); ok && tokenClaims.Valid {
			log.Infof("解析邮箱验证Token成功，用户ID=%d", claims.UserId)
			return claims, nil
		}
	}
	log.Errorf("解析邮箱验证Token失败：%s", err.Error())
	return nil, err
}

// ValidateJWT 验证请求中的 JWT 令牌，并提取用户 ID
func ValidateJWT(c *gin.Context) (uint, error) {
	// 从请求头中获取 Authorization 令牌
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		log.Errorf("缺少 Authorization 头部")
		return 0, errors.New("缺少 Authorization 头部")
	}

	// 令牌的格式通常是 "Bearer <token>"
	const BearerPrefix = "Bearer "
	tokenString := strings.TrimPrefix(authHeader, BearerPrefix)
	if tokenString == authHeader {
		log.Errorf("Authorization 格式无效")
		return 0, errors.New("Authorization 格式无效")
	}

	// 解析令牌并获取用户 ID
	claims, err := ParseToken(tokenString)
	if err != nil {
		log.Errorf("无效令牌: %s", err.Error())
		return 0, errors.New("无效令牌")
	}

	// 获取当前时间
	now := time.Now().Unix()

	// 检查令牌是否过期
	if claims.ExpiresAt < now {
		log.Errorf("令牌已过期")
		return 0, errors.New("令牌已过期")
	}

	// 返回用户 ID
	return claims.UserId, nil
}

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

// 你的签名密钥
var jwtSecret = []byte("DouyinSecret")

// Claims 定义 JWT 负载
type Claims struct {
	UserId   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateToken 创建 access_token, refresh_token 并存到 Redis
func GenerateToken(userId uint, username string) (string, string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(consts.AccessTokenExpireDuration)
	rtExpireTime := nowTime.Add(consts.RefreshTokenExpireDuration)

	// 构造 claims
	claims := Claims{
		UserId:   userId,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "douyin商城",
		},
	}

	// 生成访问令牌
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)
	if err != nil {
		log.Errorf("生成访问令牌失败：%s", err)
		return "", "", err
	}

	// 生成刷新令牌
	refreshClaims := jwt.StandardClaims{
		ExpiresAt: rtExpireTime.Unix(),
		Issuer:    "douyin商城",
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(jwtSecret)
	if err != nil {
		log.Errorf("生成刷新令牌失败：%s", err)
		return "", "", err
	}

	// 写 Redis 做吊销判断
	redisKey := fmt.Sprintf("jwt:%d", userId)
	if err := cache.RedisClient.Set(context.Background(), redisKey, accessToken, consts.AccessTokenExpireDuration).Err(); err != nil {
		log.Errorf("存储访问令牌到 Redis 失败：%s", err)
		// 看你需求决定要不要 return err
	}

	return accessToken, refreshToken, nil
}

// ParseToken 校验 token 签名 + 是否在 Redis 中仍然有效
func ParseToken(tokenString string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		log.Infof("ParseToken 使用的 secret = %s", string(jwtSecret))
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			// 校验 Redis 是否还存着
			redisKey := fmt.Sprintf("jwt:%d", claims.UserId)
			storedToken, err := cache.RedisClient.Get(context.Background(), redisKey).Result()
			if err != nil {
				log.Errorf("从 Redis 获取令牌失败：%s", err)
				return nil, errors.New("身份验证失败")
			}
			// 如果不相等，说明已经被吊销
			if storedToken != tokenString {
				log.Errorf("令牌不匹配，可能已被吊销")
				return nil, errors.New("令牌无效")
			}
			// 解析成功
			return claims, nil
		}
	}
	if err == nil {
		err = errors.New("解析 Token 失败")
	}
	return nil, err
}

// 其余 RevokeToken、ParseEmailToken、GenerateEmailToken 等和你原始写法差别不大

// RevokeToken 吊销用户 Token，将 Redis 中存储的令牌删除
func RevokeToken(userId uint) error {
	redisKey := fmt.Sprintf("jwt:%d", userId)
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

package v1

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"ginDYMall/pkg/utils/response" // Assuming this path
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// HealthController 健康检查
type HealthController struct {
	DB    *gorm.DB
	Redis *redis.Client
}

// RegisterRoutes 挂载 /healthz 路由
func (hc *HealthController) RegisterRoutes(r *gin.Engine) {
	r.GET("/healthz", hc.Healthz)
}

// Healthz 检测 MySQL 与 Redis 状态
func (hc *HealthController) Healthz(c *gin.Context) {
	ctx := c.Request.Context() // Use the request's context

	// Check MySQL
	if hc.DB != nil {
		sqlDB, err := hc.DB.DB()
		if err != nil {
			response.Fail(c, http.StatusInternalServerError, "无法获取数据库连接")
			return
		}
		if err := sqlDB.PingContext(ctx); err != nil {
			response.Fail(c, http.StatusServiceUnavailable, "MySQL 未连接")
			return
		}
	} else {
		response.Fail(c, http.StatusInternalServerError, "数据库客户端未初始化")
		return
	}

	// Check Redis
	if hc.Redis != nil {
		if err := hc.Redis.Ping(ctx).Err(); err != nil {
			response.Fail(c, http.StatusServiceUnavailable, "Redis 未连接")
			return
		}
	} else {
		response.Fail(c, http.StatusInternalServerError, "Redis客户端未初始化")
		return
	}
	// 全部正常
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

package cache

import (
	"context"
	"douyin/config"
	"fmt"
	"github.com/redis/go-redis/v9"
	logging "github.com/sirupsen/logrus"
)

// RedisClient 为 Redis 客户端单例
var RedisClient *redis.Client

// RedisContext 定义 Redis 操作的上下文
var RedisContext = context.Background()

// InitCache 初始化 Redis 连接，返回错误信息
func InitCache() error {
	// 从全局配置中获取 Redis 配置
	rConfig := config.GlobalConfig.Redis
	client := redis.NewClient(&redis.Options{
		// 组装 Redis 地址：host:port
		Addr:     fmt.Sprintf("%s:%s", rConfig.RedisHost, rConfig.RedisPort),
		Username: rConfig.RedisUsername,
		Password: rConfig.RedisPassword,
		DB:       rConfig.RedisDbName,
	})
	// Ping 测试连接是否成功
	_, err := client.Ping(RedisContext).Result()
	if err != nil {
		// 记录连接错误并返回
		logging.Errorf("Redis 连接错误：%v", err)
		return fmt.Errorf("Redis 连接失败: %v", err)
	}

	// 成功初始化 Redis 客户端
	RedisClient = client
	logging.Info("Redis 初始化成功")
	return nil
}

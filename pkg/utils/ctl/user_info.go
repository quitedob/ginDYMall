// 文件：pkg/utils/ctl/user_info.go
// 作用：提供上下文中用户信息的存取方法，统一定义用户信息结构体（使用 UserId 字段）
// 说明：所有接口中需使用该模块存取用户信息，便于统一管理和鉴权

package ctl

import (
	"context"
	"errors"

	"douyin/pkg/utils/log" // 日志工具包，用于中文日志输出
)

// 定义上下文中存储用户信息的 key 类型
type key int

var userKey key // 用户信息在上下文中的 key

// UserInfo 定义用户信息结构体，统一使用 UserId 字段
type UserInfo struct {
	UserId uint `json:"user_id"` // 用户唯一标识
}

// GetUserInfo 从上下文中获取用户信息，如果获取失败则返回错误
func GetUserInfo(ctx context.Context) (*UserInfo, error) {
	user, ok := FromContext(ctx)
	if !ok {
		errMsg := "获取用户信息失败：上下文中未找到用户数据"
		log.Errorf(errMsg)
		return nil, errors.New(errMsg)
	}
	log.Infof("成功获取用户信息：用户ID=%d", user.UserId)
	return user, nil
}

// NewContext 将用户信息存入上下文中，并返回新的上下文
func NewContext(ctx context.Context, u *UserInfo) context.Context {
	log.Infof("将用户信息存入上下文：用户ID=%d", u.UserId)
	return context.WithValue(ctx, userKey, u)
}

// FromContext 尝试从上下文中提取用户信息
func FromContext(ctx context.Context) (*UserInfo, bool) {
	u, ok := ctx.Value(userKey).(*UserInfo)
	return u, ok
}

// InitUserInfo 初始化用户信息（后续可引入缓存逻辑）
// 这里暂时仅输出初始化提示信息，未来可接入缓存逻辑以提升性能
func InitUserInfo(ctx context.Context) {
	log.Infof("正在初始化用户信息...（未来将加入缓存功能）")
}

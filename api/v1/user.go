// 文件：api/v1/user.go
// 作用：定义用户相关的 HTTP API 接口，包括创建用户、登录、注销、修改密码、修改昵称、更新用户信息和展示用户身份信息
// 使用 Gin 框架处理 HTTP 请求，调用业务逻辑层接口处理请求数据，并返回 JSON 格式响应

package v1

import (
	"douyin/pkg/utils/ctl" // 用于获取上下文中的用户信息
	"douyin/pkg/utils/log" // 日志工具包
	"douyin/service"       // 业务逻辑层
	"douyin/types"         // 数据传输对象包
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserRegisterHandler 用户注册接口（创建用户）
func UserRegisterHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.UserRegisterReq // 注册请求数据结构体
		if err := ctx.ShouldBind(&req); err != nil {
			log.LogrusObj.Infoln("绑定注册请求数据失败：", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// 调用业务层注册函数
		resp, err := service.GetUserSrv().UserRegister(ctx.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Infoln("用户注册失败：", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		log.LogrusObj.Info("注册成功，用户创建成功")
		ctx.JSON(http.StatusOK, gin.H{"data": resp})
	}
}

// UserLoginHandler 用户登录接口（输入用户ID和密码）
func UserLoginHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.UserLoginReq // 登录请求数据结构体
		if err := ctx.ShouldBind(&req); err != nil {
			log.LogrusObj.Infoln("绑定登录请求数据失败：", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// 调用业务层登录函数
		resp, err := service.GetUserSrv().UserLogin(ctx.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Infoln("用户登录失败：", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		log.LogrusObj.Info("登录成功，生成令牌")
		ctx.JSON(http.StatusOK, gin.H{"data": resp})
	}
}

// UserLogoutHandler 用户注销接口（删除 Redis 中存储的 token）
func UserLogoutHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从上下文中获取用户信息（依赖鉴权中间件提前解析 Token）
		_, err := ctl.GetUserInfo(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
			return
		}
		if err := service.GetUserSrv().UserLogout(ctx.Request.Context()); err != nil {
			log.LogrusObj.Infoln("用户注销失败：", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "注销失败"})
			return
		}
		log.LogrusObj.Info("注销成功")
		ctx.JSON(http.StatusOK, gin.H{"message": "注销成功"})
	}
}

// UserChangePasswordHandler 修改密码接口
func UserChangePasswordHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.UserChangePasswordReq // 修改密码请求数据结构体
		if err := ctx.ShouldBind(&req); err != nil {
			log.LogrusObj.Infoln("绑定修改密码请求数据失败：", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// 调用业务层修改密码函数
		resp, err := service.GetUserSrv().UserChangePassword(ctx.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Infoln("修改密码失败：", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		log.LogrusObj.Info("密码修改成功")
		ctx.JSON(http.StatusOK, gin.H{"data": resp})
	}
}

// UserChangeNicknameHandler 修改昵称接口
func UserChangeNicknameHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.UserChangeNicknameReq // 修改昵称请求数据结构体
		if err := ctx.ShouldBind(&req); err != nil {
			log.LogrusObj.Infoln("绑定修改昵称请求数据失败：", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// 调用业务层修改昵称函数
		resp, err := service.GetUserSrv().UserChangeNickname(ctx.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Infoln("修改昵称失败：", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		log.LogrusObj.Info("昵称修改成功")
		ctx.JSON(http.StatusOK, gin.H{"data": resp})
	}
}

// UserUpdateHandler 更新用户信息接口（仅更新 user_name, email, updated_at）
func UserUpdateHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.UserUpdateReq // 更新用户请求数据结构体
		if err := ctx.ShouldBind(&req); err != nil {
			log.LogrusObj.Infoln("绑定更新用户请求数据失败：", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// 调用业务层更新用户信息函数
		resp, err := service.GetUserSrv().UserUpdate(ctx.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Infoln("更新用户信息失败：", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		log.LogrusObj.Info("用户更新成功")
		ctx.JSON(http.StatusOK, gin.H{"data": resp})
	}
}

// UserInfoShowHandler 用户身份信息展示接口
// 根据用户ID和 JWT 令牌验证后返回：[用户ID, 用户名, 用户邮箱, 解密后的余额, 创建时间, 更新时间]
func UserInfoShowHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.UserInfoShowReq // 用户信息查询请求数据结构体
		if err := ctx.ShouldBindQuery(&req); err != nil {
			log.LogrusObj.Infoln("绑定用户信息查询请求数据失败：", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// 调用业务层获取用户身份信息函数
		resp, err := service.GetUserSrv().UserInfoShow(ctx.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Infoln("获取用户身份信息失败：", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		log.LogrusObj.Info("用户身份信息获取成功")
		ctx.JSON(http.StatusOK, gin.H{"data": resp})
	}
}

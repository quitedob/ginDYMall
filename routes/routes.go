package routes

import (
	"time" // For rate limiting window

	"github.com/gin-gonic/gin" // 导入 Gin 框架
	"gorm.io/gorm"             // 导入 GORM 数据库驱动

	"douyin/api/v1"     // 导入 API V1 版本的包，其中包含 CartController 的定义
	"douyin/cache"      // For cache.Rdb (Redis client)
	"douyin/middleware" // 导入中间件包（如身份验证中间件）
)

// NewRouter 根据传入的数据库实例 db 初始化并返回一个 Gin 引擎实例
func NewRouter(engine *gin.Engine, db *gorm.DB) {
	// 定义 API V1 版本分组
	apiV1 := engine.Group("/api/v1")
	{
		// 注册无需登录的接口

		// Apply to register: 2 requests per 10 minutes per IP (more strict for registration)
		// Assumes cache.Rdb is the initialized *redis.Client
		apiV1.POST("/user/register",
			middleware.RateLimitMiddleware(cache.Rdb, "register_ip", 2, 10*time.Minute),
			v1.UserRegisterHandler(), // 用户注册接口
		)
		// Apply to login: 5 requests per minute per user/IP
		apiV1.POST("/user/login",
			middleware.RateLimitMiddleware(cache.Rdb, "login", 5, 1*time.Minute),
			v1.UserLoginHandler(), // 用户登录接口
		)

		apiV1.POST("/product/get", v1.GetProduct)   // 获取单个商品接口
		apiV1.GET("/product/list", v1.ListProducts) // 获取商品列表接口

		// 定义需要登录验证的接口分组
		authGroup := apiV1.Group("")
		// 应用身份验证中间件
		authGroup.Use(middleware.AuthMiddleware())
		{
			// 用户相关接口
			authGroup.POST("user/change_password", v1.UserChangePasswordHandler()) // 修改密码接口
			authGroup.POST("user/change_nickname", v1.UserChangeNicknameHandler()) // 修改昵称接口
			authGroup.POST("user/update", v1.UserUpdateHandler())                  // 更新用户信息接口
			authGroup.GET("user/show_info", v1.UserInfoShowHandler())              // 显示用户信息接口
			authGroup.POST("user/logout", v1.UserLogoutHandler())                  // 用户登出接口

			// 订单相关接口
			authGroup.POST("order/create", v1.OrderCreateHandler()) // 创建订单接口
			authGroup.POST("order/update", v1.OrderUpdateHandler()) // 更新订单接口

			// 商品相关接口
			authGroup.POST("product/create", v1.CreateProduct)          // 创建商品接口
			authGroup.POST("product/update", v1.UpdateProduct)          // 更新商品接口
			authGroup.POST("product/delete", v1.DeleteProduct)          // 删除商品接口
			authGroup.POST("checkout/order", v1.CheckoutOrderHandler()) // 结算订单接口

			// 购物车相关接口
			// 创建 CartController 的实例，传入数据库实例 db
			cartController := v1.NewCartController(db)
			authGroup.POST("cart/create", cartController.CreateCart) // 创建购物车接口
			authGroup.POST("cart/get", cartController.GetCart)       // 获取购物车信息接口
			authGroup.POST("cart/empty", cartController.EmptyCart)   // 清空购物车接口
			authGroup.POST("cart/add", cartController.AddItem)       // 添加或更新购物车商品接口
		}
	}
}

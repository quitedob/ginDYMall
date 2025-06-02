// /api/v1/order.go
package v1

import (
	"douyin/service"
	"douyin/types"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
)

// OrderControllerType 封装订单操作
type OrderControllerType struct {
	service *service.OrderService
}

// OrderController 是全局订单控制器实例
var OrderController *OrderControllerType

// SetDB 用于初始化订单控制器，并将 db 传入
func SetDB(db *gorm.DB) {
	oc, err := NewOrderController(db)
	if err != nil {
		log.Println("OrderController 初始化失败:", err)
		return
	}
	OrderController = oc
	log.Println("OrderController 初始化成功")
}

// NewOrderController 创建新的订单控制器实例
func NewOrderController(db *gorm.DB) (*OrderControllerType, error) {
	// 创建订单服务，并处理可能的错误
	svc, err := service.NewOrderService(db)
	if err != nil {
		log.Println("OrderService 初始化失败:", err)
		return nil, err
	}
	log.Println("OrderService 初始化成功")
	return &OrderControllerType{
		service: svc,
	}, nil
}

// OrderCreateHandler 创建订单的处理函数
func OrderCreateHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 检查全局 OrderController 是否已初始化
		if OrderController == nil {
			log.Println("OrderController 未正确初始化")
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "OrderController 未正确初始化"})
			return
		}

		// 检查 OrderService 是否已初始化
		if OrderController.service == nil {
			log.Println("OrderService 未正确初始化")
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "OrderService 未正确初始化"})
			return
		}

		// 调用 CreateOrder 方法
		OrderController.CreateOrder(ctx)
	}
}

// OrderUpdateHandler 更新订单的处理函数
func OrderUpdateHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 检查全局 OrderController 是否已初始化
		if OrderController == nil {
			log.Println("OrderController 未正确初始化")
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "OrderController 未正确初始化"})
			return
		}

		// 检查 OrderService 是否已初始化
		if OrderController.service == nil {
			log.Println("OrderService 未正确初始化")
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "OrderService 未正确初始化"})
			return
		}

		// 调用 UpdateOrder 方法
		OrderController.UpdateOrder(ctx)
	}
}

// CreateOrder 调用服务层创建订单
func (c *OrderControllerType) CreateOrder(ctx *gin.Context) {
	// 检查服务是否已正确初始化
	if c.service == nil {
		log.Println("OrderService 未正确初始化")
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "OrderService 未正确初始化"})
		return
	}

	// 获取 userID，检查是否存在
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "用户未认证"})
		return
	}

	// 确保 userID 是 uint 类型
	userIDUint, ok := userID.(uint)
	if !ok {
		log.Printf("类型断言失败，userID = %v", userID)
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "无效的用户ID"})
		return
	}

	// 解析请求体
	var req types.CreateOrderReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "无效输入"})
		return
	}

	// 调用服务层创建订单
	orderID, err := c.service.CreateOrder(userIDUint, &req)
	if err != nil {
		log.Printf("创建订单时出错: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "创建订单时出错"})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{"order_id": orderID})
}

// UpdateOrder 调用服务层更新订单
func (c *OrderControllerType) UpdateOrder(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	var req types.UpdateOrderReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "无效输入"})
		return
	}

	// 将 userID 赋值给请求
	req.UserID = userID.(uint)

	if err := c.service.UpdateOrder(req.UserID, &req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "更新订单时出错"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "订单更新成功"})
}

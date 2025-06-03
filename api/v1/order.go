// /api/v1/order.go
package v1

import (
	"douyin/pkg/utils/response"
	"douyin/service"
	"douyin/types"
	"errors"
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
		if OrderController == nil || OrderController.service == nil {
			log.Println("OrderController 或 OrderService 未正确初始化")
			_ = ctx.Error(errors.New("服务内部错误：订单服务未就绪"))
			return
		}
		OrderController.CreateOrder(ctx)
	}
}

// OrderUpdateHandler 更新订单的处理函数
func OrderUpdateHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if OrderController == nil || OrderController.service == nil {
			log.Println("OrderController 或 OrderService 未正确初始化")
			_ = ctx.Error(errors.New("服务内部错误：订单服务未就绪"))
			return
		}
		OrderController.UpdateOrder(ctx)
	}
}

// CreateOrder 调用服务层创建订单
func (c *OrderControllerType) CreateOrder(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id") // Key "user_id" from AuthMiddleware
	if !exists {
		_ = ctx.Error(errors.New("用户未授权或user_id未在context中设置"))
		return
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		_ = ctx.Error(errors.New("user_id在context中的类型错误"))
		return
	}

	var req types.CreateOrderReq // Use the new types.CreateOrderReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.Fail(1001, "参数非法："+err.Error()))
		return
	}

	// 调用服务层创建订单, using the new service signature
	orderID, err := c.service.CreateOrder(ctx.Request.Context(), userID, req.AddressID, req.Items)
	if err != nil {
		// log.Printf("创建订单时出错: %v", err) // Service/DAO layer should log specifics
		_ = ctx.Error(err) // Pass to global error handler
		return
	}

	ctx.JSON(http.StatusOK, response.Success(gin.H{"order_id": orderID}))
}

// UpdateOrder 调用服务层更新订单
func (c *OrderControllerType) UpdateOrder(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		_ = ctx.Error(errors.New("用户未授权或user_id未在context中设置"))
		return
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		_ = ctx.Error(errors.New("user_id在context中的类型错误"))
		return
	}

	var req types.UpdateOrderReq // This request struct still has UserID, but we ignore it.
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.Fail(1001, "参数非法："+err.Error()))
		return
	}

	// UserID from JWT (userID variable) is authoritative, not req.UserID.
	// The service layer's UpdateOrder should accept userID from context.
	if err := c.service.UpdateOrder(ctx.Request.Context(), userID, &req); err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Success("订单更新成功"))
}

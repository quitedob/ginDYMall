package v1

import (
	// "douyin/pkg/utils/jwt" // JWT validation is now part of AuthMiddleware
	"douyin/pkg/utils/response"
	"douyin/service"
	"douyin/types"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CartController 购物车控制器
type CartController struct {
	service *service.CartService
}

// NewCartController 创建新的 CartController 实例
func NewCartController(db *gorm.DB) *CartController {
	return &CartController{
		service: service.NewCartService(db),
	}
}

// CreateCart 创建购物车接口
func (c *CartController) CreateCart(ctx *gin.Context) {
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

	// var req types.CreateCartReq // No fields in CreateCartReq anymore
	// if err := ctx.ShouldBindJSON(&req); err != nil { ... } // Not needed if req is empty

	// 创建购物车
	if err := c.service.CreateCart(ctx.Request.Context(), userID); err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Success("购物车创建成功"))
}

// GetCart 获取购物车信息接口
func (c *CartController) GetCart(ctx *gin.Context) {
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

	// var req types.GetCartReq // No fields
	// if err := ctx.ShouldBindJSON(&req); err != nil { ... }

	// 获取购物车信息
	cart, err := c.service.GetCart(ctx.Request.Context(), userID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Success(cart))
}

// EmptyCart 清空购物车接口
func (c *CartController) EmptyCart(ctx *gin.Context) {
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

	// var req types.EmptyCartReq // No fields
	// if err := ctx.ShouldBindJSON(&req); err != nil { ... }

	// 清空购物车
	if err := c.service.EmptyCart(ctx.Request.Context(), userID); err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Success("购物车已清空"))
}

// AddItem 往购物车中添加(或更新)商品接口
func (c *CartController) AddItem(ctx *gin.Context) {
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

	var req types.AddItemReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.Fail(1001, "参数非法："+err.Error()))
		return
	}

	// 调用service添加或更新
	if err := c.service.AddItem(ctx.Request.Context(), userID, req.ProductID, int32(req.Quantity)); err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Success("添加(或更新)商品成功"))
}

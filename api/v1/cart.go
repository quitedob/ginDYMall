package v1

import (
	"net/http"

	"douyin/pkg/utils/jwt"
	"douyin/service"
	"douyin/types"

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
	// 验证JWT token
	_, err := jwt.ValidateJWT(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "未授权"})
		return
	}

	var req types.CreateCartReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "无效输入"})
		return
	}

	// 创建购物车
	if err := c.service.CreateCart(req.UserID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "创建购物车时出错"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "购物车创建成功"})
}

// GetCart 获取购物车信息接口
func (c *CartController) GetCart(ctx *gin.Context) {
	// 验证JWT token
	_, err := jwt.ValidateJWT(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "未授权"})
		return
	}

	var req types.GetCartReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "无效输入"})
		return
	}

	// 获取购物车信息
	cart, err := c.service.GetCart(req.UserID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "购物车未找到"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"cart": cart})
}

// EmptyCart 清空购物车接口
func (c *CartController) EmptyCart(ctx *gin.Context) {
	// 验证JWT token
	_, err := jwt.ValidateJWT(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "未授权"})
		return
	}

	var req types.EmptyCartReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "无效输入"})
		return
	}

	// 清空购物车
	if err := c.service.EmptyCart(req.UserID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "清空购物车时出错"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "购物车已清空"})
}

// AddItem 往购物车中添加(或更新)商品接口
func (c *CartController) AddItem(ctx *gin.Context) {
	// 验证JWT token
	_, err := jwt.ValidateJWT(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "未授权"})
		return
	}

	var req types.AddItemReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "无效输入"})
		return
	}

	// 调用service添加或更新
	if err := c.service.AddItem(req.UserID, req.ProductID, req.Quantity); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "添加商品失败", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "添加(或更新)商品成功"})
}

// API/V1/CHECKOUT.GO
package v1

import (
	"douyin/pkg/utils/jwt"
	"douyin/pkg/utils/response"
	"douyin/service"
	"douyin/types"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// 声明全局结算控制器变量，供路由包装函数调用
var checkoutController *CheckoutController

// CheckoutController 结算控制器
type CheckoutController struct {
	service *service.CheckoutService
}

// NewCheckoutController 创建新的 CheckoutController 实例
func NewCheckoutController(db *gorm.DB) *CheckoutController {
	return &CheckoutController{
		service: service.NewCheckoutService(db),
	}
}

// CheckoutOrder 订单结算接口
func (c *CheckoutController) CheckoutOrder(ctx *gin.Context) {
	// 验证JWT token
	userID, err := jwt.ValidateJWT(ctx)
	if err != nil {
		_ = ctx.Error(errors.New("未授权"))
		return
	}

	// 修改处：使用 types.CreateOrderReq 代替 types.CheckoutReq
	var req types.CreateOrderReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.Fail(1001, "参数非法："+err.Error()))
		return
	}

	// 进行结算
	transactionID, err := c.service.CheckoutOrder(userID, &req)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Success(gin.H{"transaction_id": transactionID}))
}

// CheckoutOrderHandler
func CheckoutOrderHandler() gin.HandlerFunc {
	return checkoutController.CheckoutOrder
}

// SetCheckoutController
func SetCheckoutController(db *gorm.DB) {
	checkoutController = NewCheckoutController(db)
}

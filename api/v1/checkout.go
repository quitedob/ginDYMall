// API/V1/CHECKOUT.GO
package v1

import (
	"douyin/pkg/utils/jwt"
	"douyin/service"
	"douyin/types"
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
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "未授权"})
		return
	}

	// 修改处：使用 types.CreateOrderReq 代替 types.CheckoutReq
	var req types.CreateOrderReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "无效输入"})
		return
	}

	// 进行结算
	transactionID, err := c.service.CheckoutOrder(userID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "结算时出错"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"transaction_id": transactionID})
}

// CheckoutOrderHandler
func CheckoutOrderHandler() gin.HandlerFunc {
	return checkoutController.CheckoutOrder
}

// SetCheckoutController
func SetCheckoutController(db *gorm.DB) {
	checkoutController = NewCheckoutController(db)
}

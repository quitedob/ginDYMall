package service

import (
	"douyin/repository/db/dao"
	"douyin/types"
	"gorm.io/gorm"
)

// CheckoutService 结算服务
type CheckoutService struct {
	dao *dao.CheckoutDao
}

// NewCheckoutService 创建新的 CheckoutService 实例
func NewCheckoutService(db *gorm.DB) *CheckoutService {
	return &CheckoutService{
		dao: dao.NewCheckoutDao(db),
	}
}

// CheckoutOrder 订单结算接口
// 修改参数类型为 *types.CreateOrderReq 以匹配 DAO 层要求
func (s *CheckoutService) CheckoutOrder(userID uint, req *types.CreateOrderReq) (string, error) {
	return s.dao.CheckoutOrder(userID, req)
}

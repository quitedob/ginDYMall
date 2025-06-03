package service

import (
	"context"
	"douyin/repository/db/dao"
	"douyin/repository/db/model"
	"gorm.io/gorm"
)

// CartService 购物车服务，封装购物车业务逻辑
type CartService struct {
	dao *dao.CartDao
}

// NewCartService 创建 CartService 实例
func NewCartService(db *gorm.DB) *CartService {
	return &CartService{
		dao: dao.NewCartDao(db),
	}
}

// CreateCart 创建空购物车
func (s *CartService) CreateCart(ctx context.Context, userID uint) error {
	// UserID is now passed directly
	return s.dao.CreateCart(ctx, userID)
}

// GetCart 获取用户购物车信息
func (s *CartService) GetCart(ctx context.Context, userID uint) ([]model.CartItem, error) {
	// UserID is now passed directly
	return s.dao.GetCart(ctx, userID)
}

// EmptyCart 清空购物车
func (s *CartService) EmptyCart(ctx context.Context, userID uint) error {
	// UserID is now passed directly
	return s.dao.EmptyCart(ctx, userID)
}

// AddItem 往购物车中添加(或更新)商品
// 由于本例中只做简单转发给 dao，因此可以在此处加一些
// 例如商品合法性校验、库存校验等高级逻辑
func (s *CartService) AddItem(ctx context.Context, userID, productID uint, quantity int32) error {
	// 此处可添加业务逻辑，比如先查库存是否足够等...
	// The DAO layer now handles stock checking and deduction within a transaction.
	return s.dao.AddItem(ctx, userID, productID, quantity)
}

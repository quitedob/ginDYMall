package service

import (
	"douyin/pkg/utils/log"
	"douyin/repository/db/dao"
	"douyin/types"
	"gorm.io/gorm"
)

// OrderService 订单服务
type OrderService struct {
	dao *dao.OrderDao
}

// NewOrderService 创建新的 OrderService 实例
func NewOrderService(db *gorm.DB) (*OrderService, error) {
	// 获取底层的 *sql.DB 实例，并检查数据库连接是否成功
	sqlDB, err := db.DB()
	if err != nil {
		log.Errorf("获取数据库连接失败: %v", err)
		return nil, err
	}
	if err := sqlDB.Ping(); err != nil {
		log.Errorf("数据库连接失败: %v", err)
		return nil, err
	}
	log.Infof("订单服务使用的数据库连接成功")

	return &OrderService{
		dao: dao.NewOrderDao(db),
	}, nil
}

func (s *OrderService) CreateOrder(userID uint, req *types.CreateOrderReq) (string, error) {
	// 创建订单
	return s.dao.CreateOrder(userID, req)
}

func (s *OrderService) UpdateOrder(userID uint, req *types.UpdateOrderReq) error {
	// 更新订单
	return s.dao.UpdateOrder(userID, req)
}

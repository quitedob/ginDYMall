package dao

import (
	"douyin/repository/db/model"
	"douyin/types"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// OrderDao 定义订单数据访问对象
type OrderDao struct {
	db *gorm.DB
}

// NewOrderDao 根据传入的数据库连接创建新的 OrderDao 实例
func NewOrderDao(db *gorm.DB) *OrderDao {
	return &OrderDao{
		db: db,
	}
}

// CreateOrder 创建订单
func (dao *OrderDao) CreateOrder(userID uint, req *types.CreateOrderReq) (string, error) {
	orderID := uuid.New().String()
	order := model.Order{
		OrderID:       orderID,
		UserID:        userID,
		UserCurrency:  req.UserCurrency,
		Email:         req.Email,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		StreetAddress: req.StreetAddress,
		City:          req.City,
		State:         req.State,
		Country:       req.Country,
		ZipCode:       req.ZipCode,
		CreatedAt:     time.Now(),
		Status:        "pending",
	}

	tx := dao.db.Begin()

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return "", err
	}

	for _, item := range req.OrderItems {
		orderItem := model.OrderItem{
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Cost:      100.0,
		}

		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			return "", err
		}
	}

	tx.Commit()

	return orderID, nil
}

// UpdateOrder 更新订单
func (dao *OrderDao) UpdateOrder(userID uint, req *types.UpdateOrderReq) error {
	var order model.Order
	if err := dao.db.Where("user_id = ? AND order_id = ?", userID, req.OrderID).First(&order).Error; err != nil {
		return err
	}

	order.StreetAddress = req.StreetAddress
	order.City = req.City
	order.State = req.State
	order.Country = req.Country
	order.ZipCode = req.ZipCode
	order.Status = req.Status

	return dao.db.Save(&order).Error
}

package dao

import (
	"context"
	"douyin/repository/db/model"
	"douyin/types"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
func (dao *OrderDao) CreateOrder(ctx context.Context, userID uint, req *types.CreateOrderReq) (string, error) {
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
		Status:        "pending", // Initial status
	}

	tx := dao.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return "", tx.Error
	}
	defer tx.Rollback() // Rollback if not committed

	if err := tx.Create(&order).Error; err != nil {
		return "", err
	}

	var totalAmount float64 = 0 // Calculate total amount for payment

	for _, item := range req.OrderItems {
		var product model.Product
		// Lock product row for update
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", item.ProductID).First(&product).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return "", errors.New("商品不存在")
			}
			return "", err
		}

		// Check stock
		if product.Stock < item.Quantity {
			return "", errors.New("库存不足: " + product.Name)
		}

		// Update stock using optimistic locking
		result := tx.Model(&model.Product{}).Where("id = ? AND version = ?", product.ID, product.Version).Updates(map[string]interface{}{
			"stock":   gorm.Expr("stock - ?", item.Quantity),
			"version": gorm.Expr("version + 1"),
		})

		if result.Error != nil {
			return "", result.Error
		}
		if result.RowsAffected == 0 {
			return "", errors.New("并发冲突，请重试: " + product.Name)
		}
		
		totalAmount += product.Price * float64(item.Quantity) // Assume item.Cost should be product.Price

		orderItem := model.OrderItem{
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Cost:      product.Price, // Use actual product price for cost
		}

		if err := tx.Create(&orderItem).Error; err != nil {
			return "", err
		}
	}

	// Create a placeholder payment record
	payment := model.Payment{
		OrderID:   orderID,
		Amount:    totalAmount, // This should be calculated based on items
		Status:    "UNPAID",
		CreatedAt: time.Now(),
	}
	if err := tx.Create(&payment).Error; err != nil {
		return "", err
	}

	if err := tx.Commit().Error; err != nil {
		return "", err
	}

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

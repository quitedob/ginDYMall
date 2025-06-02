// model/order_items.go
package model

import (
	"time"
)

// OrderItem 订单项模型
// OrderItem 订单项模型
type OrderItem struct {
	ID        uint      `gorm:"primaryKey"`                 // 订单项ID
	OrderID   string    `gorm:"column:order_id;not null"`   // 订单ID
	ProductID uint      `gorm:"column:product_id;not null"` // 商品ID
	Quantity  int32     `gorm:"column:quantity;not null"`   // 商品数量
	Cost      float64   `gorm:"column:cost;not null"`       // 商品成本（下单时价格）
	CreatedAt time.Time `gorm:"-"`                          // 忽略创建时间字段
}

// 外键约束
func (OrderItem) TableName() string {
	return "order_items"
}

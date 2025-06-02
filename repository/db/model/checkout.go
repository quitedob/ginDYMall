package model

import (
	"time"
)

// Payment 交易记录模型
type Payment struct {
	TransactionID string    `gorm:"primaryKey;column:transaction_id;size:64" json:"transaction_id"` // 交易ID
	OrderID       string    `gorm:"not null;column:order_id;size:64" json:"order_id"`              // 订单ID
	UserID        uint      `gorm:"not null;column:user_id" json:"user_id"`                        // 用户ID
	Amount        float64   `gorm:"not null;column:amount" json:"amount"`                          // 支付金额
	Status        string    `gorm:"column:status;default:'initiated';size:50" json:"status"`       // 支付状态
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`                           // 支付时间
}

// TableName 设置表名
func (Payment) TableName() string {
	return "payments"
}

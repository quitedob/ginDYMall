package model

import (
	"time"
)

// Order 订单模型
type Order struct {
	OrderID       string      `gorm:"primaryKey;column:order_id;size:64" json:"order_id"`                // 订单ID
	UserID        uint        `gorm:"not null;column:user_id" json:"user_id"`                            // 用户ID
	UserCurrency  string      `gorm:"not null;column:user_currency;size:10" json:"user_currency"`        // 用户货币
	Email         string      `gorm:"not null;column:email;size:255" json:"email"`                       // 用户邮箱
	FirstName     string      `gorm:"column:firstname;size:50" json:"first_name"`                        // 名
	LastName      string      `gorm:"column:lastname;size:50" json:"last_name"`                          // 姓
	StreetAddress string      `gorm:"column:street_address;not null;size:255" json:"street_address"`     // 街道地址
	City          string      `gorm:"column:city;not null;size:100" json:"city"`                         // 城市
	State         string      `gorm:"column:state;not null;size:100" json:"state"`                       // 省/州
	Country       string      `gorm:"column:country;not null;size:100" json:"country"`                   // 国家
	ZipCode       string      `gorm:"column:zip_code;not null;size:20" json:"zip_code"`                  // 邮政编码
	CreatedAt     time.Time   `gorm:"column:created_at" json:"created_at"`                               // 订单创建时间
	Status        string      `gorm:"column:status;default:'pending';size:50" json:"status"`             // 订单状态
	OrderItems    []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"order_items"` // 订单项
}

// 外键约束
func (Order) TableName() string {
	return "orders"
}

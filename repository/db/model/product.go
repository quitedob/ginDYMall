// model/product.go
package model

import (
	"time"
)

// Product 商品模型
type Product struct {
	ID          uint      `gorm:"primaryKey"`            // 商品ID
	CreatedAt   time.Time `gorm:"column:created_at"`     // 商品创建时间
	UpdatedAt   time.Time `gorm:"column:updated_at"`     // 商品更新时间
	Name        string    `gorm:"column:name;not null"`  // 商品名称
	Description string    `gorm:"column:description"`    // 商品描述
	Picture     string    `gorm:"column:picture"`        // 商品图片地址
	Price       float64   `gorm:"column:price;not null"` // 商品价格
}

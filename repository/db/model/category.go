// model/category.go
package model

import (
	"time"
)

// Category 商品分类模型
type Category struct {
	ID        uint      `gorm:"primaryKey"`                  // 分类ID
	CreatedAt time.Time `gorm:"column:created_at"`           // 分类创建时间
	UpdatedAt time.Time `gorm:"column:updated_at"`           // 分类更新时间
	Name      string    `gorm:"column:name;unique;not null"` // 分类名称，唯一
}

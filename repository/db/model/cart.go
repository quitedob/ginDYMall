package model

import (
	"fmt"
	"time"
)

// CartItem 购物车项模型
// 指定与数据库表 cart_items 对应。
// 因为 GORM 默认会将结构体 CartItem 对应到表名 cart_items，
// 但这里我们也可以通过实现 TableName() 来显式指定。
type CartItem struct {
	UserID    uint      `gorm:"primaryKey;autoIncrement:false"`     // 用户ID，复合主键的一部分
	ProductID uint      `gorm:"primaryKey;autoIncrement:false"`     // 商品ID，复合主键的一部分
	Quantity  int32     `gorm:"column:quantity;not null;default:1"` // 商品数量，默认为1
	CreatedAt time.Time `gorm:"column:created_at"`                  // 添加时间
	UpdatedAt time.Time `gorm:"column:updated_at"`                  // 更新时间
}

// TableName 指定当前模型所映射的数据库表名为 cart_items
func (CartItem) TableName() string {
	return "cart_items"
}

// PrintInfo 输出购物车项的详细信息（用于调试）
func (c *CartItem) PrintInfo() {
	fmt.Printf("购物车项信息 - 用户ID: %d, 商品ID: %d, 数量: %d, 添加时间: %s, 更新时间: %s\n",
		c.UserID, c.ProductID, c.Quantity,
		c.CreatedAt.Format("2006-01-02 15:04:05"),
		c.UpdatedAt.Format("2006-01-02 15:04:05"))
}

func init() {
	// 初始化模型时输出提示信息（可省略）
	fmt.Println("初始化购物车模型 CartItem 成功")
}

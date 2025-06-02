// model/product_categories.go
package model

// ProductCategory 商品和分类的关联表
type ProductCategory struct {
	ProductID  uint `gorm:"primaryKey;not null"` // 商品ID
	CategoryID uint `gorm:"primaryKey;not null"` // 分类ID
}

// 外键约束
func (ProductCategory) TableName() string {
	return "product_categories"
}

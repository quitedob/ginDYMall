// types/category.go
package types

// 商品类别
type Category struct {
	ID   uint32 `json:"id"`   // 分类ID
	Name string `json:"name"` // 分类名称
}

// 商品分类关联
type ProductCategory struct {
	ProductID  uint32 `json:"product_id"`  // 商品ID
	CategoryID uint32 `json:"category_id"` // 分类ID
}

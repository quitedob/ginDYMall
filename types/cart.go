package types

// CreateCartReq 创建购物车请求参数
type CreateCartReq struct {
	UserID uint `json:"user_id" binding:"required"` // 用户ID
}

// GetCartReq 获取购物车请求参数
type GetCartReq struct {
	UserID uint `json:"user_id" binding:"required"` // 用户ID
}

// EmptyCartReq 清空购物车请求参数
type EmptyCartReq struct {
	UserID uint `json:"user_id" binding:"required"` // 用户ID
}

// AddItemReq 添加商品到购物车请求参数
type AddItemReq struct {
	UserID    uint  `json:"user_id" binding:"required"`    // 用户ID
	ProductID uint  `json:"product_id" binding:"required"` // 商品ID
	Quantity  int32 `json:"quantity" binding:"required"`   // 要添加的数量（可为正数）
}

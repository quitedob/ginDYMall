package types

// CreateCartReq 创建购物车请求参数 (UserID will be from JWT)
type CreateCartReq struct {
	// UserID is removed, will be extracted from JWT claims in handler
}

// GetCartReq 获取购物车请求参数 (UserID will be from JWT)
type GetCartReq struct {
	// UserID is removed, will be extracted from JWT claims in handler
}

// EmptyCartReq 清空购物车请求参数 (UserID will be from JWT)
type EmptyCartReq struct {
	// UserID is removed, will be extracted from JWT claims in handler
}

// AddItemReq 添加商品到购物车请求参数
type AddItemReq struct {
	// UserID is removed, will be extracted from JWT claims in handler
	ProductID uint `json:"product_id" binding:"required,gt=0"`
	Quantity  int  `json:"quantity" binding:"required,gt=0,lte=100"` // Max 100 per add
}

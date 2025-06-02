// types/order.go
package types

// CreateOrderReq 创建订单请求参数
type CreateOrderReq struct {
	UserCurrency  string         `json:"user_currency" binding:"required"`  // 用户货币
	Email         string         `json:"email" binding:"required"`          // 用户邮箱
	FirstName     string         `json:"first_name" binding:"required"`     // 名
	LastName      string         `json:"last_name" binding:"required"`      // 姓
	StreetAddress string         `json:"street_address" binding:"required"` // 街道地址
	City          string         `json:"city" binding:"required"`           // 城市
	State         string         `json:"state" binding:"required"`          // 省/州
	Country       string         `json:"country" binding:"required"`        // 国家
	ZipCode       string         `json:"zip_code" binding:"required"`       // 邮政编码
	OrderItems    []OrderItemReq `json:"order_items" binding:"required"`    // 订单项
}

// OrderItemReq 订单项请求参数
type OrderItemReq struct {
	ProductID uint  `json:"product_id" binding:"required"` // 商品ID
	Quantity  int32 `json:"quantity" binding:"required"`   // 商品数量
}

// UpdateOrderReq 修改订单请求参数
type UpdateOrderReq struct {
	OrderID       string `json:"order_id" binding:"required"` // 订单ID
	UserID        uint   `json:"user_id" binding:"required"`  // 用户ID
	StreetAddress string `json:"street_address"`              // 街道地址（可选）
	City          string `json:"city"`                        // 城市（可选）
	State         string `json:"state"`                       // 省/州（可选）
	Country       string `json:"country"`                     // 国家（可选）
	ZipCode       string `json:"zip_code"`                    // 邮政编码（可选）
	Status        string `json:"status"`                      // 订单状态（可选）
}

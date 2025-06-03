// types/order.go
package types

// CreateOrderReq 创建订单请求参数
type CreateOrderReq struct {
	Items     []OrderItemReq `json:"items" binding:"required,dive"`      // dive validates each item in slice
	AddressID uint           `json:"address_id" binding:"required,gt=0"` // Assuming AddressID is for shipping
	// UserCurrency, Email, FirstName, etc. might be associated with AddressID or fetched for the user
}

// OrderItemReq 订单项请求参数
type OrderItemReq struct {
	ProductID uint `json:"product_id" binding:"required,gt=0"`
	Quantity  int  `json:"quantity" binding:"required,gt=0,lte=100"`
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

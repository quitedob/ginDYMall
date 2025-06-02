package types

// CheckoutReq 结算请求参数
type CheckoutReq struct {
	UserCurrency  string         `json:"user_currency" binding:"required"`  // 用户货币
	Email         string         `json:"email" binding:"required"`          // 用户邮箱
	StreetAddress string         `json:"street_address" binding:"required"` // 街道地址
	City          string         `json:"city" binding:"required"`           // 城市
	State         string         `json:"state" binding:"required"`          // 省/州
	Country       string         `json:"country" binding:"required"`        // 国家
	ZipCode       string         `json:"zip_code" binding:"required"`       // 邮政编码
	OrderItems    []OrderItemReq `json:"order_items" binding:"required"`    // 订单项
}

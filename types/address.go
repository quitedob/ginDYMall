// Package types 定义了抖音商城项目中与数据交互相关的请求与响应结构体
// 本文件主要用于地址管理服务（例如订单结算时的地址信息）
package types

// AddressCreateReq 创建地址请求结构体，用于新增地址信息
type AddressCreateReq struct {
	StreetAddress string `form:"street_address" json:"street_address"` // 街道地址，例如 "长安街100号"
	City          string `form:"city" json:"city"`                     // 城市名称，例如 "北京"
	State         string `form:"state" json:"state"`                   // 省/州名称，例如 "北京市"
	Country       string `form:"country" json:"country"`               // 国家名称，例如 "中国"
	ZipCode       string `form:"zip_code" json:"zip_code"`             // 邮政编码，例如 "100000"
	// 以上字段用于构建完整的地址信息，便于后续订单配送与结算
}

// AddressUpdateReq 更新地址请求结构体，用于修改已有地址信息
type AddressUpdateReq struct {
	ID            uint   `form:"id" json:"id"`                         // 地址ID，唯一标识一条地址记录
	StreetAddress string `form:"street_address" json:"street_address"` // 街道地址
	City          string `form:"city" json:"city"`                     // 城市名称
	State         string `form:"state" json:"state"`                   // 省/州名称
	Country       string `form:"country" json:"country"`               // 国家名称
	ZipCode       string `form:"zip_code" json:"zip_code"`             // 邮政编码
	// 更新时必须指定 ID，其他字段为可修改的地址信息
}

// AddressGetReq 获取地址详情请求结构体，根据地址ID查询具体地址信息
type AddressGetReq struct {
	ID uint `form:"id" json:"id"` // 地址ID，用于标识查询哪条地址记录
}

// AddressDeleteReq 删除地址请求结构体，根据地址ID删除指定地址记录
type AddressDeleteReq struct {
	ID uint `form:"id" json:"id"` // 地址ID，删除时需要提供
}

// AddressListReq 地址列表查询请求结构体，支持分页查询
type AddressListReq struct {
	BasePage // 内嵌分页基础结构体，便于计算数据偏移量
	// 如前面所述，若 page = 2 且 pageSize = 10，则 offset = (2 - 1) * 10 = 10
}

// AddressResp 地址响应结构体，用于返回地址信息给客户端
type AddressResp struct {
	ID            uint   `json:"id"`             // 地址记录的唯一标识
	UserID        uint   `json:"user_id"`        // 关联的用户ID，表示该地址属于哪个用户
	StreetAddress string `json:"street_address"` // 街道地址
	City          string `json:"city"`           // 城市名称
	State         string `json:"state"`          // 省/州名称
	Country       string `json:"country"`        // 国家名称
	ZipCode       string `json:"zip_code"`       // 邮政编码
	CreatedAt     int64  `json:"created_at"`     // 记录创建时间，使用 Unix 时间戳表示
	// 以上信息用于在用户中心或订单结算时展示完整的地址详情
}

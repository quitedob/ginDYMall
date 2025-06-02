// Package types 定义了抖音商城项目中各模块的数据传输对象
// 本文件主要用于分页请求、带总数的数据结构及其他基础功能
package types

// BasePage 分页基础结构体，用于请求分页相关的参数
// 这里的 PageNum 表示当前页码，PageSize 表示每页的记录数
// 用于所有需要分页查询的接口
type BasePage struct {
	PageNum  int `form:"page_num" json:"page_num"`   // 当前页码，默认从1开始，必须传入
	PageSize int `form:"page_size" json:"page_size"` // 每页记录数，默认每页显示10条数据
}

// DataListResp 带有总数的分页数据响应结构体，用于返回带有分页信息的数据列表
// 包含 Item（具体数据项）和 Total（总记录数）
type DataListResp struct {
	Item  interface{} `json:"item"`  // 数据项，可以是任意类型，通常是查询到的具体数据
	Total int64       `json:"total"` // 总记录数，用于分页功能展示总页数等
}

// 地址信息
type Address struct {
	StreetAddress string `json:"street_address"` // 街道地址
	City          string `json:"city"`           // 城市
	State         string `json:"state"`          // 省/州
	Country       string `json:"country"`        // 国家
	ZipCode       string `json:"zip_code"`       // 邮政编码
}

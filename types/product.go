// types/product.go
package types

// 商品
type Product struct {
	ID          uint32   `json:"id"`          // 商品ID
	Name        string   `json:"name"`        // 商品名称
	Description string   `json:"description"` // 商品描述
	Picture     string   `json:"picture"`     // 商品图片
	Price       float64  `json:"price"`       // 商品价格
	Stock       int      `json:"stock"`       // 商品库存
	Version     int      `json:"version"`     // 版本号
	Categories  []string `json:"categories"`  // 商品分类
}

// 查询商品请求
type GetProductReq struct {
	ID uint32 `json:"id"` // 商品ID
}

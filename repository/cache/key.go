package cache

import (
	"fmt"
	"strconv"
)

const (
	// RankKey 每日排名的Redis键名
	RankKey = "rank"
	// SkillProductKey 单个商品的Redis键名模板，%d为商品ID占位符
	SkillProductKey = "skill:product:%d"
	// SkillProductListKey 所有商品的Redis键名
	SkillProductListKey = "skill:product_list"
	// SkillProductUserKey 用户相关的商品信息Redis键名模板，%s为用户ID占位符
	SkillProductUserKey = "skill:user:%s"
)

// ProductViewKey 返回指定商品ID的查看数Redis键名
func ProductViewKey(id uint) string {
	// 将商品ID转为字符串并格式化为Redis键名
	return fmt.Sprintf("view:product:%s", strconv.Itoa(int(id)))
}

// ProductDetailKey returns the Redis key for a single product's details.
// Example: "product:detail:123"
func ProductDetailKey(id uint) string {
	return fmt.Sprintf("product:detail:%d", id)
}

// ProductListKey returns the Redis key for a paginated list of products.
// Example: "product:list:1:10" (page 1, size 10)
func ProductListKey(page, size int) string {
	return fmt.Sprintf("product:list:%d:%d", page, size)
}

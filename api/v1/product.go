// api/v1/product.go
package v1

import (
	"douyin/repository/db/dao"
	"douyin/repository/db/model"
	"douyin/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 创建商品接口
func CreateProduct(c *gin.Context) {

	var req types.Product
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效输入"})
		return
	}

	// 将 types.Product 转换为 model.Product
	modelProduct := &model.Product{
		Name:        req.Name,
		Description: req.Description,
		Picture:     req.Picture,
		Price:       req.Price,
	}

	// 创建商品并保存
	if err := dao.CreateProduct(modelProduct); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "创建商品时出错"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "商品创建成功"})
}

// 查询单个商品信息接口
func GetProduct(c *gin.Context) {

	var req types.GetProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效输入"})
		return
	}

	// 查询商品
	product, err := dao.GetProduct(req.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "商品未找到"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": product})
}

// 修改商品信息接口
func UpdateProduct(c *gin.Context) {

	var req types.Product
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效输入"})
		return
	}

	// 将 types.Product 转换为 model.Product
	modelProduct := &model.Product{
		ID:          uint(req.ID), // 确保 ID 转换为正确的类型
		Name:        req.Name,
		Description: req.Description,
		Picture:     req.Picture,
		Price:       req.Price,
	}

	// 修改商品信息
	if err := dao.UpdateProduct(modelProduct); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "更新商品时出错"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "商品更新成功"})
}

// 删除商品接口
func DeleteProduct(c *gin.Context) {

	var req types.GetProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效输入"})
		return
	}

	// 删除商品
	if err := dao.DeleteProduct(req.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "删除商品时出错"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "商品删除成功"})
}

// 查询商品列表接口
func ListProducts(c *gin.Context) {

	var req types.BasePage
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效输入"})
		return
	}

	// 查询商品列表
	products, total, err := dao.ListProducts(req.PageNum, req.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "获取商品列表时出错"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"total":    total,
	})
}

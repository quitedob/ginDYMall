// api/v1/product.go
package v1

import (
	"douyin/pkg/utils/ctl" // Assuming GetUserID gets current user ID
	"douyin/pkg/utils/response"
	// "douyin/repository/db/dao" // DAO calls are now in service layer
	// "douyin/repository/db/model" // Model conversions are now in service layer
	"douyin/service"
	"douyin/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv" // For parsing ID from path if needed, or use ShouldBindUri
)

// 创建商品接口
func CreateProduct(c *gin.Context) {
	var req types.Product
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Fail(1001, "参数非法："+err.Error()))
		return
	}

	// Assuming userID is needed for CreateProduct service call, get it from context
	userID, err := ctl.GetUserID(c) // Example: GetUserID extracts from JWT
	if err != nil {
		_ = c.Error(err) // Let error handler deal with it
		return
	}

	// 调用服务层创建商品
	if err := service.CreateProduct(c.Request.Context(), userID, &req); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.Success("商品创建成功"))
}

// 查询单个商品信息接口
func GetProduct(c *gin.Context) {
	// Assuming product ID is passed as a path parameter, e.g., /product/:id
	// Or if it's a query param, use c.Query("id")
	// For this example, let's assume it's a request body for consistency with current code
	var req types.GetProductReq
	if err := c.ShouldBindJSON(&req); err != nil { // If ID is in path: c.ShouldBindUri(&uriReq)
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Fail(1001, "参数非法："+err.Error()))
		return
	}

	// 查询商品
	product, err := service.GetProductByID(c.Request.Context(), req.ID) // req.ID is uint32
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.Success(product))
}

// 修改商品信息接口
func UpdateProduct(c *gin.Context) {
	var req types.Product // Assuming types.Product contains ID for update
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Fail(1001, "参数非法："+err.Error()))
		return
	}

	userID, err := ctl.GetUserID(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 修改商品信息
	if err := service.UpdateProduct(c.Request.Context(), userID, &req); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.Success("商品更新成功"))
}

// 删除商品接口
func DeleteProduct(c *gin.Context) {
	var req types.GetProductReq // Assuming this just needs product ID
	if err := c.ShouldBindJSON(&req); err != nil { // Or c.ShouldBindUri if ID in path
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Fail(1001, "参数非法："+err.Error()))
		return
	}

	userID, err := ctl.GetUserID(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 删除商品
	if err := service.DeleteProduct(c.Request.Context(), userID, req.ID); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.Success("商品删除成功"))
}

// 查询商品列表接口
func ListProducts(c *gin.Context) {
	var req types.BasePage
	if err := c.ShouldBindQuery(&req); err != nil { // Query params for pagination
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Fail(1001, "参数非法："+err.Error()))
		return
	}

	// 查询商品列表
	products, total, err := service.ListProducts(c.Request.Context(), req.PageNum, req.PageSize)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"products": products,
		"total":    total,
	}))
}

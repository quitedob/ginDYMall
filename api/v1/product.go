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

// @Summary      创建商品
// @Description  新增一条商品记录. 需要用户认证.
// @Tags         商品 (Product)
// @Accept       json
// @Produce      json
// @Param        data  body      types.Product        true  "商品信息 (Product Information)"
// @Success      200   {object}  response.APIResponse{message=string} "创建成功"
// @Failure      400   {object}  response.APIResponse "参数校验失败 (Bad Request)"
// @Failure      401   {object}  response.APIResponse "用户未认证 (Unauthorized)"
// @Failure      500   {object}  response.APIResponse "服务器内部错误 (Internal Server Error)"
// @Router       /product [post] // Assuming base path /api/v1 is set globally
func CreateProduct(c *gin.Context) {
	var req types.Product
	if err := c.ShouldBindJSON(&req); err != nil {
		// Using response.Fail directly as per task example, though current code uses AbortWithStatusJSON
		response.Fail(c, http.StatusBadRequest, "参数非法: "+err.Error())
		return
	}

	// Assuming userID is needed for CreateProduct service call, get it from context
	userID, err := ctl.GetUserID(c) // Example: GetUserID extracts from JWT
	if err != nil {
		// response.Fail(c, http.StatusUnauthorized, "用户未认证: "+err.Error()) // Example for Swaggo
		_ = c.Error(err) // Let error handler deal with it
		return
	}

	// 调用服务层创建商品
	ctx := c.Request.Context()
	if err := service.CreateProduct(ctx, userID, &req); err != nil {
		// response.Fail(c, http.StatusInternalServerError, "创建商品失败: "+err.Error()) // Example for Swaggo
		_ = c.Error(err)
		return
	}

	response.Success(c, "商品创建成功") // Adjusted to match task's response.Success
}

// @Summary      获取商品
// @Description  根据 ID 查询单个商品详情. ID 在请求体中提供.
// @Tags         商品 (Product)
// @Accept       json
// @Produce      json
// @Param        data  body      types.GetProductReq  true  "商品 ID 请求 (Product ID Request)"
// @Success      200   {object}  response.APIResponse{data=types.Product} "返回商品详情"
// @Failure      400   {object}  response.APIResponse "参数错误 (Bad Request - e.g., invalid ID format or missing ID)"
// @Failure      404   {object}  response.APIResponse "商品不存在 (Not Found)"
// @Failure      500   {object}  response.APIResponse "服务器内部错误 (Internal Server Error)"
// @Router       /product/detail [post] // Changed to POST as ID is in body, or use GET with query param for /product
func GetProduct(c *gin.Context) {
	var req types.GetProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数非法: "+err.Error())
		return
	}

	// 查询商品
	ctx := c.Request.Context()
	product, err := service.GetProductByID(ctx, req.ID) // req.ID is uint32
	if err != nil {
		// Handle specific errors like not found if service.GetProductByID provides them
		// For example: if errors.Is(err, gorm.ErrRecordNotFound) { ... }
		// response.Fail(c, http.StatusNotFound, "商品未找到") // Example for Swaggo
		_ = c.Error(err) // Let error handler deal with it
		return
	}
	// Ensure product is not nil if error is nil, though service should guarantee this
	if product == nil {
		response.Fail(c, http.StatusNotFound, "商品未找到")
        return
	}

	response.Success(c, product)
}

// @Summary      修改商品信息
// @Description  修改现有商品的信息. 需要用户认证. 商品ID应在请求体中.
// @Tags         商品 (Product)
// @Accept       json
// @Produce      json
// @Param        data  body      types.Product        true  "要更新的商品信息 (Product Information to Update)"
// @Success      200   {object}  response.APIResponse{message=string} "更新成功"
// @Failure      400   {object}  response.APIResponse "参数校验失败 (Bad Request)"
// @Failure      401   {object}  response.APIResponse "用户未认证 (Unauthorized)"
// @Failure      404   {object}  response.APIResponse "商品不存在 (Not Found)"
// @Failure      500   {object}  response.APIResponse "服务器内部错误 (Internal Server Error)"
// @Router       /product [put] // Assuming base path /api/v1 is set globally
func UpdateProduct(c *gin.Context) {
	var req types.Product // Assuming types.Product contains ID for update
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数非法: "+err.Error())
		return
	}

	userID, err := ctl.GetUserID(c)
	if err != nil {
		// response.Fail(c, http.StatusUnauthorized, "用户未认证: "+err.Error()) // Example for Swaggo
		_ = c.Error(err)
		return
	}

	// 修改商品信息
	ctx := c.Request.Context()
	if err := service.UpdateProduct(ctx, userID, &req); err != nil {
		// response.Fail(c, http.StatusInternalServerError, "更新商品失败: "+err.Error()) // Example for Swaggo
		// Handle not found error from service if applicable
		_ = c.Error(err)
		return
	}

	response.Success(c, "商品更新成功")
}

// @Summary      删除商品
// @Description  根据 ID 删除商品. ID 在请求体中提供. 需要用户认证.
// @Tags         商品 (Product)
// @Accept       json
// @Produce      json
// @Param        data  body      types.GetProductReq  true  "商品 ID 请求 (Product ID Request)"
// @Success      200   {object}  response.APIResponse{message=string} "删除成功"
// @Failure      400   {object}  response.APIResponse "参数错误 (Bad Request)"
// @Failure      401   {object}  response.APIResponse "用户未认证 (Unauthorized)"
// @Failure      404   {object}  response.APIResponse "商品不存在 (Not Found)"
// @Failure      500   {object}  response.APIResponse "服务器内部错误 (Internal Server Error)"
// @Router       /product [delete] // Assuming base path /api/v1 is set globally
func DeleteProduct(c *gin.Context) {
	var req types.GetProductReq // Assuming this just needs product ID
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数非法: "+err.Error())
		return
	}

	userID, err := ctl.GetUserID(c)
	if err != nil {
		// response.Fail(c, http.StatusUnauthorized, "用户未认证: "+err.Error()) // Example for Swaggo
		_ = c.Error(err)
		return
	}

	// 删除商品
	ctx := c.Request.Context()
	if err := service.DeleteProduct(ctx, userID, req.ID); err != nil {
		// response.Fail(c, http.StatusInternalServerError, "删除商品失败: "+err.Error()) // Example for Swaggo
		// Handle not found error from service if applicable
		_ = c.Error(err)
		return
	}

	response.Success(c, "商品删除成功")
}

// @Summary      查询商品列表
// @Description  分页获取商品列表.
// @Tags         商品 (Product)
// @Accept       json
// @Produce      json
// @Param        pageNum   query     int                  false "页码 (Page Number)" default(1)
// @Param        pageSize  query     int                  false "每页数量 (Page Size)" default(10)
// @Success      200   {object}  response.APIResponse{data=object{products=[]types.Product,total=int}} "返回商品列表和总数"
// @Failure      400   {object}  response.APIResponse "参数错误 (Bad Request)"
// @Failure      500   {object}  response.APIResponse "服务器内部错误 (Internal Server Error)"
// @Router       /product/list [get] // Or /product [get] if it's the standard list endpoint
func ListProducts(c *gin.Context) {
	var req types.BasePage
	if err := c.ShouldBindQuery(&req); err != nil { // Query params for pagination
		response.Fail(c, http.StatusBadRequest, "参数非法: "+err.Error())
		return
	}

	// Set defaults for pagination if not provided or invalid
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}


	// 查询商品列表
	ctx := c.Request.Context()
	products, total, err := service.ListProducts(ctx, req.PageNum, req.PageSize)
	if err != nil {
		// response.Fail(c, http.StatusInternalServerError, "查询商品列表失败: "+err.Error()) // Example for Swaggo
		_ = c.Error(err)
		return
	}

	response.Success(c, gin.H{
		"products": products,
		"total":    total,
	})
}

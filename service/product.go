// service/product.go
package service

import (
	"douyin/repository/db/dao"
	"douyin/repository/db/model"
	"douyin/types"
	"errors"
	"fmt"
	"log"
)

func CreateProduct(userID uint32, product *types.Product) error {
	// 1. 验证用户身份是否有效，或者是否有权限创建商品
	if userID == 0 {
		return errors.New("用户身份无效")
	}

	// 2. 将 types.Product 转换为 model.Product
	modelProduct := &model.Product{
		Name:        product.Name,
		Description: product.Description,
		Picture:     product.Picture,
		Price:       product.Price,
	}

	// 3. 调用DAO层，进行商品创建
	err := dao.CreateProduct(modelProduct)
	if err != nil {
		log.Printf("创建商品失败：%v", err)
		return err
	}
	fmt.Println("商品创建成功")
	return nil
}

func GetProductByID(id uint32) (*types.Product, error) {
	// 调用DAO层查询商品
	product, err := dao.GetProduct(id)
	if err != nil {
		log.Printf("查询商品失败：%v", err)
		return nil, err
	}

	// 将 model.Product 转换为 types.Product
	typesProduct := &types.Product{
		ID:          uint32(product.ID),
		Name:        product.Name,
		Description: product.Description,
		Picture:     product.Picture,
		Price:       product.Price,
	}

	return typesProduct, nil
}
func ListProducts(pageNum, pageSize int) ([]types.Product, int64, error) {
	// 调用DAO层查询商品列表
	products, total, err := dao.ListProducts(pageNum, pageSize)
	if err != nil {
		log.Printf("获取商品列表失败：%v", err)
		return nil, 0, err
	}

	// 转换 model.Product 为 types.Product
	var result []types.Product
	for _, p := range products {
		result = append(result, types.Product{
			ID:          uint32(p.ID),
			Name:        p.Name,
			Description: p.Description,
			Picture:     p.Picture,
			Price:       p.Price,
		})
	}

	return result, total, nil
}

func UpdateProduct(userID uint32, product *types.Product) error {
	// 验证用户身份或权限
	if userID == 0 {
		return errors.New("用户身份无效")
	}

	// 2. 将 types.Product 转换为 model.Product
	modelProduct := &model.Product{
		ID:          uint(product.ID),
		Name:        product.Name,
		Description: product.Description,
		Picture:     product.Picture,
		Price:       product.Price,
	}

	// 调用DAO层修改商品
	err := dao.UpdateProduct(modelProduct) // 确保这里传递的是 *model.Product 类型
	if err != nil {
		log.Printf("修改商品失败：%v", err)
		return err
	}

	fmt.Println("商品信息修改成功")
	return nil
}

func DeleteProduct(userID uint32, productID uint32) error {
	// 验证用户身份或权限
	if userID == 0 {
		return errors.New("用户身份无效")
	}

	// 调用DAO层删除商品
	err := dao.DeleteProduct(uint32(productID)) // 如果 DAO 需要 int 类型，这里做类型转换
	if err != nil {
		log.Printf("删除商品失败：%v", err)
		return err
	}

	fmt.Println("商品删除成功")
	return nil
}

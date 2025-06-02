package dao

import (
	"errors"
	"fmt"

	"douyin/repository/db/model"
	"gorm.io/gorm"
)

// CartDao 购物车数据访问对象，封装对 cart_items 表的增删改查操作
type CartDao struct {
	db *gorm.DB
}

// NewCartDao 创建并返回一个新的 CartDao 实例
func NewCartDao(db *gorm.DB) *CartDao {
	return &CartDao{
		db: db,
	}
}

// CreateCart 创建空购物车（实际上不做任何插入操作，表示用户在 cart_items 中还没有商品）
func (dao *CartDao) CreateCart(userID uint) error {
	// 对于只有 cart_items 一张表的情况，“创建购物车”意味着此时用户暂时没有任何商品记录。
	// 如果有需要，可以在此插入一些默认记录；当前示例中不插入，返回nil即可。
	fmt.Println("CreateCart: 创建空购物车成功（不插入任何记录）")
	return nil
}

// GetCart 获取用户的购物车信息（即获取 cart_items 表中该用户所有商品项）
func (dao *CartDao) GetCart(userID uint) ([]model.CartItem, error) {
	var items []model.CartItem
	if err := dao.db.Where("user_id = ?", userID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// EmptyCart 清空购物车（删除 cart_items 表里所有匹配 user_id 的记录）
func (dao *CartDao) EmptyCart(userID uint) error {
	if err := dao.db.Where("user_id = ?", userID).Delete(&model.CartItem{}).Error; err != nil {
		return err
	}
	fmt.Printf("EmptyCart: 已清空用户 %d 的购物车\n", userID)
	return nil
}

// AddItem 往购物车中添加(或更新)商品
// userID: 用户ID
// productID: 商品ID
// quantity: 要添加的商品数量（可能为正数，表示往购物车中增加）
func (dao *CartDao) AddItem(userID, productID uint, quantity int32) error {
	var cartItem model.CartItem

	// 先查询cart_items中是否已有此商品
	err := dao.db.Where("user_id = ? AND product_id = ?", userID, productID).
		First(&cartItem).Error

	if err != nil {
		// 如果记录未找到，表示购物车里还没有这个商品，执行插入操作
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newItem := model.CartItem{
				UserID:    userID,
				ProductID: productID,
				Quantity:  quantity,
			}
			if createErr := dao.db.Create(&newItem).Error; createErr != nil {
				return createErr
			}
			fmt.Printf("AddItem: 用户 %d 的购物车中新增商品 %d，数量为 %d\n", userID, productID, quantity)
			return nil
		}
		// 如果是其他错误，则直接返回
		return err
	}

	// 如果购物车已存在该商品，则进行数量累加
	cartItem.Quantity += quantity

	// 如果因为业务需要可限制最大数量或最小数量，可在此处做校验
	if cartItem.Quantity < 1 {
		// 如果数量不合法（例如 < 1），也可根据业务逻辑决定是否直接删除该条记录，或报错返回
		return fmt.Errorf("AddItem: 数量小于1，操作非法或请使用删除接口处理")
	}

	// 保存更新后的数量
	if saveErr := dao.db.Save(&cartItem).Error; saveErr != nil {
		return saveErr
	}
	fmt.Printf("AddItem: 用户 %d 的购物车中更新商品 %d，数量已变更为 %d\n", userID, productID, cartItem.Quantity)
	return nil
}

package dao

import (
	"context"
	"douyin/repository/db/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
func (dao *CartDao) CreateCart(ctx context.Context, userID uint) error {
	// 对于只有 cart_items 一张表的情况，“创建购物车”意味着此时用户暂时没有任何商品记录。
	// 如果有需要，可以在此插入一些默认记录；当前示例中不插入，返回nil即可。
	fmt.Println("CreateCart: 创建空购物车成功（不插入任何记录）")
	return nil
}

// GetCart 获取用户的购物车信息（即获取 cart_items 表中该用户所有商品项）
func (dao *CartDao) GetCart(ctx context.Context, userID uint) ([]model.CartItem, error) {
	var items []model.CartItem
	if err := dao.db.WithContext(ctx).Where("user_id = ?", userID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// EmptyCart 清空购物车（删除 cart_items 表里所有匹配 user_id 的记录）
func (dao *CartDao) EmptyCart(ctx context.Context, userID uint) error {
	if err := dao.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&model.CartItem{}).Error; err != nil {
		return err
	}
	fmt.Printf("EmptyCart: 已清空用户 %d 的购物车\n", userID)
	return nil
}

// AddItem 往购物车中添加(或更新)商品
// userID: 用户ID
// productID: 商品ID
// quantity: 要添加的商品数量（可能为正数，表示往购物车中增加）
func (dao *CartDao) AddItem(ctx context.Context, userID, productID uint, quantity int32) error {
	tx := dao.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	var product model.Product
	// Lock product row for update
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", productID).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("商品不存在")
		}
		return err
	}

	// Check stock - assuming adding to cart deducts stock
	if product.Stock < int(quantity) { // quantity is int32, product.Stock is int
		return errors.New("库存不足: " + product.Name)
	}

	// Update stock using optimistic locking
	result := tx.Model(&model.Product{}).Where("id = ? AND version = ?", product.ID, product.Version).Updates(map[string]interface{}{
		"stock":   gorm.Expr("stock - ?", quantity),
		"version": gorm.Expr("version + 1"),
	})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("并发冲突，请重试: " + product.Name)
	}

	var cartItem model.CartItem
	// 先查询cart_items中是否已有此商品
	err := tx.Where("user_id = ? AND product_id = ?", userID, productID).
		First(&cartItem).Error

	if err != nil {
		// 如果记录未找到，表示购物车里还没有这个商品，执行插入操作
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newItem := model.CartItem{
				UserID:    userID,
				ProductID: productID,
				Quantity:  quantity,
			}
			if createErr := tx.Create(&newItem).Error; createErr != nil {
				return createErr
			}
			fmt.Printf("AddItem: 用户 %d 的购物车中新增商品 %d，数量为 %d\n", userID, productID, quantity)
		} else {
			// 如果是其他错误，则直接返回
			return err
		}
	} else {
		// 如果购物车已存在该商品，则进行数量累加
		// Note: Stock was already deducted for the new quantity. If the item exists,
		// this logic might need adjustment if stock is only deducted on checkout.
		// For now, assuming stock is deducted when adding to cart.
		cartItem.Quantity += quantity
		if cartItem.Quantity < 1 {
			// If quantity becomes less than 1, consider deleting the item or erroring
			// For now, let's assume quantity will always be positive when adding.
			// If logic allows reducing quantity which might make it < 1, handle appropriately.
			// e.g., tx.Delete(&cartItem) or return error
			return fmt.Errorf("AddItem: 最终商品数量小于1，操作非法")
		}
		if saveErr := tx.Save(&cartItem).Error; saveErr != nil {
			return saveErr
		}
		fmt.Printf("AddItem: 用户 %d 的购物车中更新商品 %d，数量已变更为 %d\n", userID, productID, cartItem.Quantity)
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

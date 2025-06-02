package dao

import (
	"douyin/repository/db/model"
	"douyin/types"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// CheckoutDao 定义结算数据访问对象
type CheckoutDao struct {
	db *gorm.DB
}

// NewCheckoutDao 根据传入的数据库连接创建新的 CheckoutDao 实例
func NewCheckoutDao(db *gorm.DB) *CheckoutDao {
	return &CheckoutDao{
		db: db,
	}
}

// CheckoutOrder 进行订单结算
// 此方法先创建订单记录，再创建支付记录，确保支付记录的订单外键引用存在
func (dao *CheckoutDao) CheckoutOrder(userID uint, req *types.CreateOrderReq) (string, error) {
	// 计算订单总金额
	var totalAmount float64
	for _, item := range req.OrderItems {
		var product model.Product
		if err := dao.db.Where("id = ?", item.ProductID).First(&product).Error; err != nil {
			return "", err
		}
		totalAmount += float64(item.Quantity) * product.Price
	}

	// 生成订单 ID 和支付交易 ID
	orderID := uuid.New().String()
	transactionID := uuid.New().String()

	// 开启事务
	tx := dao.db.Begin()

	// 插入订单记录
	order := model.Order{
		OrderID:       orderID,           // 订单ID
		UserID:        userID,            // 用户ID
		UserCurrency:  req.UserCurrency,  // 用户货币
		Email:         req.Email,         // 用户邮箱
		FirstName:     req.FirstName,     // 名（注意字段名大小写）
		LastName:      req.LastName,      // 姓
		StreetAddress: req.StreetAddress, // 街道地址
		City:          req.City,          // 城市
		State:         req.State,         // 省/州
		Country:       req.Country,       // 国家
		ZipCode:       req.ZipCode,       // 邮编
		CreatedAt:     time.Now(),        // 创建时间
		Status:        "paid",            // 订单状态设为已支付
	}
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		println("创建订单记录失败:", err.Error())
		return "", err
	}

	// 插入支付记录，使用相同的订单ID
	payment := model.Payment{
		TransactionID: transactionID, // 支付交易ID
		OrderID:       orderID,       // 关联的订单ID
		UserID:        userID,        // 用户ID
		Amount:        totalAmount,   // 支付金额
		Status:        "initiated",   // 初始支付状态
		CreatedAt:     time.Now(),    // 支付时间
	}
	if err := tx.Create(&payment).Error; err != nil {
		tx.Rollback()
		println("创建支付记录失败:", err.Error())
		return "", err
	}

	// 提交事务
	tx.Commit()
	println("订单结算成功，交易ID:", transactionID)

	return transactionID, nil
}

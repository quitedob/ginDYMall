package service

import (
	"douyin/pkg/utils/log"
	"context"
	"douyin/pkg/utils/log"
	"douyin/repository/db/dao"
	"douyin/types"
	"gorm.io/gorm"
)

// OrderService 订单服务
type OrderService struct {
	orderDao   *dao.OrderDao   // Renamed field for clarity
	addressDao *dao.AddressDao // Added AddressDao
	// productDao *dao.ProductDao // Might be needed if product logic moves here
}

// NewOrderService 创建新的 OrderService 实例
func NewOrderService(db *gorm.DB) (*OrderService, error) {
	// 获取底层的 *sql.DB 实例，并检查数据库连接是否成功
	sqlDB, err := db.DB()
	if err != nil {
		log.Errorf("获取数据库连接失败: %v", err)
		return nil, err
	}
	if err := sqlDB.Ping(); err != nil {
		log.Errorf("数据库连接失败: %v", err)
		return nil, err
	}
	log.Infof("订单服务使用的数据库连接成功")

	return &OrderService{
		orderDao:   dao.NewOrderDao(db),
		addressDao: dao.NewAddressDao(db), // Initialize AddressDao
	}, nil
}

// CreateOrder handles the business logic for creating an order.
// It fetches address details and then calls the order DAO.
// The transactionality is currently handled within s.orderDao.CreateOrder.
func (s *OrderService) CreateOrder(ctx context.Context, userID uint, addressID uint, items []types.OrderItemReq) (string, error) {
	// Fetch address details using addressID
	address, err := s.addressDao.GetAddressByID(ctx, userID, addressID)
	if err != nil {
		log.Errorf("获取地址失败 (userID: %d, addressID: %d): %v", userID, addressID, err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("无效的地址ID或地址不属于该用户")
		}
		return "", errors.New("获取地址信息时出错") // Generic error for other DB issues
	}

	// Adapt the service layer request (with AddressID and simplified items)
	// to the DAO layer's expected types.CreateOrderReq structure.
	daoOrderItems := make([]types.OrderItemReq, len(items))
	for i, item := range items {
		// Assuming the types.OrderItemReq for DAO is the same as for service input after previous refactor.
		// If DAO expects a different structure or type for Quantity (e.g. int32 vs int), adjust here.
		daoOrderItems[i] = types.OrderItemReq{
			ProductID: item.ProductID,
			Quantity:  item.Quantity, // This now uses the updated types.OrderItemReq with int Quantity
		}
	}

	// Populate the DAO request DTO with fetched address details.
	// UserCurrency might come from user profile or be determined by region/address.
	// For now, using placeholder or from address if available.
	userCurrency := "USD" // Placeholder, or determine from address.Country or user profile.
	// if address.UserCurrency != "" { userCurrency = address.UserCurrency }


	daoReq := &types.CreateOrderReq{ // This refers to the types.CreateOrderReq used by orderDao.CreateOrder
		UserCurrency:  userCurrency,
		Email:         address.Email, // Assuming Email is part of the address model, or fetch from user profile
		FirstName:     address.FirstName,
		LastName:      address.LastName,
		StreetAddress: address.StreetAddress,
		City:          address.City,
		State:         address.State,
		Country:       address.Country,
		ZipCode:       address.ZipCode,
		OrderItems:    daoOrderItems,
	}

	// Call the existing transactional DAO method
	log.Infof("Service CreateOrder calling DAO with userID: %d, using addressID: %d", userID, addressID)
	return s.orderDao.CreateOrder(ctx, userID, daoReq)
}

func (s *OrderService) UpdateOrder(ctx context.Context, userID uint, req *types.UpdateOrderReq) error {
	// The DAO's UpdateOrder signature was: func (dao *OrderDao) UpdateOrder(userID uint, req *types.UpdateOrderReq) error
	// It should be updated to accept context.
	// Assuming it's updated to: func (dao *OrderDao) UpdateOrder(ctx context.Context, userID uint, req *types.UpdateOrderReq) error
	// return s.orderDao.UpdateOrder(ctx, userID, req)

	// If DAO not updated yet:
	log.Warnf("Calling orderDao.UpdateOrder without context. Consider updating DAO if context propagation is needed.")
	return s.orderDao.UpdateOrder(userID, req) 
}

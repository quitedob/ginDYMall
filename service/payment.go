package service

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"gorm.io/gorm"

	"github.com/CocaineCong/gin-mall/consts"
	"github.com/CocaineCong/gin-mall/pkg/utils/ctl"
	"github.com/CocaineCong/gin-mall/pkg/utils/log"
	"github.com/CocaineCong/gin-mall/repository/db/dao"
	"github.com/CocaineCong/gin-mall/repository/db/model"
	"github.com/CocaineCong/gin-mall/types"
)

var PaymentSrvIns *PaymentSrv
var PaymentSrvOnce sync.Once

// PaymentSrv 支付服务层对象
type PaymentSrv struct {
}

// GetPaymentSrv 单例模式获取支付服务实例
func GetPaymentSrv() *PaymentSrv {
	PaymentSrvOnce.Do(func() {
		PaymentSrvIns = &PaymentSrv{}
	})
	return PaymentSrvIns
}

// PayDown 支付操作
// 该操作包括：扣除买家金币、增加卖家金币、扣减商品库存、更新订单状态、生成买家商品记录等
func (s *PaymentSrv) PayDown(ctx context.Context, req *types.PaymentDownReq) (resp interface{}, err error) {
	// 获取当前用户信息（买家信息）
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		log.LogrusObj.Error("获取用户信息失败：", err)
		return nil, err
	}

	// 开启事务处理支付流程
	err = dao.NewOrderDao(ctx).Transaction(func(tx *gorm.DB) error {
		uId := u.Id

		// 获取订单信息（支付订单）
		payment, err := dao.NewOrderDaoByDB(tx).GetOrderById(req.OrderId, uId)
		if err != nil {
			log.LogrusObj.Error("查询订单信息失败：", err)
			return err
		}

		// 计算订单总金额：单价 * 数量
		money := payment.Money
		num := payment.Num
		money = money * float64(num)

		// 获取买家信息，并解密余额
		userDao := dao.NewUserDaoByDB(tx)
		user, err := userDao.GetUserById(uId)
		if err != nil {
			log.LogrusObj.Error("查询买家信息失败：", err)
			return err
		}

		moneyFloat, err := user.DecryptMoney(req.Key)
		if err != nil {
			log.LogrusObj.Error("解密买家余额失败：", err)
			return err
		}

		// 检查买家余额是否充足
		if moneyFloat-money < 0.0 {
			log.LogrusObj.Error("买家余额不足，当前余额：", moneyFloat, " 订单金额：", money)
			return errors.New("金币不足")
		}

		// 更新买家余额：扣除订单金额
		newBuyerBalance := moneyFloat - money
		finMoney := fmt.Sprintf("%f", newBuyerBalance)
		user.Money = finMoney
		user.Money, err = user.EncryptMoney(req.Key)
		if err != nil {
			log.LogrusObj.Error("加密买家余额失败：", err)
			return err
		}

		err = userDao.UpdateUserById(uId, user)
		if err != nil {
			log.LogrusObj.Error("更新买家余额失败：", err)
			return err
		}

		// 获取卖家信息
		boss, err := userDao.GetUserById(uint(req.BossID))
		if err != nil {
			log.LogrusObj.Error("查询卖家信息失败：", err)
			return err
		}

		// 解密卖家余额，并增加订单金额
		bossBalance, err := boss.DecryptMoney(req.Key)
		if err != nil {
			log.LogrusObj.Error("解密卖家余额失败：", err)
			return err
		}

		newBossBalance := bossBalance + money
		finMoney = fmt.Sprintf("%f", newBossBalance)
		boss.Money = finMoney
		boss.Money, err = boss.EncryptMoney(req.Key)
		if err != nil {
			log.LogrusObj.Error("加密卖家余额失败：", err)
			return err
		}

		err = userDao.UpdateUserById(uint(req.BossID), boss)
		if err != nil {
			log.LogrusObj.Error("更新卖家余额失败：", err)
			return err
		}

		// 更新商品库存：扣减已购买数量
		productDao := dao.NewProductDaoByDB(tx)
		product, err := productDao.GetProductById(uint(req.ProductID))
		if err != nil {
			log.LogrusObj.Error("查询商品信息失败：", err)
			return err
		}
		if product.Num < num {
			log.LogrusObj.Error("商品库存不足，当前库存：", product.Num, " 购买数量：", num)
			return errors.New("商品库存不足")
		}
		product.Num -= num
		err = productDao.UpdateProduct(uint(req.ProductID), product)
		if err != nil {
			log.LogrusObj.Error("更新商品库存失败：", err)
			return err
		}

		// 更新订单状态为待发货
		payment.Type = consts.OrderTypePendingShipping
		err = dao.NewOrderDaoByDB(tx).UpdateOrderById(req.OrderId, uId, payment)
		if err != nil {
			log.LogrusObj.Error("更新订单状态失败：", err)
			return err
		}

		// 创建买家商品记录（将已购买的商品加入买家商品列表）
		productUser := model.Product{
			Name:          product.Name,
			CategoryID:    product.CategoryID,
			Title:         product.Title,
			Info:          product.Info,
			ImgPath:       product.ImgPath,
			Price:         product.Price,
			DiscountPrice: product.DiscountPrice,
			Num:           num,
			OnSale:        false,
			BossID:        uId,
			BossName:      user.UserName,
			BossAvatar:    user.Avatar,
		}

		err = productDao.CreateProduct(&productUser)
		if err != nil {
			log.LogrusObj.Error("创建买家商品记录失败：", err)
			return err
		}

		return nil
	})

	if err != nil {
		log.LogrusObj.Error("支付操作失败：", err)
		return nil, err
	}

	fmt.Println("支付操作成功") // 控制台输出中文提示
	return
}

// dao\product.go
package dao

import (
	"context"
	"douyin/config"
	"douyin/repository/db/model" // 商品模型包
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// init 在包初始化时建立数据库连接
func init() {
	// 加载配置文件
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		panic(fmt.Sprintf("加载配置失败：%v", err))
	}
	mysqlConfig := cfg.MySql.Default
	// 构造 DSN 字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		mysqlConfig.UserName,
		mysqlConfig.Password,
		mysqlConfig.DbHost,
		mysqlConfig.DbPort,
		mysqlConfig.DbName,
		mysqlConfig.Charset)
	// 连接数据库
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("连接数据库失败：%v", err))
	}
	fmt.Println("数据库连接成功！")
}

// ProductDao 定义商品数据访问对象结构体
type ProductDao struct {
	ctx context.Context // 当前请求上下文
	db  *gorm.DB        // 数据库连接对象
}

// NewProductDao 根据传入的上下文创建新的 ProductDao 实例
func NewProductDao(ctx context.Context) *ProductDao {
	return &ProductDao{
		ctx: ctx,
		db:  db,
	}
}

// 创建商品
func CreateProduct(product *model.Product) error {
	// 执行插入操作
	if err := db.Create(&product).Error; err != nil {
		fmt.Printf("创建商品时出错：%v\n", err)
		return err
	}
	fmt.Println("商品创建成功！")
	return nil
}

// 获取单个商品信息
func GetProduct(id uint32) (*model.Product, error) {
	var dbProduct model.Product
	// 根据商品ID查询商品
	if err := db.Where("id = ?", id).First(&dbProduct).Error; err != nil {
		fmt.Printf("查询商品时出错：%v\n", err)
		return nil, err
	}
	return &dbProduct, nil
}

// 获取商品列表（带分页）
func ListProducts(pageNum, pageSize int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	// 获取商品总数
	if err := db.Model(&model.Product{}).Count(&total).Error; err != nil {
		fmt.Printf("获取商品总数时出错：%v\n", err)
		return nil, 0, err
	}

	// 查询商品列表（分页）
	if err := db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&products).Error; err != nil {
		fmt.Printf("查询商品列表时出错：%v\n", err)
		return nil, 0, err
	}

	return products, total, nil
}

// 修改商品
func UpdateProduct(product *model.Product) error {
	// 执行更新操作，忽略 created_at 字段
	if err := db.Model(&model.Product{}).Where("id = ?", product.ID).Omit("created_at").Updates(product).Error; err != nil {
		fmt.Printf("更新商品时出错：%v\n", err)
		return err
	}
	fmt.Println("商品更新成功！")
	return nil
}

// 删除商品
func DeleteProduct(id uint32) error {
	// 执行删除操作
	if err := db.Where("id = ?", id).Delete(&model.Product{}).Error; err != nil {
		fmt.Printf("删除商品时出错：%v\n", err)
		return err
	}
	fmt.Println("商品删除成功！")
	return nil
}

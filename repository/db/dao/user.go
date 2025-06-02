// 文件：repository/db/dao/user.go
// 作用：定义用户数据访问层，操作 MySQL 数据库中 users 表
// 注意：查询条件中的字段名使用反引号包裹，确保与数据库建表一致

package dao

import (
	"context"
	"douyin/config"
	"douyin/repository/db/model" // 用户模型包
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 全局数据库连接对象
var db *gorm.DB

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

// UserDao 定义用户数据访问对象结构体
type UserDao struct {
	ctx context.Context // 当前请求上下文
	db  *gorm.DB        // 数据库连接对象
}

// NewUserDao 根据传入的上下文创建新的 UserDao 实例
func NewUserDao(ctx context.Context) *UserDao {
	return &UserDao{
		ctx: ctx,
		db:  db,
	}
}

// ExistOrNotByUserName 根据用户名查询用户是否存在
func (dao *UserDao) ExistOrNotByUserName(userName string) (*model.User, bool, error) {
	var user model.User
	if err := dao.db.Where("`user_name` = ?", userName).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &user, true, nil
}

// CreateUser 在数据库中创建新的用户记录
func (dao *UserDao) CreateUser(user *model.User) error {
	return dao.db.Create(user).Error
}

// GetUserById 根据用户ID查询用户详细信息
func (dao *UserDao) GetUserById(id uint) (*model.User, error) {
	var user model.User
	if err := dao.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUserById 根据用户ID更新用户记录
func (dao *UserDao) UpdateUserById(id uint, user *model.User) error {
	return dao.db.Model(&model.User{}).Where("id = ?", id).Updates(user).Error
}

// GetUserByUsername 根据用户名查询用户记录
func (dao *UserDao) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	if err := dao.db.Where("`user_name` = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail 根据邮箱查询用户记录（用于登录时）
func (dao *UserDao) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	if err := dao.db.Where("`email` = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

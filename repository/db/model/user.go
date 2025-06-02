package model

import (
	"strconv"
	"time"

	conf "douyin/config"
	"douyin/consts"
	"golang.org/x/crypto/bcrypt"
)

// User 用户模型（自定义字段，不使用 gorm.Model，以免引入 DeletedAt 字段）
type User struct {
	ID             uint      `gorm:"primaryKey"`                                 // 用户ID
	CreatedAt      time.Time `gorm:"column:created_at"`                          // 创建时间
	UpdatedAt      time.Time `gorm:"column:updated_at"`                          // 更新时间
	UserName       string    `gorm:"column:user_name;unique"`                    // 用户名
	Email          string    `gorm:"type:varchar(255);unique;not null"`          // 用户邮箱
	PasswordDigest string    `gorm:"column:password;type:varchar(255);not null"` // 加密后密码，对应数据库字段 "password"
	NickName       string    `gorm:"column:nick_name;type:varchar(255)"`         // 昵称
	Status         string    `gorm:"type:varchar(50);default:'active'"`          // 用户状态，默认激活
	Avatar         string    `gorm:"type:varchar(1000)"`                         // 头像
	Money          string    `gorm:"type:varchar(255)"`                          // 用户余额（直接存储字符串，不做加解密）
	Relations      []User    `gorm:"many2many:relation;"`                        // 用户之间的关系
}

const (
	PassWordCost = 12       // 密码加密难度
	Active       = "active" // 激活状态
)

// SetPassword 使用 bcrypt 对密码进行加密，并存入 PasswordDigest 字段
func (u *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PassWordCost)
	if err != nil {
		return err
	}
	u.PasswordDigest = string(bytes)
	return nil
}

// CheckPassword 校验密码是否正确
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordDigest), []byte(password))
	return err == nil
}

// AvatarURL 获取头像地址
func (u *User) AvatarURL() string {
	if conf.GlobalConfig.System.UploadModel == consts.UploadModelOss {
		return u.Avatar
	}
	pConfig := conf.GlobalConfig.PhotoPath
	return pConfig.PhotoHost + conf.GlobalConfig.System.HttpPort + pConfig.AvatarPath + u.Avatar
}

// EncryptMoney 返回原始余额，不进行加密处理
func (u *User) EncryptMoney(key string) (string, error) {
	// 本项目不对余额进行加密，直接返回余额
	return u.Money, nil
}

// DecryptMoney 直接将余额字符串转换为 float64，不进行解密处理
func (u *User) DecryptMoney(key string) (float64, error) {
	// 尝试将余额转换为 float64
	return strconv.ParseFloat(u.Money, 64)
}

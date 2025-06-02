package consts

import (
	"time"
)

const (
	EmailOperationBinding = iota + 1
	EmailOperationNoBinding
	EmailOperationUpdatePassword
)

var EmailOperationMap = map[uint]string{
	EmailOperationBinding:        "您正在绑定邮箱, 请点击链接确定身份 %s",
	EmailOperationNoBinding:      "您正在解邦邮箱, 请点击链接确定身份 %s",
	EmailOperationUpdatePassword: "您正在修改密码, 请点击链接校验身份 %s",
}

const (
	AccessTokenHeader    = "access_token"
	RefreshTokenHeader   = "refresh_token"
	HeaderForwardedProto = "X-Forwarded-Proto"
	MaxAge               = 3600 * 24
)

const (
	AccessToken = iota
	RefreshToken
)

const (
	AccessTokenExpireDuration  = 24 * time.Hour
	RefreshTokenExpireDuration = 10 * 24 * time.Hour
)

const EncryptMoneyKeyLength = 6

const UserInitMoney = "10000" // 初始金额

// MoneyDecryptKey 定义余额解密密钥（本项目余额不加密，不做解密，所以可为空）
const MoneyDecryptKey = ""

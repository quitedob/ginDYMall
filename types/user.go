package types // 定义 types 包，包含所有用户相关数据传输对象

// UserServiceReq 用户服务请求结构体（例如登录）
type UserServiceReq struct {
	NickName string `form:"nick_name" json:"nick_name"` // 用户昵称
	UserName string `form:"user_name" json:"user_name"` // 用户名（唯一）
	Password string `form:"password" json:"password"`   // 用户密码
}

type UserRegisterReq struct {
	NickName string `form:"nickname" json:"nickname"` // 用户昵称
	UserName string `form:"username" json:"username"` // 用户名
	Password string `form:"password" json:"password"` // 用户密码
	Email    string `form:"email" json:"email"`       // 用户邮箱
}

// UserTokenData 用户令牌数据结构，登录成功后返回 token 数据
type UserTokenData struct {
	User         interface{} `json:"user"`          // 用户信息对象
	AccessToken  string      `json:"access_token"`  // 访问令牌
	RefreshToken string      `json:"refresh_token"` // 刷新令牌
}

// UserLoginReq 用户登录请求结构体
type UserLoginReq struct {
	UserName string `form:"user_name" json:"user_name"` // 用户名
	Password string `form:"password" json:"password"`   // 用户密码
}

// UserInfoUpdateReq 用户信息更新请求结构体（更新昵称时使用）
type UserInfoUpdateReq struct {
	NickName string `form:"nick_name" json:"nick_name"` // 用户昵称
}

// UserUpdateReq 用户更新请求结构体（用于更新用户名和邮箱，仅更新这两个字段以及更新时间）
type UserUpdateReq struct {
	UserId   uint   `json:"user_id"`                    // 用户ID（必传，用于标识要更新的用户）
	UserName string `form:"user_name" json:"user_name"` // 新的用户名
	Email    string `form:"email" json:"email"`         // 新的邮箱
}

// UserInfoShowReq 用户信息展示请求结构体
type UserInfoShowReq struct {
	// 当前通过上下文获取用户信息，无需额外字段
}

// UserIdentityInfo 用户身份信息响应结构体，返回给客户端的详细信息
// 包括：用户ID、用户名、用户邮箱、余额、创建时间、更新时间
type UserIdentityInfo struct {
	UserID   uint   `json:"user_id"`   // 用户ID
	UserName string `json:"user_name"` // 用户名
	Email    string `json:"email"`     // 用户邮箱
	Money    string `json:"money"`     // 余额（直接展示，不进行加解密）
	CreateAt int64  `json:"create_at"` // 创建时间（Unix 时间戳）
	UpdateAt int64  `json:"update_at"` // 更新时间（Unix 时间戳）
}

// UserFollowingReq 用户关注请求结构体
type UserFollowingReq struct {
	Id uint `json:"id" form:"id"` // 被关注的用户 ID
}

// UserUnFollowingReq 用户取消关注请求结构体
type UserUnFollowingReq struct {
	Id uint `json:"id" form:"id"` // 取消关注的用户 ID
}

// SendEmailServiceReq 邮件服务请求结构体（如绑定邮箱、修改密码等）
type SendEmailServiceReq struct {
	Email         string `form:"email" json:"email"`                   // 用户邮箱
	Password      string `form:"password" json:"password"`             // 用户密码
	OperationType uint   `form:"operation_type" json:"operation_type"` // 操作类型
}

// ValidEmailServiceReq 邮箱验证请求结构体
type ValidEmailServiceReq struct {
	Token string `json:"token" form:"token"` // 邮箱验证令牌
}

// UserInfoResp 用户信息响应结构体，返回给客户端的用户详细信息
type UserInfoResp struct {
	ID       uint   `json:"id"`        // 用户 ID
	UserName string `json:"user_name"` // 用户名
	NickName string `json:"nickname"`  // 用户昵称
	Type     int    `json:"type"`      // 用户类型（例如普通用户、商家等）
	Email    string `json:"email"`     // 用户邮箱
	Status   string `json:"status"`    // 用户状态
	Avatar   string `json:"avatar"`    // 用户头像地址
	CreateAt int64  `json:"create_at"` // 用户创建时间（Unix 时间戳）
}

// UserChangePasswordReq 修改密码请求结构体
type UserChangePasswordReq struct {
	OldPassword string `form:"old_password" json:"old_password"` // 旧密码
	NewPassword string `form:"new_password" json:"new_password"` // 新密码
}

// UserChangeNicknameReq 修改昵称请求结构体
type UserChangeNicknameReq struct {
	NickName string `form:"nick_name" json:"nick_name"` // 新昵称
}

// UserSearchReq 用户查询请求结构体
type UserSearchReq struct {
	Username string `form:"username" json:"username"` // 查询时使用的用户名字段
}

// 文件：service/user.go
// 作用：实现用户业务逻辑，包括注册、登录、注销、修改密码、修改昵称、更新用户信息和展示用户身份信息

package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"douyin/consts"
	"douyin/pkg/utils/ctl"
	"douyin/pkg/utils/jwt"
	"douyin/pkg/utils/log"
	"douyin/repository/db/dao"
	"douyin/repository/db/model"
	"douyin/types"
)

var UserSrvIns *UserSrv   // 全局用户服务实例
var UserSrvOnce sync.Once // 确保只初始化一次

// UserSrv 定义用户服务层结构体
type UserSrv struct{}

// GetUserSrv 单例模式获取用户服务实例
func GetUserSrv() *UserSrv {
	UserSrvOnce.Do(func() {
		UserSrvIns = &UserSrv{}
	})
	return UserSrvIns
}

// UserRegister 用户注册业务逻辑
func (s *UserSrv) UserRegister(ctx context.Context, req *types.UserRegisterReq) (resp interface{}, err error) {
	userDao := dao.NewUserDao(ctx)
	// 检查用户名是否存在
	_, exist, err := userDao.ExistOrNotByUserName(req.UserName)
	if err != nil {
		log.LogrusObj.Error("查询用户是否存在失败：", err)
		return nil, err
	}
	if exist {
		err = errors.New("用户已存在")
		log.LogrusObj.Error("注册失败：", err)
		return nil, err
	}
	// 创建用户对象并设置初始字段（创建时间、更新时间取当前时间）
	user := &model.User{
		UserName:  req.UserName,
		NickName:  req.NickName,
		Email:     req.Email,
		Status:    model.Active,
		Money:     consts.UserInitMoney,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	// 加密密码
	if err = user.SetPassword(req.Password); err != nil {
		log.LogrusObj.Error("密码加密失败：", err)
		return nil, err
	}
	// 写入数据库
	if err = userDao.CreateUser(user); err != nil {
		log.LogrusObj.Error("创建用户失败：", err)
		return nil, err
	}
	fmt.Println("注册成功，用户创建成功")
	return user, nil
}

// UserLogin 用户登录业务逻辑（输入用户名和密码）
func (s *UserSrv) UserLogin(ctx context.Context, req *types.UserLoginReq) (resp interface{}, err error) {
	userDao := dao.NewUserDao(ctx)
	// 根据用户名查询用户记录
	user, exist, err := userDao.ExistOrNotByUserName(req.UserName)
	if err != nil {
		log.LogrusObj.Error("查询用户失败：", err)
		return nil, err
	}
	if !exist {
		log.LogrusObj.Error("登录失败，用户不存在")
		return nil, errors.New("用户不存在")
	}
	// 校验密码
	if !user.CheckPassword(req.Password) {
		log.LogrusObj.Error("登录失败，密码不正确")
		return nil, errors.New("账号/密码不正确")
	}
	// 生成 JWT 令牌，并存储到 Redis
	accessToken, refreshToken, err := jwt.GenerateToken(user.ID, user.UserName)
	if err != nil {
		log.LogrusObj.Error("生成令牌失败：", err)
		return nil, err
	}
	// 构造返回给前端的用户信息响应结构体
	userResp := &types.UserInfoResp{
		ID:       user.ID,
		UserName: user.UserName,
		NickName: user.NickName,
		Email:    user.Email,
		Status:   user.Status,
		CreateAt: user.CreatedAt.Unix(),
	}
	tokenData := types.UserTokenData{
		User:         userResp,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	fmt.Println("登录成功，生成令牌")
	return tokenData, nil
}

// UserLogout 用户注销业务逻辑（吊销 JWT 令牌，从 Redis 中删除 token）
func (s *UserSrv) UserLogout(ctx context.Context) error {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		log.LogrusObj.Error("获取用户信息失败：", err)
		return err
	}
	if err := jwt.RevokeToken(u.UserId); err != nil {
		log.LogrusObj.Error("注销失败：", err)
		return err
	}
	fmt.Println("注销成功")
	return nil
}

// UserChangePassword 修改密码业务逻辑
func (s *UserSrv) UserChangePassword(ctx context.Context, req *types.UserChangePasswordReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		log.LogrusObj.Error("获取用户信息失败：", err)
		return nil, err
	}
	userDao := dao.NewUserDao(ctx)
	user, err := userDao.GetUserById(u.UserId)
	if err != nil {
		log.LogrusObj.Error("查询用户失败：", err)
		return nil, err
	}
	if !user.CheckPassword(req.OldPassword) {
		log.LogrusObj.Error("旧密码不正确")
		return nil, errors.New("旧密码不正确")
	}
	if err = user.SetPassword(req.NewPassword); err != nil {
		log.LogrusObj.Error("密码加密失败：", err)
		return nil, err
	}
	user.UpdatedAt = time.Now()
	if err = userDao.UpdateUserById(u.UserId, user); err != nil {
		log.LogrusObj.Error("更新用户信息失败：", err)
		return nil, err
	}
	fmt.Println("密码修改成功")
	return "密码修改成功", nil
}

// UserChangeNickname 修改昵称业务逻辑
func (s *UserSrv) UserChangeNickname(ctx context.Context, req *types.UserChangeNicknameReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		log.LogrusObj.Error("获取用户信息失败：", err)
		return nil, err
	}
	userDao := dao.NewUserDao(ctx)
	user, err := userDao.GetUserById(u.UserId)
	if err != nil {
		log.LogrusObj.Error("查询用户失败：", err)
		return nil, err
	}
	user.NickName = req.NickName
	user.UpdatedAt = time.Now()
	if err = userDao.UpdateUserById(u.UserId, user); err != nil {
		log.LogrusObj.Error("更新昵称失败：", err)
		return nil, err
	}
	fmt.Println("昵称修改成功")
	return "昵称修改成功", nil
}

// UserUpdate 更新用户信息业务逻辑（仅更新 UserName, Email, UpdatedAt）
func (s *UserSrv) UserUpdate(ctx context.Context, req *types.UserUpdateReq) (resp interface{}, err error) {
	userDao := dao.NewUserDao(ctx)
	// 查询原始用户数据
	user, err := userDao.GetUserById(req.UserId)
	if err != nil {
		log.LogrusObj.Error("查询用户失败：", err)
		return nil, err
	}
	// 仅更新允许修改的字段
	user.UserName = req.UserName
	user.Email = req.Email
	user.UpdatedAt = time.Now()
	if err = userDao.UpdateUserById(req.UserId, user); err != nil {
		log.LogrusObj.Error("更新用户信息失败：", err)
		return nil, err
	}
	fmt.Println("用户更新成功")
	return "用户更新成功", nil
}

// UserInfoShow 获取用户身份信息业务逻辑
// 返回内容：[用户ID, 用户名, 用户邮箱, 用户余额, 创建时间, 更新时间]
// 余额直接返回，不进行加解密处理
func (s *UserSrv) UserInfoShow(ctx context.Context, req *types.UserInfoShowReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		log.LogrusObj.Error("获取用户信息失败：", err)
		return nil, err
	}
	userDao := dao.NewUserDao(ctx)
	user, err := userDao.GetUserById(u.UserId)
	if err != nil {
		log.LogrusObj.Error("查询用户失败：", err)
		return nil, err
	}
	userResp := &types.UserIdentityInfo{
		UserID:   user.ID,
		UserName: user.UserName,
		Email:    user.Email,
		Money:    user.Money,
		CreateAt: user.CreatedAt.Unix(),
		UpdateAt: user.UpdatedAt.Unix(),
	}
	fmt.Println("用户身份信息获取成功")
	return userResp, nil
}

package accountSrv

import (
	"context"
	"errors"

	"github.com/wethedevelop/account/cache"
	"github.com/wethedevelop/account/model"
	pb "github.com/wethedevelop/proto/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

const (
	// 账号或密码为空
	ACCOUNT_OR_PWD_EMPTY = "ACCOUNT_OR_PWD_EMPTY"
	// 已经被注册
	ACCOUNT_REGISTERED = "ACCOUNT_REGISTERED"
	// 账号或密码错误
	ACCOUNT_OR_PWD_NOT_MATCH = "ACCOUNT_OR_PWD_NOT_MATCH"
	// token无效
	ACCOUNT_INVALID_TOKEN = "ACCOUNT_INVALID_TOKEN"
)

// 实现第一方账号注册和登陆服务
type AccountServer struct {
	pb.UnsafeAccountAuthServer
}

// Signup 注册接口
func (s *AccountServer) Signup(context context.Context, in *pb.SignupRequest) (*pb.User, error) {
	account := in.GetAccount()
	password := in.GetPassword()
	// 不允许为空
	if account == "" || password == "" {
		return nil, status.Error(codes.InvalidArgument, ACCOUNT_OR_PWD_EMPTY)
	}
	// 检测重名
	checked, err := model.CheckRegistered(account)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if checked {
		return nil, status.Error(codes.InvalidArgument, ACCOUNT_REGISTERED)
	}
	user := model.User{
		Account:  account,
		Nickname: "沉默的开发者",
	}
	if err := user.SetPassword(password); err != nil {
		return nil, err
	}
	err = user.Create()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.User{Id: int64(user.ID), Nickname: user.Nickname}, nil
}

// Signin 登录接口
func (s *AccountServer) Signin(context context.Context, in *pb.SigninRequest) (*pb.TokenResponse, error) {
	account := in.GetAccount()
	password := in.GetPassword()

	var user model.User
	if err := model.DB.Where("account = ?", account).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.PermissionDenied, ACCOUNT_OR_PWD_NOT_MATCH)
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	if !user.CheckPassword(password) {
		return nil, status.Error(codes.PermissionDenied, ACCOUNT_OR_PWD_NOT_MATCH)
	}

	token, tokenExpire, err := user.MakeToken()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.TokenResponse{Token: token, TokenExpire: tokenExpire}, nil
}

// 通过Token验证用户身份
func (s *AccountServer) Check(context context.Context, in *pb.TokenRequest) (*pb.User, error) {
	uid, err := cache.GetUserByToken(in.Token)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if uid == "" {
		return nil, status.Error(codes.PermissionDenied, ACCOUNT_INVALID_TOKEN)
	}
	user, err := model.GetUser(uid)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if uid == "" {
		return nil, status.Error(codes.PermissionDenied, ACCOUNT_INVALID_TOKEN)
	}
	return &pb.User{ID: Nickname: user.Nickname}, nil
}

// 更新用户资料
func (s *AccountServer) Update(context.Context, *pb.UpdateRequest) (*pb.User, error) {
	return &pb.User{}, nil
}

// 通过用户ID获取用户资料
func (s *AccountServer) GetUser(context.Context, *pb.GetUserRequest) (*pb.User, error) {
	return &pb.User{}, nil
}

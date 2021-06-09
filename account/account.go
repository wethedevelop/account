package accountSrv

import (
	"context"
	"errors"

	"github.com/wethedevelop/account/cache"
	"github.com/wethedevelop/account/model"
	"github.com/wethedevelop/account/serializer"
	pb "github.com/wethedevelop/proto/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// 实现第一方账号注册和登陆服务
type AccountServer struct {
	pb.UnsafeAccountAuthServer
}

// 通用获取用户身份接口
func (s *AccountServer) authToken(context context.Context, in *pb.TokenRequest) (*model.User, error) {
	uid, err := cache.GetUserByToken(in.Token)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if uid == "" {
		return nil, status.Error(codes.PermissionDenied, serializer.ACCOUNT_INVALID_TOKEN)
	}
	user, err := model.GetUser(uid)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if uid == "" {
		return nil, status.Error(codes.PermissionDenied, serializer.ACCOUNT_INVALID_TOKEN)
	}
	return &user, nil
}

// Signup 注册接口
func (s *AccountServer) Signup(context context.Context, in *pb.SignupRequest) (*pb.User, error) {
	account := in.GetAccount()
	password := in.GetPassword()
	// 不允许为空
	if account == "" || password == "" {
		return nil, status.Error(codes.InvalidArgument, serializer.ACCOUNT_OR_PWD_EMPTY)
	}
	// 检测重名
	checked, err := model.CheckRegistered(account)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if checked {
		return nil, status.Error(codes.InvalidArgument, serializer.ACCOUNT_REGISTERED)
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
	return serializer.BuildUser(user), nil
}

// Signin 登录接口
func (s *AccountServer) Signin(context context.Context, in *pb.SigninRequest) (*pb.TokenResponse, error) {
	account := in.GetAccount()
	password := in.GetPassword()

	var user model.User
	if err := model.DB.Where("account = ?", account).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.PermissionDenied, serializer.ACCOUNT_OR_PWD_NOT_MATCH)
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	if !user.CheckPassword(password) {
		return nil, status.Error(codes.PermissionDenied, serializer.ACCOUNT_OR_PWD_NOT_MATCH)
	}

	token, tokenExpire, err := user.MakeToken()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.TokenResponse{Token: token, TokenExpire: tokenExpire}, nil
}

// 通过Token验证用户身份
func (s *AccountServer) Check(context context.Context, in *pb.TokenRequest) (*pb.User, error) {
	user, err := s.authToken(context, &pb.TokenRequest{
		Token: in.Token,
	})
	if err != nil {
		return nil, err
	}
	return serializer.BuildUser(*user), nil
}

// 更新用户资料
func (s *AccountServer) Update(context context.Context, in *pb.UpdateRequest) (*pb.User, error) {
	user, err := s.authToken(context, &pb.TokenRequest{
		Token: in.Token,
	})
	if err != nil {
		return nil, err
	}
	user.Nickname = in.Nickname
	if err := user.Save(); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return serializer.BuildUser(*user), nil
}

// 特权：通过用户ID获取用户资料
func (s *AccountServer) DevGetUser(context context.Context, in *pb.GetUserRequest) (*pb.User, error) {
	user, err := model.GetUser(in.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "NotFound")
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return serializer.BuildUser(user), nil
}

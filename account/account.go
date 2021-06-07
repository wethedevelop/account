package account

import (
	"context"
	"fmt"

	"github.com/wethedevelop/account/model"
	pb "github.com/wethedevelop/proto/auth"
)

// 实现第一方账号注册和登陆服务
type AccountServer struct {
	pb.UnsafeAccountAuthServer
}

// Signup 是我们的注册服务
func (s *AccountServer) Signup(ctx context.Context, in *pb.SignupRequest) (*pb.User, error) {
	account := in.GetAccount()
	password := in.GetPassword()
	if account == "" || password == "" {
		return nil, fmt.Errorf("account or password is empty")
	}
	user := model.User{
		Account: account,
	}
	if err := user.SetPassword(password); err != nil {
		return nil, err
	}
	err := user.Create()
	if err != nil {
		return nil, err
	}
	return &pb.User{Nickname: "新用户"}, nil
}

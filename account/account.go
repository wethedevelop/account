package account

import (
	"context"
	"log"

	"github.com/wethedevelop/account/model"
	pb "github.com/wethedevelop/proto/auth"
)

// 实现第一方账号注册和登陆服务
type AccountServer struct {
	pb.UnsafeAccountAuthServer
}

// Signup 是我们的注册服务
func (s *AccountServer) Signup(ctx context.Context, in *pb.SignupRequest) (*pb.User, error) {
	log.Printf("Received: %s - %s", in.GetAccount(), in.GetPassword())

	account := in.GetAccount()
	password := in.GetPassword()
	user := model.User{
		Account: account,
	}
	user.SetPassword(password)
	err := user.Create()
	if err != nil {
		return nil, err
	}
	return &pb.User{Nickname: "新用户"}, nil
}

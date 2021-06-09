package test

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"

	accountSrv "github.com/wethedevelop/account/account"
	"github.com/wethedevelop/account/conf"
	"github.com/wethedevelop/account/serializer"
	"github.com/wethedevelop/account/util"
	pb "github.com/wethedevelop/proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	// 连接数据库
	conf.Init()

	// 启动grpc服务
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterAccountAuthServer(s, &accountSrv.AccountServer{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

// 账号注册测试
func TestSignUp(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewAccountAuthClient(conn)
	account := util.RandStringRunes(10)
	password := util.RandStringRunes(10)
	// 测试空密码无法注册
	rsp, err := client.Signup(ctx, &pb.SignupRequest{Account: account, Password: ""})
	if err == nil {
		t.Fatalf("Signup Password empty should not success %v", rsp)
	} else {
		e, ok := status.FromError(err)
		if !ok || e.Code() != codes.InvalidArgument || e.Message() != serializer.ACCOUNT_OR_PW_EMPTY {
			t.Fatalf("Signup Password empty should not success: %v", err)
		}
	}
	// 正常注册
	rsp, err = client.Signup(ctx, &pb.SignupRequest{Account: account, Password: password})
	if err != nil {
		t.Fatalf("Signup failed: %v", err)
	}
	if rsp.Id == 0 {
		t.Fatalf("Signup failed: %v", rsp)
	}
	// 用户名被占用应该无法注册
	rsp, err = client.Signup(ctx, &pb.SignupRequest{Account: account, Password: password})
	if err == nil {
		t.Fatalf("ACCOUNT_REGISTERED should not success %v", rsp)
	} else {
		e, ok := status.FromError(err)
		if !ok || e.Code() != codes.InvalidArgument || e.Message() != serializer.ACCOUNT_REGISTERED {
			t.Fatalf("ACCOUNT_REGISTERED should not success: %v", rsp)
		}
	}
}

// 账号登录测试
func TestSignin(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewAccountAuthClient(conn)
	account := util.RandStringRunes(10)
	password := util.RandStringRunes(10)
	// 正常注册
	rsp, err := client.Signup(ctx, &pb.SignupRequest{Account: account, Password: password})
	if err != nil {
		t.Fatalf("Signup failed: %v", err)
	}
	if rsp.Id == 0 {
		t.Fatalf("Signup failed: %v", rsp)
	}
	// 正常登录
	token, err := client.Signin(ctx, &pb.SigninRequest{Account: account, Password: password})
	if err != nil {
		t.Fatalf("Signin failed: %v", err)
	}
	if token.Token == "" || token.TokenExpire == 0 {
		t.Fatalf("Signin failed: %v", token)
	}
	// 瞎传个密码
	token, err = client.Signin(ctx, &pb.SigninRequest{Account: account, Password: "123456"})
	if err == nil {
		t.Fatalf("ACCOUNT_OR_PWD_NOT_MATCH should not success %v", token)
	} else {
		e, ok := status.FromError(err)
		if !ok || e.Code() != codes.PermissionDenied || e.Message() != serializer.ACCOUNT_OR_PWD_NOT_MATCH {
			t.Fatalf("ACCOUNT_OR_PWD_NOT_MATCH should not success: %v", err)
		}
	}
}

// 获取当前用户身份接口
func TestCheckUser(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewAccountAuthClient(conn)
	account := util.RandStringRunes(10)
	password := util.RandStringRunes(10)
	// 正常注册
	rsp, err := client.Signup(ctx, &pb.SignupRequest{Account: account, Password: password})
	if err != nil {
		t.Fatalf("Signup failed: %v", err)
	}
	if rsp.Id == 0 {
		t.Fatalf("Signup failed: %v", rsp)
	}
	// 正常登录
	token, err := client.Signin(ctx, &pb.SigninRequest{Account: account, Password: password})
	if err != nil {
		t.Fatalf("Signin failed: %v", err)
	}
	if token.Token == "" || token.TokenExpire == 0 {
		t.Fatalf("Signin failed: %v", token)
	}
	user, err := client.Check(ctx, &pb.TokenRequest{Token: token.Token})
	if err != nil {
		t.Fatalf("User Check failed: %v", err)
	}
	if user.Id == 0 {
		t.Fatalf("User Check failed: %v", user)
	}
}

// 账号资料更新接口
func TestUpdateUser(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewAccountAuthClient(conn)
	account := util.RandStringRunes(10)
	password := util.RandStringRunes(10)
	// 正常注册
	rsp, err := client.Signup(ctx, &pb.SignupRequest{Account: account, Password: password})
	if err != nil {
		t.Fatalf("Signup failed: %v", err)
	}
	if rsp.Id == 0 {
		t.Fatalf("Signup failed: %v", rsp)
	}
	// 正常登录
	token, err := client.Signin(ctx, &pb.SigninRequest{Account: account, Password: password})
	if err != nil {
		t.Fatalf("Signin failed: %v", err)
	}
	if token.Token == "" || token.TokenExpire == 0 {
		t.Fatalf("Signin failed: %v", token)
	}
	name := util.RandStringRunes(10)
	// 修改昵称
	user, err := client.Update(ctx, &pb.UpdateRequest{Token: token.Token, Nickname: name})
	if err != nil {
		t.Fatalf("User Update failed: %v", err)
	}
	if user.Nickname != name {
		t.Fatalf("User Update failed: %v", user)
	}
}

// 账号资料更新接口
func TestGetUser(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewAccountAuthClient(conn)
	account := util.RandStringRunes(10)
	password := util.RandStringRunes(10)
	// 正常注册
	rsp, err := client.Signup(ctx, &pb.SignupRequest{Account: account, Password: password})
	if err != nil {
		t.Fatalf("Signup failed: %v", err)
	}
	if rsp.Id == 0 {
		t.Fatalf("Signup failed: %v", rsp)
	}
	// 修改昵称
	user, err := client.DevGetUser(ctx, &pb.GetUserRequest{Id: rsp.Id})
	if err != nil {
		t.Fatalf("DevGetUser failed: %v", err)
	}
	if user.Nickname != rsp.Nickname {
		t.Fatalf("DevGetUser failed: %v", user)
	}
}

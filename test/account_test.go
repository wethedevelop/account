package test

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/wethedevelop/account/account"
	"github.com/wethedevelop/account/conf"
	pb "github.com/wethedevelop/proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	// 连接数据库
	conf.Init()

	// 启动grpc服务
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterAccountAuthServer(s, &account.AccountServer{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestSignUp(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewAccountAuthClient(conn)
	resp, err := client.Signup(ctx, &pb.SignupRequest{Account: "chengka", Password: "12345678"})
	if err != nil {
		t.Fatalf("Signup failed: %v", err)
	}
	log.Printf("Response: %+v", resp)
}

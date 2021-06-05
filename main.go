package main

import (
	"log"
	"net"

	"github.com/wethedevelop/account/account"
	"github.com/wethedevelop/account/conf"
	pb "github.com/wethedevelop/proto/auth"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

func main() {
	conf.Init()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAccountAuthServer(s, &account.AccountServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

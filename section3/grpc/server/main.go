package main

import (
	"context"
	"goKitRPC/grpc/pb"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct{}

func (s *server) CheckPassword(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	if in.Username == "13sai" && in.Password == "123456" {
		return &pb.LoginResponse{Ret: "success"}, nil
	} else {
		return &pb.LoginResponse{Ret: "fail"}, nil
	}
}

func main() {
	lis, err := net.Listen("tcp", "127.0.0.1:9234")
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}

	gServer := grpc.NewServer()
	pb.RegisterUserServiceServer(gServer, &server{})

	if err := gServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve %v", err)
	}
}

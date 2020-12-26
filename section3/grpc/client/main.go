package main

import (
	"context"
	"fmt"
	"goKitRPC/grpc/pb"

	"google.golang.org/grpc"
)

func main() {
	serviceAddress := "127.0.0.1:9234"
	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure())
	if err != nil {
		panic("connect error")
	}

	defer conn.Close()

	userClient := pb.NewUserServiceClient(conn)
	userReq := &pb.LoginRequest{Username: "13sai", Password: "123456"}

	reply, err := userClient.CheckPassword(context.Background(), userReq)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(reply)
}

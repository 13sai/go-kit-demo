package main

import (
	"context"
	"fmt"
	"goKitRPC/endpoint"
	"goKitRPC/grpc/pb"
	"goKitRPC/service"
	"goKitRPC/transport"
	"log"
	"net"
	"time"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

var addr = "127.0.0.1:8009"

func main() {
	go serverRun()
	time.Sleep(3 * time.Second)
	clientTest()
}

func clientTest() {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("dial error")
	}

	defer conn.Close()
	client := pb.NewUserServiceClient(conn)
	userReq := &pb.LoginRequest{Username: "sai", Password: "11111111"}

	reply, err := client.CheckPassword(context.Background(), userReq)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(reply)
	time.Sleep(3 * time.Second)
}

func serverRun() {
	ctx := context.Background()

	var s service.UserService
	s = service.UserServiceImpl{}

	end := endpoint.MakeUserEndpoint(s)

	rateLimit := rate.NewLimiter(rate.Every(time.Second*1), 100)

	end = endpoint.NewTokenBucketLimit(rateLimit)(end)

	ends := endpoint.UEndpoint{
		UserEndpoint: end,
	}

	hand := transport.NewUserServer(ctx, ends)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("listen error")
	}

	gServer := grpc.NewServer()

	pb.RegisterUserServiceServer(gServer, hand)
	gServer.Serve(lis)
}

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/google/uuid"

	"go-kit-demo/discovery"
	"go-kit-demo/endpoint"
	"go-kit-demo/service"
	"go-kit-demo/transport"
)

func main() {
	consulAddr := flag.String("consul.addr", "localhost", "consul address")
	consulPort := flag.Int("consul.port", 8500, "consul port")
	serviceName := flag.String("service.name", "register", "service name")
	serviceAddr := flag.String("service.addr", "localhost", "service addr")
	servicePort := flag.Int("service.port", 12312, "service port")

	flag.Parse()

	// consul客户端
	client := discovery.NewDiscoveryClient(*consulAddr, *consulPort)

	// model.GetDB()
	// model.InitRedis()

	ctx := context.Background()
	errChan := make(chan error)

	srv := service.NewRegisterServiceImpl(client)

	endpoints := endpoint.RegisterEndpoints{
		DiscoveryEndpoint:   endpoint.MakeDiscoveryEndpoint(srv),
		HealthCheckEndpoint: endpoint.MakeHealthCheckEndpoint(srv),
	}
	r := transport.NewHttpHandle(ctx, &endpoints)

	go func() {
		errChan <- http.ListenAndServe(":"+strconv.Itoa(*servicePort), r)
	}()
	go func() {
		// 监控系统信号，等待 ctrl + c 系统信号通知服务关闭
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	instanceId := *serviceName + "-" + uuid.New().String()
	err := client.Register(ctx, *serviceName, instanceId, "/health", *serviceAddr, *servicePort, nil, nil)
	if err != nil {
		log.Printf("register service error: %s", err)
		os.Exit(-1)
	}

	err = <-errChan
	fmt.Println("over!!!!")
	fmt.Println(err)

	client.Deregister(ctx, instanceId)
}

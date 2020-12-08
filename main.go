package main

import (
	"fmt"
	"net/http"
	"log"
	"os"
	"os/signal"
	"syscall"
	"context"

	"github.com/spf13/viper"
	"github.com/spf13/pflag"

	"local.com/13sai/go-kit-demo/service"
	"local.com/13sai/go-kit-demo/endpoint"
	"local.com/13sai/go-kit-demo/transport"
	"local.com/13sai/go-kit-demo/config"
	"local.com/13sai/go-kit-demo/model"
)

var (
	conf = pflag.StringP("config", "c", "", "config filepath")
)

func main() {
	pflag.Parse()

	// 初始化配置
	if err := config.Init(*conf); err != nil {
		panic(err)
	}
	model.GetDB()
	model.InitRedis()

	ctx := context.Background()
	errChan := make(chan error)

	userService := service.MakeUserServiceImpl(&model.UserDAOImpl{})

	userEndpoints := &endpoint.UserEndpoints{ 
		endpoint.MakeRegisterEndpoint(userService), 
		endpoint.MakeLoginEndPoint(userService), 
	} 
	r := transport.NewHttpHandle(ctx, userEndpoints) 
	go func() {
		fmt.Println("server on " + viper.GetString("addr"))
		errChan <- http.ListenAndServe(viper.GetString("addr"), r) 
	}() 
	go func() { 
		// 监控系统信号，等待 ctrl + c 系统信号通知服务关闭 
		c := make(chan os.Signal, 1) 
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM) 
		errChan <- fmt.Errorf("%s", <-c) 
	}() 
	error := <-errChan 
	log.Println(error) 
}
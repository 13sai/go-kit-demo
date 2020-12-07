package main

import (
	"fmt"
	"net/http"

	"github.com/spf13/viper"
	"github.com/spf13/pflag"

	"local.com/13sai/go-kit-demo/service"
	"local.com/13sai/go-kit-demo/endpoint"
	"local.com/13sai/go-kit-demo/transport"
	"local.com/13sai/go-kit-demo/config"
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
	server := service.NewService()
	endpoints := endpoint.NewEndPointServer(server)
	httpHandle := transport.NewHttpHandle(endpoints)

	fmt.Println("server run ", viper.GetString("addr"))
	http.ListenAndServe(viper.GetString("addr"), httpHandle)
}
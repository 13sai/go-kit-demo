package main

import (
	"fmt"
	"net/http"

	"local.com/13sai/go-kit-demo/service"
	"local.com/13sai/go-kit-demo/endpoint"
	"local.com/13sai/go-kit-demo/transport"
)

func main() {
	server := service.NewService()
	endpoints := endpoint.NewEndPointServer(server)
	httpHandle := transport.NewHttpHandle(endpoints)

	fmt.Println("server run :9008")
	http.ListenAndServe("0.0.0.0:9008", httpHandle)
}
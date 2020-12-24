package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

type StringReq struct {
	A string
	B string
}

type Service interface {
	ConCat(req StringReq, ret *string) error
}

type StringService struct{}

func (s StringService) Concat(req StringReq, ret *string) error {
	if len(req.A)+len(req.B) > 10 {
		*ret = ""
		return errors.New("max")
	}
	*ret = req.A + req.B
	return nil
}

func main() {
	s := new(StringService)
	rpc.Register(s)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", "127.0.0.1:8989")
	if e != nil {
		log.Fatal("listen error", e)
	}
	go client()
	http.Serve(l, nil)
}

func client() {
	fmt.Println('a')
	time.Sleep(3 * time.Second)
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:8989")
	if err != nil {
		log.Fatal("dial error", err)
	}

	sReq := &StringReq{"aa", "c"}
	var res string
	err = client.Call("StringService.Concat", sReq, &res)

	if err != nil {
		log.Fatal("error", err)
	}
	fmt.Println(res)
}

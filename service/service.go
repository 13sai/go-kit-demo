package service

import (
	"context"
	"errors"
	"fmt"

	"local.com/13sai/go-kit-demo/util"
)

const ContextReqUUid = "uuid"

type Service interface {
	Add(ctx context.Context, in Add) AddAck
	Login(ctx context.Context, in Login) (ack LoginAck, err error)
}

type baseServer struct {}

func NewService() Service {
	return &baseServer{}
}

func (s baseServer) Add(ctx context.Context, in Add) AddAck {
	return AddAck{in.A+in.B}
}

func (s baseServer) Login(ctx context.Context, in Login) (ack LoginAck, err error) {
	if in.Name != "sai" || in.Pass != "123456" {
		err = errors.New("not matched")
	}
	fmt.Println(in)

	ack.Token, _ = util.JWTCreate(in.Name, 13)
	return 
}
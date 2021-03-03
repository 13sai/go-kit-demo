package service

import (
	"context"
	"errors"
)

// service 单纯处理业务逻辑

type UserService interface {
	CheckPassword(ctx context.Context, username string, password string) (bool, error)
}

type UserServiceImpl struct{}

func (s UserServiceImpl) CheckPassword(ctx context.Context, username string, password string) (bool, error) {
	if username == "sai" && password == "111111" {
		return true, nil
	}

	return false, errors.New("deny")
}

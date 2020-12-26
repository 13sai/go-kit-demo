package endpoint

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"

	"go-kit-demo/service"
)

type UserEndpoints struct {
	RegisterEndpoint endpoint.Endpoint
	LoginEndPoint    endpoint.Endpoint
}

type LoginRequest struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type LoginResponse struct {
	UserInfo *service.UserInfoDTO `json:"userinfo"`
}

func MakeLoginEndPoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (res interface{}, err error) {
		req := request.(*LoginRequest)
		fmt.Println(req)
		userInfo, err := s.Login(ctx, req.Email, req.Password)
		return &LoginResponse{UserInfo: userInfo}, err
	}
}

type RegisterRequest struct {
	Username string `json:"username" form:"username"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type RegisterResponse struct {
	UserInfo *service.UserInfoDTO `json:"userinfo"`
}

func MakeRegisterEndpoint(userService service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(RegisterRequest)
		userInfo, err := userService.Register(ctx, &service.RegisterUserVO{
			Username: req.Username,
			Password: req.Password,
			Email:    req.Email,
		})
		return &RegisterResponse{UserInfo: userInfo}, err

	}
}

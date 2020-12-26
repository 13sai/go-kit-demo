package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"goKitRPC/service"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRes struct {
	Ret bool  `json:"ret"`
	Err error `json:"err"`
}
type UEndpoint struct {
	UserEndpoint endpoint.Endpoint
}

func MakeUserEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, form interface{}) (result interface{}, err error) {
		req := form.(LoginRequest)
		ret, err := s.CheckPassword(ctx, req.Username, req.Password)
		return LoginRes{ret, err}, nil
	}
}

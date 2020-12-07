package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"local.com/13sai/go-kit-demo/service"
)

type EndPointServer struct {
	AddEndPoint endpoint.Endpoint
	LoginEndPoint endpoint.Endpoint
}

func NewEndPointServer(s service.Service) EndPointServer {
	var addEndPoint endpoint.Endpoint
	{
		addEndPoint = MakeAddEndPoint(s)
		addEndPoint = AuthMiddleware()(addEndPoint)
	}

	var loginEndPoint endpoint.Endpoint
	{
		loginEndPoint = MakeLoginEndPoint(s)
		// loginEndPoint = AuthMiddleware()(loginEndPoint)
	}

	return EndPointServer{addEndPoint, loginEndPoint}
}

func (s EndPointServer) Add(ctx context.Context, in service.Add) service.AddAck {
	res, _ := s.AddEndPoint(ctx, in)
	return res.(service.AddAck)
}

func MakeAddEndPoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (res interface{}, err error) {
		req := request.(service.Add)
		res = s.Add(ctx, req)
		return
	}
}

func MakeLoginEndPoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (res interface{}, err error) {
		req := request.(service.Login)
		res, _ = s.Login(ctx, req)
		return
	}
}
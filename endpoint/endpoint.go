package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"local.com/13sai/go-kit-demo/service"
)

type EndPointServer struct {
	AddEndPoint endpoint.Endpoint
}

func NewEndPointServer(s service.Service) EndPointServer {
	var addEndPoint endpoint.Endpoint

	addEndPoint = MakeAddEndPoint(s)

	return EndPointServer{addEndPoint}
}

func (s EndPointServer) Add(ctx context.Context, in service.Add) service.AddAck {
	res, _ := s.AddEndPoint(ctx, in)
	return res.(service.AddAck)
}

func MakeAddEndPoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (res interface{}, err error) {
		req := request.(service.Add)
		res = s.TestAdd(ctx, req)
		return
	}
}
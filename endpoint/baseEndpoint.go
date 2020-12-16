package endpoint

/**
 抽象层

 两个路由：
- health 健康检查
- Discovery 服务发现
*/

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"go-kit-demo/discovery"
	"go-kit-demo/service"
)

type RegisterEndpoints struct {
	DiscoveryEndpoint   endpoint.Endpoint
	HealthCheckEndpoint endpoint.Endpoint
}

// 服务发现请求结构体
type DiscoveryRequest struct {
	ServiceName string
}

// 服务发现响应结构体
type DiscoveryResponse struct {
	Instances []*discovery.InstanceInfo `json:"instances"`
	Error     string                    `json:"error"`
}

func MakeDiscoveryEndpoint(s service.DiscoveryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (res interface{}, err error) {
		req := request.(DiscoveryRequest)
		instances, err := s.DiscoveryService(ctx, req.ServiceName)
		var errString = ""

		if err != nil {
			errString = err.Error()
		}
		return &DiscoveryResponse{
			Instances: instances,
			Error:     errString,
		}, nil
	}
}

// HealthRequest 健康检查请求结构
type HealthRequest struct{}

// HealthResponse 健康检查响应结构
type HealthResponse struct {
	Status string `json:"status"`
}

// MakeHealthCheckEndpoint 创建健康检查Endpoint
func MakeHealthCheckEndpoint(svc service.DiscoveryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		status := svc.HealthCheck()
		return HealthResponse{
			Status: status,
		}, nil
	}
}

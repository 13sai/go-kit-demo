package service

/**
service，业务逻辑层

提供两个接口，健康检查和服务发现
*/
import (
	"context"
	"errors"
	"log"

	"go-kit-demo/discovery"
)

type DiscoveryService interface {
	// 健康检查接口
	HealthCheck() string
	// 服务发现接口
	DiscoveryService(ctx context.Context, serviceName string) ([]*discovery.InstanceInfo, error)
}

var ErrNotServiceInstances = errors.New("instances are not existed")

type RegisterServiceImpl struct {
	discoveryClient *discovery.DiscoveryClient
}

func NewRegisterServiceImpl(discoveryClient *discovery.DiscoveryClient) DiscoveryService {
	return &RegisterServiceImpl{
		discoveryClient: discoveryClient,
	}
}

func (service *RegisterServiceImpl) DiscoveryService(ctx context.Context, serviceName string) ([]*discovery.InstanceInfo, error) {
	instances, err := service.discoveryClient.DiscoverServices(ctx, serviceName)

	if err != nil {
		log.Printf("get service info err: %s", err)
	}
	if instances == nil || len(instances) == 0 {
		return nil, ErrNotServiceInstances
	}

	return instances, nil
}

func (*RegisterServiceImpl) HealthCheck() string {
	return "OK"
}

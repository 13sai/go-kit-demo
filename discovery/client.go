package discovery

/**
  consul客户端
  字段较为复杂，但逻辑简单，提供3个方法，注册、注销、发现，都是通过http请求（consul）
*/

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type InstanceInfo struct {
	ID                string            `json:"ID"`                // 服务实例ID
	Service           string            `json:"Service,omitempty"` // 服务发现时返回的服务名
	Name              string            `json:"Name"`              // 服务名
	Tags              []string          `json:"Tags,omitempty"`    // 标签，可用于进行服务过滤
	Address           string            `json:"Address"`           // 服务实例HOST
	Port              int               `json:"Port"`              // 服务实例端口
	Meta              map[string]string `json:"Meta,omitempty"`    // 元数据
	EnableTagOverride bool              `json:"EnableTagOverride"` // 是否允许标签覆盖
	Check             Check             `json:"Check,omitempty"`   // 健康检查相关配置
	Weights           Weights           `json:"Weights,omitempty"` // 权重
}

type Check struct {
	DeregisterCriticalServiceAfter string   `json:"DeregisterCriticalServiceAfter"` // 多久之后注销服务
	Args                           []string `json:"Args,omitempty"`                 // 请求参数
	HTTP                           string   `json:"HTTP"`                           // 健康检查地址
	Interval                       string   `json:"Interval,omitempty"`             // Consul 主动检查间隔
	TTL                            string   `json:"TTL,omitempty"`                  // 服务实例主动维持心跳间隔，与Interval只存其一
}

type Weights struct {
	Passing int `json:"Passing"`
	Warning int `json:"Warning"`
}

type DiscoveryClient struct {
	host string // Consul 的 Host
	port int    // Consul 的 端口
}

func NewDiscoveryClient(host string, port int) *DiscoveryClient {
	return &DiscoveryClient{
		host: host,
		port: port,
	}
}

// service register
func (consulClient *DiscoveryClient) Register(ctx context.Context, serviceName, instanceId, healthCheckUrl string, instanceHost string, instancePort int, meta map[string]string, weights *Weights) error {
	instanceInfo := &InstanceInfo{
		ID:                instanceId,
		Name:              serviceName,
		Address:           instanceHost,
		Port:              instancePort,
		Meta:              meta,
		EnableTagOverride: false,
		Check: Check{
			DeregisterCriticalServiceAfter: "20s",
			HTTP:                           "http://" + instanceHost + ":" + strconv.Itoa(instancePort) + healthCheckUrl,
			Interval:                       "15s",
		},
	}

	if weights != nil {
		instanceInfo.Weights = *weights
	} else {
		instanceInfo.Weights = Weights{
			Passing: 10,
			Warning: 1,
		}
	}

	byteData, err := json.Marshal(instanceInfo)
	if err != nil {
		log.Printf("json format error:%s", err)
		return err
	}

	req, err := http.NewRequest("PUT", "http://"+consulClient.host+":"+strconv.Itoa(consulClient.port)+"/v1/agent/service/register", bytes.NewReader(byteData))
	if err != nil {
		log.Printf("http error:%s", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	client := http.Client{}
	client.Timeout = 2 * time.Second
	res, err := client.Do(req)
	if err != nil {
		log.Printf("http res error:%s", err)
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Printf("register service http request errCode : %v", res.StatusCode)
		return fmt.Errorf("register service http request errCode : %v", res.StatusCode)
	}

	fmt.Println("register service success")
	return nil
}

// service deregister
func (consulClient *DiscoveryClient) Deregister(ctx context.Context, instanceId string) error {
	req, err := http.NewRequest("PUT",
		"http://"+consulClient.host+":"+strconv.Itoa(consulClient.port)+"/v1/agent/service/deregister/"+instanceId, nil)

	if err != nil {
		log.Printf("req format err: %s", err)
		return err
	}

	client := http.Client{}
	client.Timeout = time.Second * 2

	res, err := client.Do(req)

	if err != nil {
		log.Printf("deregister service err : %s", err)
		return err
	}

	res.Body.Close()

	if res.StatusCode != 200 {
		log.Printf("deresigister service http request errCode : %v", res.StatusCode)
		return fmt.Errorf("deresigister service http request errCode : %v", res.StatusCode)
	}

	fmt.Println("deregister service success")
	return nil
}

// discover service
func (consulClient *DiscoveryClient) DiscoverServices(ctx context.Context, serviceName string) ([]*InstanceInfo, error) {
	req, err := http.NewRequest("GET",
		"http://"+consulClient.host+":"+strconv.Itoa(consulClient.port)+"/v1/health/service/"+serviceName, nil)

	if err != nil {
		log.Printf("req format err: %s", err)
		return nil, err
	}

	client := http.Client{}
	client.Timeout = time.Second * 2

	res, err := client.Do(req)
	if err != nil {
		log.Printf("discover service err : %s", err)
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Printf("discover service http request errCode : %v", res.StatusCode)
		return nil, fmt.Errorf("discover service http request errCode : %v", res.StatusCode)
	}

	var serviceList []struct {
		Service InstanceInfo `json:"Service"`
	}
	err = json.NewDecoder(res.Body).Decode(&serviceList)
	if err != nil {
		log.Printf("format service info err : %s", err)
		return nil, err
	}

	instances := make([]*InstanceInfo, len(serviceList))
	for i := 0; i < len(instances); i++ {
		instances[i] = &serviceList[i].Service
	}
	return instances, nil

}

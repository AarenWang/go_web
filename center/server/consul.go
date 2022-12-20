package server

import (
	"center/pb"
	"fmt"
	"github.com/hashicorp/consul/api"
)

var serviceName = pb.ServiceType_name[int32(pb.ServiceType_CENTER)]

var ConsulClient *api.Client

func InitConsulClient(id string, myIp string, myport int, consulIp string, consulPort int) {
	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%d", consulIp, consulPort)
	client, err := api.NewClient(config)
	if err != nil {
		panic("server client connect failed")
	}
	checkIp := myIp
	if consulIp == "127.0.0.1" {
		checkIp = "127.0.0.1"
	}
	// 创建注册到consul的服务到
	registration := new(api.AgentServiceRegistration)
	registration.ID = id            // 服务节点的名称
	registration.Name = serviceName // 服务名称
	registration.Port = myport      // 服务端口
	registration.Address = myIp     // 服务 IP 要确保consul可以访问这个ip
	// 健康检查
	registration.Check = &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d/Health", checkIp, myport),
		Timeout:                        "1s",
		Interval:                       "15s",
		DeregisterCriticalServiceAfter: "30s",
	}
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic("service register failed")
	}
	ConsulClient = client
}

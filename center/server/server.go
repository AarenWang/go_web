package server

import (
	"center/pb"
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
	"strings"
)

type server struct {
	pb.UnimplementedCenterServiceServer
}

func (s *server) Register(c context.Context, service *pb.Service) (*emptypb.Empty, error) {
	registration := new(api.AgentServiceRegistration)
	id := service.GetId()
	port := int(service.GetPort())
	addr := service.GetAddr()
	if addr == "" {
		addr = getPeerAddr(c)
	}
	serviceName := pb.ServiceType_name[int32(service.GetType())]
	check := &api.AgentServiceCheck{
		Timeout:                        "1s",
		Interval:                       "15s",
		DeregisterCriticalServiceAfter: "30s",
	}
	scheme := pb.Scheme_name[int32(service.GetScheme())]
	switch service.GetScheme() {
	case pb.Scheme_HTTP, pb.Scheme_HTTPS:
		check.HTTP = fmt.Sprintf("%s://%s:%d/health", scheme, addr, port)
	case pb.Scheme_GRPC:
		check.GRPC = fmt.Sprintf("%s:%d/Health", addr, port)
	}
	meta := make(map[string]string)
	meta["scheme"] = scheme
	registration.ID = id
	registration.Name = serviceName
	registration.Port = port
	registration.Address = addr
	registration.Check = check
	registration.Meta = meta
	err := ConsulClient.Agent().ServiceRegister(registration)
	return &emptypb.Empty{}, err
}
func (s *server) Logout(c context.Context, req *pb.LogoutReq) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, ConsulClient.Agent().ServiceDeregister(req.GetServiceId())
}

func (s *server) GetServices(c context.Context, req *pb.ServicesReq) (*pb.Services, error) {
	serviceType := req.GetType()
	serviceName = pb.ServiceType_name[int32(serviceType)]
	serviceMap, err := ConsulClient.Agent().ServicesWithFilter(fmt.Sprintf("Service==%s", serviceName))
	if err != nil {
		return nil, err
	}
	list := make([]*pb.Service, len(serviceMap))
	i := 0
	for _, value := range serviceMap {
		var scheme pb.Scheme
		if value.Meta != nil {
			t := value.Meta["scheme"]
			if t != "" {
				t = strings.ToUpper(t)
				scheme = pb.Scheme(pb.Scheme_value[t])
			}
		}
		list[i] = &pb.Service{
			Id:     value.ID,
			Type:   serviceType,
			Addr:   value.Address,
			Port:   int32(value.Port),
			Scheme: scheme,
		}
		i++
	}
	return &pb.Services{List: list}, nil
}

func NewServer() *server {
	return &server{}
}

// GetRealAddr get real client ip
func getRealAddr(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	rips := md.Get("x-real-ip")
	if len(rips) == 0 {
		return ""
	}

	return rips[0]
}

// GetPeerAddr get peer addr
func getPeerAddr(ctx context.Context) string {
	var addr string
	if pr, ok := peer.FromContext(ctx); ok {
		if tcpAddr, ok := pr.Addr.(*net.TCPAddr); ok {
			addr = tcpAddr.IP.String()
		} else {
			addr = pr.Addr.String()
		}
	}
	return addr
}

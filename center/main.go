package main

import (
	"center/pb"
	"center/server"
	"common"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
)

func main() {
	consulAddr := flag.String("consul-addr", "127.0.0.1", "server agent ip address")
	consulPort := flag.Int("consul-port", 8500, "server agent pb port")
	id := flag.String("id", "none", "node id")
	ip := flag.String("ip", "127.0.0.1", "my ip address")
	port := flag.Int("port", 18888, "my ip address")
	flag.Parse()

	server.InitConsulClient(*id, *ip, *port, *consulAddr, *consulPort)

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		panic("listen start fail")
	}
	s := grpc.NewServer()
	pb.RegisterCenterServiceServer(s, server.NewServer())
	grpc_health_v1.RegisterHealthServer(s, &common.HealthService{})
	if err := s.Serve(listen); err != nil {
		panic("server start fail")
	}
}

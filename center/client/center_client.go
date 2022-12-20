package client

import (
	"center/pb"
	"context"
	"errors"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var logger *zap.Logger

type CenterServiceClientHelper struct {
	conn      *grpc.ClientConn
	service   pb.CenterServiceClient
	ServiceDb *serviceDb
	my        *Service
}

func (c *CenterServiceClientHelper) Connect() error {
	list := c.ServiceDb.List
	if len(list) == 0 {
		return errors.New("no center address available")
	}
	for i := 0; i < len(list); i++ {
		service := list[i]
		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", service.Addr, service.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err == nil {
			c.conn = conn
			c.service = pb.NewCenterServiceClient(conn)
			return nil
		}
	}
	logger.Error("connect fail")
	return errors.New("center service connect fail")
}

func (c *CenterServiceClientHelper) Close() error {
	return c.conn.Close()
}

func (c *CenterServiceClientHelper) Register(ctx context.Context) error {
	_, err := c.service.Register(ctx, (*pb.Service)(c.my))
	return err
}

func (c *CenterServiceClientHelper) Logout(ctx context.Context) error {
	_, err := c.service.Logout(ctx, &pb.LogoutReq{ServiceId: c.my.Id})
	return err
}

func (c *CenterServiceClientHelper) GetServices(ctx context.Context, serviceType pb.ServiceType) ([]*Service, error) {
	req, err := c.service.GetServices(ctx, &pb.ServicesReq{Type: serviceType})
	if err != nil {
		return nil, err
	}
	reqList := req.GetList()
	list := make([]*Service, len(reqList))
	for idx, item := range reqList {
		list[idx] = (*Service)(item)
	}
	return list, err
}

func (c *CenterServiceClientHelper) SyncCenterServices(ctx context.Context) error {
	list, err := c.GetServices(ctx, pb.ServiceType_CENTER)
	if err != nil {
		return err
	}
	c.ServiceDb.List = list
	c.ServiceDb.Save()
	return nil
}

func (c *CenterServiceClientHelper) Init(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if c.conn != nil {
		c.Close()
		c.Logout(ctx)
	}
	err := c.Connect()
	if err != nil {
		logger.Warn("connect fail")
		return err
	}
	err = c.Register(ctx)
	if err != nil {
		logger.Warn("register fail")
		return err
	}
	c.SyncCenterServices(ctx)
	return nil
}

func NewCenterServiceClientHelper(my *Service) CenterServiceClientHelper {
	//读取参数
	centerAddr := flag.String("center-addr", "", "center service addr")
	centerPort := flag.Int("center-port", 18888, "center service port")
	flag.Parse()

	serviceDb := NewServiceDb(pb.ServiceType_CENTER)
	if *centerAddr != "" {
		serviceDb.List = append(serviceDb.List, &Service{Scheme: pb.Scheme_GRPC, Type: pb.ServiceType_CENTER, Addr: *centerAddr, Port: int32(*centerPort)})
	}
	helper := CenterServiceClientHelper{
		my:        my,
		ServiceDb: serviceDb,
	}
	return helper
}

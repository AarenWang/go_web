package main

import (
	cc "center/client"
	cp "center/pb"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

var centerHelper cc.CenterServiceClientHelper

var serverPort = 80

func main() {
	my := cc.Service{
		Type:   cp.ServiceType_DOWNLOAD_CLIENT,
		Id:     "test123",
		Port:   int32(serverPort),
		Scheme: cp.Scheme_HTTP,
	}
	centerHelper = cc.NewCenterServiceClientHelper(&my)
	defer centerHelper.Close()
	err := centerHelper.Init(nil)
	if err != nil {
		panic(err)
	}
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	//conn, err := pb.Dial("127.0.0.1:8099", pb.WithTransportCredentials(insecure.NewCredentials()))
	//if err != nil {
	//	log.Fatalf("上游连接失败")
	//}
	//defer conn.Close()
	//downloadServiceClient = common.NewDownloadServiceClient(conn)
	//ctx, _ := context.WithTimeout(context.Background(), time.Second)
	//aa, err := downloadServiceClient.GetAddress(ctx, &common.DownloadAddressQuery{})
	//
	//if err != nil {
	//	log.Fatalf("请求失败")
	//}
	r.Run(fmt.Sprintf(":%d", serverPort))
}

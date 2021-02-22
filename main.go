package main

import (
	"iaas-api-server/common"
	"iaas-api-server/proto/image"
	"iaas-api-server/proto/natgateway"
	"iaas-api-server/proto/securitygroup"
	"iaas-api-server/proto/tenant"
	"iaas-api-server/service/imagesvc"
	"iaas-api-server/service/natgatewaysvc"
	"os"

	//"iaas-api-server/service/routesvc"
	"net"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	//	"iaas-api-server/service/imagesvc"
	//	"iaas-api-server/service/instancesvc"
	//	"iaas-api-server/service/nasdisksvc"
	//	"iaas-api-server/service/osssvc"
	//	"iaas-api-server/service/routesvc"
	"iaas-api-server/service/securitygroupsvc"
	//	"iaas-api-server/service/tenantsvc"
	"iaas-api-server/service/tenantsvc"
	//	"iaas-api-server/service/vpcsvc"
	//	"iaas-api-server/service/routesvc"
	//"iaas-api-server/proto/peerlink"
	//"iaas-api-server/service/peerlinksvc"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{}) //设置日志的输出格式为json格式，还可以设置为text格式
	log.SetOutput(os.Stdout)               //设置日志的输出为标准输出
	log.SetLevel(log.InfoLevel)            //设置日志的显示级别，这一级别以及更高级别的日志信息将会输出
	log.SetReportCaller(true)
	if !common.InitDb() {                  //数据库初始化
		panic("数据库初始化失败")
	}
}

func main() {
	rpcServer := grpc.NewServer()
	//注册服务
	//	flavor.RegisterFlavorServiceServer(rpcServer, &flavorsvc.FlavorService{})
	image.RegisterImageServiceServer(rpcServer, &imagesvc.ImageService{})
	//	nasdisk.RegisterNasDiskServiceServer(rpcServer, &nasdisksvc.NasDiskService{})
	//	oss.RegisterOSSServiceServer(rpcServer, &osssvc.OssService{})
	securitygroup.RegisterSecurityGroupServiceServer(rpcServer, &securitygroupsvc.SecurityGroupService{})
	tenant.RegisterTenantServiceServer(rpcServer, &tenantsvc.TenantService{})
	//	vpc.RegisterVpcServiceServer(rpcServer, &vpcsvc.VpcService{})
	//	route.RegisterRouteServiceServer(rpcServer, &routesvc.RouteService{})
	natgateway.RegisterNatGatewayServiceServer(rpcServer, &natgatewaysvc.NatGatewayService{})
	//peerlink.RegisterPeerLinkServiceServer(rpcServer, &peerlinksvc.PeerLinkService{})

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("服务监听端口失败", err)
	}

	err = rpcServer.Serve(listener)
	if err != nil {
		log.Fatal("服务启动失败", err)
	}
}

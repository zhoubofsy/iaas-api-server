package main

import (
	"iaas-api-server/proto/image"
	"iaas-api-server/proto/tenant"
	"iaas-api-server/service/imagesvc"
	"iaas-api-server/service/tenantsvc"
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
)

	log.SetFormatter(&log.JSONFormatter{}) //设置日志的输出格式为json格式，还可以设置为text格式
	log.SetOutput(os.Stdout)               //设置日志的输出为标准输出
	log.SetReportCaller(true)
}
func main() {


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

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("服务监听端口失败", err)
	}

	err = rpcServer.Serve(listener)
	if err != nil {
		log.Fatal("服务启动失败", err)
	}
}

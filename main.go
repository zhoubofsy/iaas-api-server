package main

import (
	"flag"
	"iaas-api-server/common/config"
	"iaas-api-server/proto/clouddisk"
	"iaas-api-server/proto/flavor"
	"iaas-api-server/proto/floatip"
	"iaas-api-server/proto/image"
	"iaas-api-server/proto/instance"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	//"iaas-api-server/proto/nasdisk"
	"iaas-api-server/proto/natgateway"
	"iaas-api-server/proto/oss"
	"iaas-api-server/proto/peerlink"
	"iaas-api-server/proto/route"
	"iaas-api-server/proto/securitygroup"
	"iaas-api-server/proto/tenant"
	"iaas-api-server/proto/vpc"
	"iaas-api-server/service/clouddisksvc"
	"iaas-api-server/service/flavorsvc"
	"iaas-api-server/service/floatipsvc"
	"iaas-api-server/service/imagesvc"
	"iaas-api-server/service/instancesvc"

	//"iaas-api-server/service/nasdisksvc"
	"iaas-api-server/service/natgatewaysvc"
	"iaas-api-server/service/osssvc"
	"iaas-api-server/service/peerlinksvc"
	"iaas-api-server/service/routesvc"
	"iaas-api-server/service/securitygroupsvc"
	"iaas-api-server/service/tenantsvc"
	"iaas-api-server/service/vpcsvc"
	"net"
	_ "net/http/pprof"
	"os"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{}) //设置日志的输出格式为json格式，还可以设置为text格式
	log.SetOutput(os.Stdout)               //设置日志的输出为标准输出
	log.SetLevel(log.InfoLevel)            //设置日志的显示级别，这一级别以及更高级别的日志信息将会输出
	log.SetReportCaller(true)              //设置日志的调用文件，调用函数
	log.SetFormatter(&log.JSONFormatter{}) //设置日志格式
}

var (
	conf = flag.String("conf", "", "config file")
)

func main() {
	flag.Parse()
	if *conf != "" {
		config.InitConfig(*conf)
	} else {
		panic("no config file.. usage:\n\t./serv.exe -conf xx.conf")
	}

	//go func() {
	//	log.Println(http.ListenAndServe(":10001", nil))
	//}()
	rpcServer := grpc.NewServer()
	//注册服务
	clouddisk.RegisterCloudDiskServiceServer(rpcServer, &clouddisksvc.CloudDiskService{})
	flavor.RegisterFlavorServiceServer(rpcServer, &flavorsvc.FlavorService{})
	instance.RegisterInstanceServiceServer(rpcServer, &instancesvc.InstanceService{})
	image.RegisterImageServiceServer(rpcServer, &imagesvc.ImageService{})
	//nasdisk.RegisterNasDiskServiceServer(rpcServer, &nasdisksvc.NasDiskService{})
	oss.RegisterOSSServiceServer(rpcServer, &osssvc.OSSService{})
	securitygroup.RegisterSecurityGroupServiceServer(rpcServer, &securitygroupsvc.SecurityGroupService{})
	tenant.RegisterTenantServiceServer(rpcServer, &tenantsvc.TenantService{})
	vpc.RegisterVpcServiceServer(rpcServer, &vpcsvc.VpcService{})
	route.RegisterRouterServiceServer(rpcServer, &routesvc.RouteService{})
	natgateway.RegisterNatGatewayServiceServer(rpcServer, &natgatewaysvc.NatGatewayService{})
	peerlink.RegisterPeerLinkServiceServer(rpcServer, &peerlinksvc.PeerLinkService{})
	floatip.RegisterFloatIpServiceServer(rpcServer, &floatipsvc.FloatIpService{})

	addr, _ := config.GetString("listening_addr")
	if addr == "" {
		addr = ":8080"
	}
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("服务监听端口失败", err)
	}

	err = rpcServer.Serve(listener)
	if err != nil {
		log.Fatal("服务启动失败", err)
	}
}

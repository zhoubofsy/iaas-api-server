package main

import (
	"google.golang.org/grpc"
	"iaas-api-server/proto/cloud_disk"
	"iaas-api-server/proto/flavor"
	"iaas-api-server/proto/image"
	"iaas-api-server/proto/instance"
	"iaas-api-server/proto/nas_disk"
	"iaas-api-server/proto/oss"
	"os"

	//"iaas-api-server/proto/route"
	"iaas-api-server/proto/security_group"
	"iaas-api-server/proto/tenant"
	"iaas-api-server/proto/vpc"
	"iaas-api-server/service/cloud_disk_service"
	"iaas-api-server/service/flavor_service"
	"iaas-api-server/service/image_service"
	"iaas-api-server/service/instance_service"
	"iaas-api-server/service/nas_disk_service"
	"iaas-api-server/service/oss_service"
	//"iaas-api-server/service/route_service"
	"iaas-api-server/service/security_group_service"
	"iaas-api-server/service/tenant_service"
	"iaas-api-server/service/vpc_service"
	log "github.com/sirupsen/logrus"
	"net"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})//设置日志的输出格式为json格式，还可以设置为text格式
	log.SetOutput(os.Stdout)//设置日志的输出为标准输出
	log.SetLevel(log.InfoLevel)//设置日志的显示级别，这一级别以及更高级别的日志信息将会输出
}

func main() {

	rpcServer := grpc.NewServer()

    //注册服务
	cloud_disk.RegisterCloudDiskServiceServer(rpcServer, &cloud_disk_service.CloudDiskService{})
	flavor.RegisterFlavorServiceServer(rpcServer, &flavor_service.FlavorService{})
	image.RegisterImageServiceServer(rpcServer, &image_service.ImageService{})
	instance.RegisterInstanceServiceServer(rpcServer, &instance_service.InstanceService{})
	nas_disk.RegisterNasDiskServiceServer(rpcServer, &nas_disk_service.NasDiskService{})
	oss.RegisterOSSServiceServer(rpcServer, &oss_service.OssService{})
	security_group.RegisterSecurityGroupServiceServer(rpcServer, &security_group_service.SecurityGroupService{})
	tenant.RegisterTenantServiceServer(rpcServer, &tenant_service.TenantService{})
	vpc.RegisterVpcServiceServer(rpcServer, &vpc_service.VpcService{})

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("服务监听端口失败", err)
	}

	err = rpcServer.Serve(listener)
	if err != nil {
		log.Fatal("服务启动失败", err)
	}
}

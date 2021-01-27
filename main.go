/*================================================================
*
*  文件名称：main.go
*  创 建 者: mongia
*  创建日期：2021年01月27日
*
================================================================*/

package main

import (
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	pb "iaas-api-server/proto/security_group"
	sgs "iaas-api-server/service/security_group_service"
	"net"
	"os"
)

const (
	port = ":50051"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to listen.")
		return
	}

	log.Info("listen port success.")

	s := grpc.NewServer()
	pb.RegisterSecurityGroupServiceServer(s, &sgs.SecurityGroupService{})

	log.Info("register grpc over.")
	if err := s.Serve(lis); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to serve.")
		return
	}
}

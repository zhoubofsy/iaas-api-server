/*================================================================
*
*  文件名称：main.go
*  创 建 者: mongia
*  创建日期：2021年02月02日
*
================================================================*/

package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"time"

	pb "iaas-api-server/proto/securitygroup"

	"google.golang.org/grpc"
)

var (
	address = flag.String("address", "localhost:8080", "grpc server address, localhost:8080")
	param   = flag.String("param", "", "grpc request param")
	method  = flag.String("method", "", "call method")
)

func main() {
	flag.Parse()

	// Set up a connection to the server.
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSecurityGroupServiceClient(conn)

	// Contact the server and print out its response.
	if "" == *param {
		log.Fatalf("input request param is null")
	}

	var data []byte = []byte(*param)
	//req := &pb.DeleteSecurityGroupReq{}
	req := &pb.CreateSecurityGroupReq{}

	err = json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("request is : %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	startTime := time.Now()
	defer func() {
		duration := time.Now().Sub(startTime)
		log.Printf("rpc waste time: %+v", duration)
	}()
	//r, err := c.DeleteSecurityGroup(ctx, req)
	r, err := c.CreateSecurityGroup(ctx, req)
	if err != nil {
		log.Fatalf("could not rpc request: %v", err)
	}
	log.Printf("rpc result: %+v", r)
}

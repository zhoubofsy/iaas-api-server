/*================================================================
*
*  文件名称：nat_gateway_cli.go
*  创 建 者: mongia
*  创建日期：2021年03月04日
*
================================================================*/

package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"google.golang.org/grpc"

	ngpb "iaas-api-server/proto/natgateway"
)

func init() {
	registerFunc("ngDelete", ngDelete)
	registerFunc("ngGet", ngGet)
	registerFunc("ngCreate", ngCreate)
}

func ngDelete(conn grpc.ClientConnInterface, data []byte) {
	c := ngpb.NewNatGatewayServiceClient(conn)

	req := &ngpb.DeleteNatGatewayReq{}
	err := json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("request is : %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	r, err := c.DeleteNatGateway(ctx, req)
	if err != nil {
		log.Println("could not rpc request: %v", err)
		return
	}
	log.Println("rpc result: %+v", r)
}

func ngCreate(conn grpc.ClientConnInterface, data []byte) {
	c := ngpb.NewNatGatewayServiceClient(conn)

	req := &ngpb.CreateNatGatewayReq{}
	err := json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("request is : %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	r, err := c.CreateNatGateway(ctx, req)
	if err != nil {
		log.Println("could not rpc request: %v", err)
		return
	}
	log.Println("rpc result: %+v", r)
}

func ngGet(conn grpc.ClientConnInterface, data []byte) {
	c := ngpb.NewNatGatewayServiceClient(conn)

	req := &ngpb.GetNatGatewayReq{}
	err := json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("request is : %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	r, err := c.GetNatGateway(ctx, req)
	if err != nil {
		log.Println("could not rpc request: %v", err)
		return
	}
	log.Println("rpc result: %+v", r)
}

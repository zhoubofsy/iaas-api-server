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
	"fmt"
	"log"
	"time"

	ngpb "iaas-api-server/proto/natgateway"
	plpb "iaas-api-server/proto/peerlink"
	sgpb "iaas-api-server/proto/securitygroup"

	"google.golang.org/grpc"
)

func sgOperate(conn grpc.ClientConnInterface, data []byte) {
	c := sgpb.NewSecurityGroupServiceClient(conn)

	req := &sgpb.OperateSecurityGroupReq{}
	err := json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("request is : %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.OperateSecurityGroup(ctx, req)
	if err != nil {
		log.Println("could not rpc request: %v", err)
		return
	}
	log.Println("rpc result: %+v", r)
}

func ngDelete(conn grpc.ClientConnInterface, data []byte) {
	c := ngpb.NewNatGatewayServiceClient(conn)

	req := &ngpb.DeleteNatGatewayReq{}
	err := json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("request is : %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.GetNatGateway(ctx, req)
	if err != nil {
		log.Println("could not rpc request: %v", err)
		return
	}
	log.Println("rpc result: %+v", r)
}

func plCreate(conn grpc.ClientConnInterface, data []byte) {
	c := plpb.NewPeerLinkServiceClient(conn)

	req := &plpb.PeerLinkReq{}
	err := json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("request is : %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	r, err := c.CreatePeerLink(ctx, req)
	if err != nil {
		log.Println("could not rpc request: %v", err)
		return
	}
	log.Println("rpc result: %+v", r)
}

func plDelete(conn grpc.ClientConnInterface, data []byte) {
	c := plpb.NewPeerLinkServiceClient(conn)

	req := &plpb.PeerLinkReq{}
	err := json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("request is : %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	r, err := c.DeletePeerLink(ctx, req)
	if err != nil {
		log.Println("could not rpc request: %v", err)
	}
	log.Println("rpc result: %+v", r)
}

func plGet(conn grpc.ClientConnInterface, data []byte) {
	c := plpb.NewPeerLinkServiceClient(conn)

	req := &plpb.PeerLinkReq{}
	err := json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("request is : %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	r, err := c.GetPeerLink(ctx, req)
	if err != nil {
		log.Println("could not rpc request: %v", err)
		return
	}
	log.Println("rpc result: %+v", r)
}

func main() {
	flag.Parse()

	// Set up a connection to the server.
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	// Contact the server and print out its response.
	if "" == *param {
		log.Fatalf("input request param is null")
	}

	var data []byte = []byte(*param)

	startTime := time.Now()
	defer func() {
		duration := time.Now().Sub(startTime)
		log.Printf("rpc waste time: %+v", duration)
	}()

	switch *method {
	case "sgOperate":
		sgOperate(conn, data)
		break
	case "ngDelete":
		ngDelete(conn, data)
		break
	case "ngCreate":
		ngCreate(conn, data)
		break
	case "ngGet":
		ngGet(conn, data)
		break
	case "plCreate":
		plCreate(conn, data)
		break
	case "plDelete":
		plDelete(conn, data)
		break
	case "plGet":
		plGet(conn, data)
		break
	default:
		fmt.Println("input method not found")
		break
	}
}

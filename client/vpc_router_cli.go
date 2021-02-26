package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	route "iaas-api-server/proto/route"
	vpc "iaas-api-server/proto/vpc"

	"google.golang.org/grpc"
)

func vpcCreate(conn grpc.ClientConnInterface, data []byte) {
	c := vpc.NewVpcServiceClient(conn)

	req := &vpc.CreateVpcReq{}
	err := json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("request is : %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.CreateVpc(ctx, req)
	if err != nil {
		log.Fatalf("could not rpc request: %v", err)
	}
	log.Printf("rpc result: %+v", r)
}

func vpcGetInfo(conn grpc.ClientConnInterface, data []byte) {
	c := vpc.NewVpcServiceClient(conn)

	req := &vpc.GetVpcInfoReq{}
	err := json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("request is : %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.GetVpcInfo(ctx, req)
	if err != nil {
		log.Fatalf("could not rpc request: %v", err)
	}
	log.Printf("rpc result: %+v", r)
}

func vpcSetInfo(conn grpc.ClientConnInterface, data []byte) {
	c := vpc.NewVpcServiceClient(conn)

	req := &vpc.SetVpcInfoReq{}
	err := json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("request is : %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SetVpcInfo(ctx, req)
	if err != nil {
		log.Fatalf("could not rpc request: %v", err)
	}
	log.Printf("rpc result: %+v", r)
}

func routerGet(conn grpc.ClientConnInterface, data []byte) {
	c := route.NewRouterServiceClient(conn)

	req := &route.GetRouterReq{}
	err := json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("request is : %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.GetRouter(ctx, req)
	if err != nil {
		log.Fatalf("could not rpc request: %v", err)
	}
	log.Printf("rpc result: %+v", r)
}

func routerSet(conn grpc.ClientConnInterface, data []byte) {
	c := route.NewRouterServiceClient(conn)

	req := &route.SetRoutesReq{}
	err := json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("request is : %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SetRoutes(ctx, req)
	if err != nil {
		log.Fatalf("could not rpc request: %v", err)
	}
	log.Printf("rpc result: %+v", r)
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
	case "vpcCreate":
		vpcCreate(conn, data)
		break
	case "vpcGetInfo":
		vpcGetInfo(conn, data)
		break
	case "vpcSetInfo":
		vpcSetInfo(conn, data)
		break
	case "routerGet":
		routerGet(conn, data)
		break
	case "routerSet":
		routerSet(conn, data)
		break
	default:
		fmt.Print("method not found")
		break
	}
}

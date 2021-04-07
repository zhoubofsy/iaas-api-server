/*================================================================
*
*  文件名称：firewall_cli.go
*  创 建 者: mongia
*  创建日期：2021年04月06日
*
================================================================*/
package main

import (
	"context"
	"encoding/json"
	fpb "iaas-api-server/proto/firewall"
	"log"
	"time"

	"google.golang.org/grpc"
)

func init() {
	registerFunc("fOperate", fOperate)
	registerFunc("fGet", fGet)
	registerFunc("fCreate", fCreate)
	registerFunc("fDelete", fDelete)
	registerFunc("fUpdate", fUpdate)
}

func fOperate(conn grpc.ClientConnInterface, data []byte) {
	c := fpb.NewFirewallServiceClient(conn)

	req := &fpb.OperateFirewallReq{}
	err := json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("request is : %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	r, err := c.OperateFirewall(ctx, req)
	if err != nil {
		log.Println("could not rpc request: %v", err)
		return
	}
	log.Println("rpc result: %+v", r)
}

func fCreate(conn grpc.ClientConnInterface, data []byte) {
	c := fpb.NewFirewallServiceClient(conn)

	req := &fpb.CreateFirewallReq{}
	err := json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("request is : %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	r, err := c.CreateFirewall(ctx, req)
	if err != nil {
		log.Println("could not rpc request: %v", err)
		return
	}
	log.Println("rpc result: %+v", r)
}

func fGet(conn grpc.ClientConnInterface, data []byte) {
	c := fpb.NewFirewallServiceClient(conn)

	req := &fpb.GetFirewallReq{}
	err := json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("request is : %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	r, err := c.GetFirewall(ctx, req)
	if err != nil {
		log.Println("could not rpc request: %v", err)
		return
	}
	log.Println("rpc result: %+v", r)
}

func fDelete(conn grpc.ClientConnInterface, data []byte) {
	c := fpb.NewFirewallServiceClient(conn)

	req := &fpb.DeleteFirewallReq{}
	err := json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("request is : %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	r, err := c.DeleteFirewall(ctx, req)
	if err != nil {
		log.Println("could not rpc request: %v", err)
		return
	}
	log.Println("rpc result: %+v", r)
}

func fUpdate(conn grpc.ClientConnInterface, data []byte) {
	c := fpb.NewFirewallServiceClient(conn)

	req := &fpb.UpdateFirewallReq{}
	err := json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("request is : %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	r, err := c.UpdateFirewall(ctx, req)
	if err != nil {
		log.Println("could not rpc request: %v", err)
		return
	}
	log.Println("rpc result: %+v", r)
}

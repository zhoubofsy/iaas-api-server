/*================================================================
*
*  文件名称：security_group_cli.go
*  创 建 者: mongia
*  创建日期：2021年03月04日
*
================================================================*/

package main

import (
	"context"
	"encoding/json"
	sgpb "iaas-api-server/proto/securitygroup"
	"log"
	"time"

	"google.golang.org/grpc"
)

func init() {
	registerFunc("sgOperate", sgOperate)
	registerFunc("sgGet", sgGet)
	registerFunc("sgCreate", sgCreate)
	registerFunc("sgDelete", sgDelete)
	registerFunc("sgUpdate", sgUpdate)
}

func sgOperate(conn grpc.ClientConnInterface, data []byte) {
	c := sgpb.NewSecurityGroupServiceClient(conn)

	req := &sgpb.OperateSecurityGroupReq{}
	err := json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("request is : %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	r, err := c.OperateSecurityGroup(ctx, req)
	if err != nil {
		log.Println("could not rpc request: %v", err)
		return
	}
	log.Println("rpc result: %+v", r)
}

func sgCreate(conn grpc.ClientConnInterface, data []byte) {
	c := sgpb.NewSecurityGroupServiceClient(conn)

	req := &sgpb.CreateSecurityGroupReq{}
	err := json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("request is : %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	r, err := c.CreateSecurityGroup(ctx, req)
	if err != nil {
		log.Println("could not rpc request: %v", err)
		return
	}
	log.Println("rpc result: %+v", r)
}

func sgGet(conn grpc.ClientConnInterface, data []byte) {
	c := sgpb.NewSecurityGroupServiceClient(conn)

	req := &sgpb.GetSecurityGroupReq{}
	err := json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("request is : %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	r, err := c.GetSecurityGroup(ctx, req)
	if err != nil {
		log.Println("could not rpc request: %v", err)
		return
	}
	log.Println("rpc result: %+v", r)
}

func sgDelete(conn grpc.ClientConnInterface, data []byte) {
	c := sgpb.NewSecurityGroupServiceClient(conn)

	req := &sgpb.DeleteSecurityGroupReq{}
	err := json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("request is : %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	r, err := c.DeleteSecurityGroup(ctx, req)
	if err != nil {
		log.Println("could not rpc request: %v", err)
		return
	}
	log.Println("rpc result: %+v", r)
}

func sgUpdate(conn grpc.ClientConnInterface, data []byte) {
	c := sgpb.NewSecurityGroupServiceClient(conn)

	req := &sgpb.UpdateSecurityGroupReq{}
	err := json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("request is : %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	r, err := c.UpdateSecurityGroup(ctx, req)
	if err != nil {
		log.Println("could not rpc request: %v", err)
		return
	}
	log.Println("rpc result: %+v", r)
}

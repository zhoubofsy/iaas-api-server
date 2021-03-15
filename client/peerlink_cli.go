/*================================================================
*
*  文件名称：peerlink_cli.go
*  创 建 者: mongia
*  创建日期：2021年03月04日
*
================================================================*/

package main

import (
	"context"
	"encoding/json"
	plpb "iaas-api-server/proto/peerlink"
	"log"
	"time"

	"google.golang.org/grpc"
)

func init() {
	registerFunc("plCreate", plCreate)
	registerFunc("plGet", plGet)
	registerFunc("plDelete", plDelete)
}

func plCreate(conn grpc.ClientConnInterface, data []byte) {
	c := plpb.NewPeerLinkServiceClient(conn)

	req := &plpb.PeerLinkReq{}
	err := json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("request is : %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	r, err := c.GetPeerLink(ctx, req)
	if err != nil {
		log.Println("could not rpc request: %v", err)
		return
	}
	log.Println("rpc result: %+v", r)
}

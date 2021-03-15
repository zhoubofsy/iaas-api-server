package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"iaas-api-server/proto/floatip"
	"log"
	"time"

	"google.golang.org/grpc"
)

func BindFloatIpToInstance(conn grpc.ClientConnInterface, data []byte) {
	c := floatip.NewFloatIpServiceClient(conn)

	req := &floatip.BindFloatIpToInstanceReq{}
	err := json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("request is : %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout)*time.Second)
	defer cancel()

	r, err := c.BindFloatIpToInstance(ctx, req)
	if err != nil {
		log.Fatalf("could not rpc request: %v", err)
	}
	log.Printf("rpc result: %+v", r)
}

func RevokeFloatIpFromInstance(conn grpc.ClientConnInterface, data []byte) {
	c := floatip.NewFloatIpServiceClient(conn)

	req := &floatip.RevokeFloatIpFromInstanceReq{}
	err := json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("request is : %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout)*time.Second)
	defer cancel()

	r, err := c.RevokeFloatIpFromInstance(ctx, req)
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
	case "bind":
		BindFloatIpToInstance(conn, data)
	case "revoke":
		RevokeFloatIpFromInstance(conn, data)
	default:
		fmt.Print("method not found")
		break
	}
}

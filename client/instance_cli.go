package main

import (
	"context"
	"encoding/json"
	"flag"
	"google.golang.org/grpc"
	"iaas-api-server/proto/instance"
	"log"
	"time"
)

func createInstance(){
	// Set up a connection to the server.
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connect to %s failed: %v", *address, err)
	}
	defer conn.Close()

	c := instance.NewInstanceServiceClient(conn)

	var data []byte = []byte(*param)
	req := &instance.CreateInstanceReq{}
	err = json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("rpc req: %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout) * time.Second)
	defer cancel()

	res, err := c.CreateInstance(ctx, req)
	if err != nil {
		log.Fatalf("rpc request failed: %v", err)
	}

	log.Printf("rpc result: %+v", res)
}

func getInstance(){
	// Set up a connection to the server.
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connect to %s failed: %v", *address, err)
	}
	defer conn.Close()

	c := instance.NewInstanceServiceClient(conn)

	var data []byte = []byte(*param)
	req := &instance.GetInstanceReq{}
	err = json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("rpc req: %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout) * time.Second)
	defer cancel()

	res, err := c.GetInstance(ctx, req)
	if err != nil {
		log.Fatalf("rpc request failed: %v", err)
	}

	log.Printf("rpc result: %+v", res)
}

func updateInstanceFlavor(){
	// Set up a connection to the server.
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connect to %s failed: %v", *address, err)
	}
	defer conn.Close()

	c := instance.NewInstanceServiceClient(conn)

	var data []byte = []byte(*param)
	req := &instance.UpdateInstanceFlavorReq{}
	err = json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("rpc req: %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout) * time.Second)
	defer cancel()

	res, err := c.UpdateInstanceFlavor(ctx, req)
	if err != nil {
		log.Fatalf("rpc request failed: %v", err)
	}

	log.Printf("rpc result: %+v", res)
}

func deleteInstance(){
	// Set up a connection to the server.
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connect to %s failed: %v", *address, err)
	}
	defer conn.Close()

	c := instance.NewInstanceServiceClient(conn)

	var data []byte = []byte(*param)
	req := &instance.DeleteInstanceReq{}
	err = json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("rpc req: %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout) * time.Second)
	defer cancel()

	res, err := c.DeleteInstance(ctx, req)
	if err != nil {
		log.Fatalf("rpc request failed: %v", err)
	}

	log.Printf("rpc result: %+v", res)
}

func operateInstance(){
	// Set up a connection to the server.
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connect to %s failed: %v", *address, err)
	}
	defer conn.Close()

	c := instance.NewInstanceServiceClient(conn)

	var data []byte = []byte(*param)
	req := &instance.OperateInstanceReq{}
	err = json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("rpc req: %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout) * time.Second)
	defer cancel()

	res, err := c.OperateInstance(ctx, req)
	if err != nil {
		log.Fatalf("rpc request failed: %v", err)
	}

	log.Printf("rpc result: %+v", res)
}

func main() {
	flag.Parse()

	if "" == *param {
		log.Print("usage:\n    ./cli.exe -address localhost:8080 -method CreateInstance -param '{\"xx\":\"123\"}' -timeout 1")
		return
	}

	if *method == "CreateInstance" {
		createInstance()
	} else if *method == "GetInstance" {
		getInstance()
	} else if *method == "UpdateInstanceFlavor" {
		updateInstanceFlavor()
	} else if *method == "DeleteInstance" {
		deleteInstance()
	} else if *method == "OperateInstance" {
		operateInstance()
	} else if *method == "" {
		log.Fatalf("error: rpc method not found")
	} else {
		log.Fatalf("error: unknown rpc method: ", *method)
	}
}

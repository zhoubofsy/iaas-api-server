package main

import (
	"context"
	"encoding/json"
	"flag"
	"google.golang.org/grpc"
	"iaas-api-server/common"
	"iaas-api-server/proto/instance"
	"log"
	"time"
)

var insConn *grpc.ClientConn
var insCli instance.InstanceServiceClient

// Set up a connection to the server.
func createInstanceServiceClient() (instance.InstanceServiceClient, *grpc.ClientConn) {
	timer := common.NewTimer()
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connect to %s failed: %v", *address, err)
		return nil, nil
	}

	c := instance.NewInstanceServiceClient(conn)
	log.Printf("connect to %s ok, time elapse: %v", *address, timer.Elapse())
	return c, conn
}

func createInstance(){
	timer := common.NewTimer()

	var data []byte = []byte(*param)
	req := &instance.CreateInstanceReq{}
	err := json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("rpc req: %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout) * time.Second)
	defer cancel()

	res, err := insCli.CreateInstance(ctx, req)
	if err != nil {
		log.Fatalf("rpc request failed: %v", err)
	}

	log.Printf("rpc result: %+v, time elpase: %v", res, timer.Elapse())
}

func getInstance(){
	timer := common.NewTimer()
	var data []byte = []byte(*param)
	req := &instance.GetInstanceReq{}
	err := json.Unmarshal(data, req)
	if err != nil {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("rpc req: %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout) * time.Second)
	defer cancel()

	res, err := insCli.GetInstance(ctx, req)
	if err != nil {
		log.Fatalf("rpc request failed: %v", err)
	}

	log.Printf("rpc result: %+v, time elapse: %v", res, timer.Elapse())
}

func updateInstanceFlavor(){
	timer := common.NewTimer()

	var data []byte = []byte(*param)
	req := &instance.UpdateInstanceFlavorReq{}
	err := json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("rpc req: %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout) * time.Second)
	defer cancel()

	res, err := insCli.UpdateInstanceFlavor(ctx, req)
	if err != nil {
		log.Fatalf("rpc request failed: %v", err)
	}

	log.Printf("rpc result: %+v, time elapse: %v", res, timer.Elapse())
}

func deleteInstance(){
	timer := common.NewTimer()

	var data []byte = []byte(*param)
	req := &instance.DeleteInstanceReq{}
	err := json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("rpc req: %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout) * time.Second)
	defer cancel()

	res, err := insCli.DeleteInstance(ctx, req)
	if err != nil {
		log.Fatalf("rpc request failed: %v", err)
	}

	log.Printf("rpc result: %+v, time elapse: %v", res, timer.Elapse())
}

func operateInstance(){
	timer := common.NewTimer()

	var data []byte = []byte(*param)
	req := &instance.OperateInstanceReq{}
	err := json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("rpc req: %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout) * time.Second)
	defer cancel()

	res, err := insCli.OperateInstance(ctx, req)
	if err != nil {
		log.Fatalf("rpc request failed: %v", err)
	}

	log.Printf("rpc result: %+v, time elapse: %v", res, timer.Elapse())
}

func main() {
	flag.Parse()

	if "" == *param {
		log.Print("usage:\n    ./cli.exe -address localhost:8080 -method CreateInstance -param '{\"xx\":\"123\"}' -timeout 1")
		return
	}

	insCli, insConn = createInstanceServiceClient()
	if insConn == nil {
		return
	}
	defer insConn.Close()

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

package main

import (
	"context"
	"encoding/json"
	"flag"
	"google.golang.org/grpc"
	"iaas-api-server/proto/clouddisk"
	"log"
	"time"
)

func createVolume(){
	// Set up a connection to the server.
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connect to %s failed: %v", *address, err)
	}
	defer conn.Close()

	c := clouddisk.NewCloudDiskServiceClient(conn)

	var data []byte = []byte(*param)
	req := &clouddisk.CreateCloudDiskReq{}
	err = json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("rpc req: %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout) * time.Second)
	defer cancel()

	res, err := c.CreateCloudDisk(ctx, req)
	if err != nil {
		log.Fatalf("rpc request failed: %v", err)
	}

	log.Printf("rpc result: %+v", res)
}

func getVolume(){
	// Set up a connection to the server.
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connect to %s failed: %v", *address, err)
	}
	defer conn.Close()

	c := clouddisk.NewCloudDiskServiceClient(conn)

	var data []byte = []byte(*param)
	req := &clouddisk.GetCloudDiskReq{}
	err = json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("rpc req: %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout) * time.Second)
	defer cancel()

	res, err := c.GetCloudDisk(ctx, req)
	if err != nil {
		log.Fatalf("rpc request failed: %v", err)
	}

	log.Printf("rpc result: %+v", res)
}

func updateVolumeInfo(){
	// Set up a connection to the server.
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connect to %s failed: %v", *address, err)
	}
	defer conn.Close()

	c := clouddisk.NewCloudDiskServiceClient(conn)

	var data []byte = []byte(*param)
	req := &clouddisk.ModifyCloudDiskInfoReq{}
	err = json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("rpc req: %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout) * time.Second)
	defer cancel()

	res, err := c.ModifyCloudDiskInfo(ctx, req)
	if err != nil {
		log.Fatalf("rpc request failed: %v", err)
	}

	log.Printf("rpc result: %+v", res)
}

func deleteVolume(){
	// Set up a connection to the server.
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connect to %s failed: %v", *address, err)
	}
	defer conn.Close()

	c := clouddisk.NewCloudDiskServiceClient(conn)

	var data []byte = []byte(*param)
	req := &clouddisk.DeleteCloudDiskReq{}
	err = json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("rpc req: %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout) * time.Second)
	defer cancel()

	res, err := c.DeleteCloudDisk(ctx, req)
	if err != nil {
		log.Fatalf("rpc request failed: %v", err)
	}

	log.Printf("rpc result: %+v", res)
}

func extendsizeVolume(){
	// Set up a connection to the server.
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connect to %s failed: %v", *address, err)
	}
	defer conn.Close()

	c := clouddisk.NewCloudDiskServiceClient(conn)

	var data []byte = []byte(*param)

	req := &clouddisk.ReqizeCloudDiskReq{}

	err = json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("rpc req: %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout) * time.Second)
	defer cancel()

	res, err := c.ReqizeCloudDisk(ctx, req)
	if err != nil {
		log.Fatalf("rpc request failed: %v", err)
	}

	log.Printf("rpc result: %+v", res)
}

func operateVolume(){
	// Set up a connection to the server.
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connect to %s failed: %v", *address, err)
	}
	defer conn.Close()

	c := clouddisk.NewCloudDiskServiceClient(conn)

	var data []byte = []byte(*param)

	req := &clouddisk.OperateCloudDiskReq{}

	err = json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("rpc req: %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout) * time.Second)
	defer cancel()

	res, err := c.OperateCloudDisk(ctx, req)
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

	if *method == "createVolume" {
		createInstance()
	} else if *method == "getVolume" {
		getInstance()
	} else if *method == "updateVolumeInfo" {
		updateVolumeInfo()
	} else if *method == "deleteVolume" {
		deleteVolume()
	} else if *method == "operateVolume" {
		operateInstance()
	}else if *method == "extendsizeVolume" {
		operateInstance()
	}else if *method == "" {
		log.Fatalf("error: rpc method not found")
	} else {
		log.Fatalf("error: unknown rpc method: ", *method)
	}
}

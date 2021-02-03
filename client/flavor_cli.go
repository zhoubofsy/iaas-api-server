package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"time"

	"iaas-api-server/proto/flavor"
	"google.golang.org/grpc"
)

func do_list_flavors() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connect to %s failed: %v", *address, err)
	}
	defer conn.Close()

	c := flavor.NewFlavorServiceClient(conn)

	var data []byte = []byte(*param)
	req := &flavor.ListFlavorsReq{}
	err = json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("rpc req: %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout) * time.Second)
	defer cancel()

	res, err := c.ListFlavors(ctx, req)
	if err != nil {
		log.Fatalf("rpc request failed: %v", err)
	}

	log.Printf("rpc result: %+v", res)
}

func do_get_flavor() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connect to %s failed: %v", *address, err)
	}
	defer conn.Close()

	c := flavor.NewFlavorServiceClient(conn)

	var data []byte = []byte(*param)
	req := &flavor.GetFlavorReq{}
	err = json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("rpc req: %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout) * time.Second)
	defer cancel()

	res, err := c.GetFlavor(ctx, req)
	if err != nil {
		log.Fatalf("rpc request failed: %v", err)
	}

	log.Printf("rpc result: %+v", res)
}

func main() {
	flag.Parse()

	if "" == *param {
		log.Print("usage:\n    ./exe -address localhost:8080 -param '{\"xx\":\"123\"}' -timeout 1")
		return
	}

	if *method == "ListFlavors" {
		do_list_flavors()
	} else if *method == "GetFlavor" {
		do_get_flavor()
	} else if *method == "" {
		log.Fatalf("error: rpc method not found")
	} else {
		log.Fatalf("error: unknown rpc method: ", *method)
	}
}

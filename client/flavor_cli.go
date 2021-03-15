package main

import (
	"context"
	"encoding/json"
	"flag"
	"iaas-api-server/common"
	"log"
	"time"

	"google.golang.org/grpc"
	"iaas-api-server/proto/flavor"
)

var flvConn *grpc.ClientConn
var flvCli flavor.FlavorServiceClient

// Set up a connection to the server.
func createFlavorServiceClient() (flavor.FlavorServiceClient, *grpc.ClientConn) {
	timer := common.NewTimer()
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connect to %s failed: %v", *address, err)
		return nil, nil
	}

	c := flavor.NewFlavorServiceClient(conn)
	log.Printf("connect to %s ok, time elapse: %v", *address, timer.Elapse())
	return c, conn
}

func listFlavors() {
	timer := common.NewTimer()

	var data []byte = []byte(*param)
	req := &flavor.ListFlavorsReq{}
	err := json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("rpc req: %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout) * time.Second)
	defer cancel()

	res, err := flvCli.ListFlavors(ctx, req)
	if err != nil {
		log.Fatalf("rpc request failed: %v", err)
	}

	log.Printf("rpc result: %+v, time elapse: %v", res, timer.Elapse())
}

func getFlavor() {
	timer := common.NewTimer()

	var data []byte = []byte(*param)
	req := &flavor.GetFlavorReq{}
	err := json.Unmarshal(data, req)
	if nil != err {
		log.Fatalf("input param is not json or format error, err: %+v", err)
	}

	log.Printf("rpc req: %+v", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout) * time.Second)
	defer cancel()

	res, err := flvCli.GetFlavor(ctx, req)
	if err != nil {
		log.Fatalf("rpc request failed: %v", err)
	}

	log.Printf("rpc result: %+v, time elapse: %v", res, timer.Elapse())
}

func main() {
	flag.Parse()

	if "" == *param {
		log.Print("usage:\n    ./cli.exe -address localhost:8080 -method GetFlavor -param '{\"xx\":\"123\"}' -timeout 1")
		return
	}

	flvCli, flvConn = createFlavorServiceClient()
	if flvConn == nil {
		return;
	}
	defer flvConn.Close()

	if *method == "ListFlavors" {
		listFlavors()
	} else if *method == "GetFlavor" {
		getFlavor()
	} else if *method == "" {
		log.Fatalf("error: rpc method not found")
	} else {
		log.Fatalf("error: unknown rpc method: ", *method)
	}
}

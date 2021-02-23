package main

import "flag"

var (
	address = flag.String("address", "localhost:8080", "grpc server address, localhost:8080")
	method = flag.String("method", "", "rpc method: ListFlavors...")
	param   = flag.String("param", "", "grpc request param in Json")
	timeout = flag.Int("timeout", 1, "rpc timeout in seconds")
)

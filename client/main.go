/*================================================================
*
*  文件名称：main.go
*  创 建 者: mongia
*  创建日期：2021年02月02日
*
================================================================*/

package main

import (
	"errors"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
)

type subFunc func(conn grpc.ClientConnInterface, data []byte)

var (
	funcs = make(map[string]subFunc)
)

func registerFunc(name string, function subFunc) error {
	if nil != funcs[name] {
		return errors.New("function is exists")
	}

	funcs[name] = function

	return nil
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

	if nil != funcs[*method] {
		funcs[*method](conn, data)
	} else {
		log.Printf("method not found")
	}
}

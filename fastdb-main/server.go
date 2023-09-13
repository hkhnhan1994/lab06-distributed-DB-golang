// server.go

package main

import (
	"fmt"
	"net"
	"net/rpc"
	"github.com/marcelloh/fastdb"
)

func main() {
	db := &fastdb.DB{}
	rpcHandler := &fastdb.FastDBRPC{db}

	rpc.Register(rpcHandler)
	rpc.HandleHTTP()

	listener, err := net.Listen("tcp", "localhost:1234")
	if err != nil {
		fmt.Println("Error starting RPC server:", err)
		return
	}

	fmt.Println("RPC server listening on localhost:1234 ...")

	err = fastdb.Open(":memory:", 0) // Open an in-memory database for testing purposes
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}

	select {}
}

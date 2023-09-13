// client.go

package main

import (
	"fmt"
	"net/rpc"
	"os"
	"github.com/marcelloh/fastdb"
)

func main() {
	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		fmt.Println("Error connecting to RPC server:", err)
		return
	}

	args := &fastdb.RPCArgs{Method: "OpenDB", Bucket: ":memory:", SyncIime: 0}
	result := &fastdb.RPCResult{}

	err = client.Call("FastDBRPC."+args.Method, args, result)
	if err != nil {
		fmt.Println("RPC Call error:", err)
		return
	}

	fmt.Println("Result:", result.Info)

	setArgs := &fastdb.RPCArgs{Method: "SetDB", Bucket: "bucket1", Key: 1, Value: []byte("Hello, RPC!"), SyncIime: 0}
	setResult := &fastdb.RPCResult{}

	err = client.Call("FastDBRPC."+setArgs.Method, setArgs, setResult)
	if err != nil {
		fmt.Println("RPC Call error:", err)
		return
	}

	fmt.Println("Result:", setResult.Info)

	closeArgs := &fastdb.RPCArgs{Method: "CloseDB"}
	closeResult := &fastdb.RPCResult{}
	err = client.Call("FastDBRPC."+closeArgs.Method, closeArgs, closeResult)
	if err != nil {
		fmt.Println("RPC Call error:", err)
		return
	}

	fmt.Println("Result:", closeResult.Info)
}

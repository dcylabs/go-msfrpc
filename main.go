package main

import (
	"fmt"
	"os"

	"github.com/dcylabs/go-msfrpc/msfrpc"
)

func main() {

	rpc := msfrpc.NewMsfrpc(os.Args[1], os.Args[2], "/api", os.Args[3], os.Args[4], true)
	err := rpc.Login()
	fmt.Printf("Error  of 'auth.login':      %v\n", err)
	result, err := rpc.Call("console.create", []interface{}{})
	fmt.Printf("Result of 'console.create':  %v\n", result)
	fmt.Printf("Error  of 'console.create':  %v\n", err)

}

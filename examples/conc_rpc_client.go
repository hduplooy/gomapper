package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc/jsonrpc"

	mp "github.com/hduplooy/gomapper"
)

// Function to do the actual calling of the json-rpc server to summate a slice of integers for us
func rpcAdd(add string, args []int) int {
	client, err := net.Dial("tcp", add)
	if err != nil {
		log.Fatal("dialing:", err)
		return -1
	}
	var reply int
	c := jsonrpc.NewClient(client)
	err = c.Call("Agg.Sum", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
		return -1
	}
	return reply
}

func main() {
	// IP Addresses of machines that we want to use to process our function
	// Obviously the correct ones have to be used here
	addrs := []string{"127.0.0.1:9991", "127.0.0.1:9992", "127.0.0.1:9993"}

	// Call MapConc to distribute the work for us
	sums, err := mp.MapConc(3, func(vals []interface{}, pos int) (interface{}, error) {
		// Convert the interface{} to ints for the call to the rpc function
		ints := make([]int, len(vals))
		for i := 0; i < len(vals); i++ {
			ints[i] = vals[i].(int)
		}
		// let rpcAdd do the json-rpc call to the different servers to do the summations for us (use pos to select the server)
		return rpcAdd(addrs[pos], ints), nil
	}, []int{1, 2, 3}, []int{5, 6, 7}, []int{9, 10, 11})
	if err != nil {
		log.Println(err)
		return
	}

	// Just add the results together
	ans := mp.Fold(func(val1, val2 interface{}) interface{} {
		return val1.(int) + val2.(int)
	}, mp.ToInterfaceArr(sums)...)
	fmt.Printf("Sum=%d\n", ans.(int))
}

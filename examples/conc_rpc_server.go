package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
)

type Agg struct{}

// Actual function that does the summation of integers
func (t *Agg) Sum(args *[]int, reply *int) error {
	var ans int
	for _, val := range *args {
		ans += val
	}
	*reply = ans
	return nil
}

func main() {
	fmt.Printf("all args: %d first arg: %s\n", len(os.Args), os.Args[1])
	if len(os.Args) < 2 {
		log.Fatal("Please provide a port number for the server")
	}
	cal := new(Agg)
	server := rpc.NewServer()
	server.Register(cal)
	server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
	listener, e := net.Listen("tcp", ":"+os.Args[1])
	if e != nil {
		log.Fatal("listen error:", e)
	}
	for {
		if conn, err := listener.Accept(); err != nil {
			log.Fatal("accept error: " + err.Error())
		} else {
			log.Printf("new connection established\n")
			go server.ServeCodec(jsonrpc.NewServerCodec(conn))
		}
	}
}

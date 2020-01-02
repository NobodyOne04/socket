package main

import (
	"flag"
	"fmt"
	"sync"
	"./server"
	"./client"
)

var (
	host     = flag.String("h", "localhost", "host")
	port     = flag.Int("p", 9090, "port")
	asserver   = flag.Bool("s",false,"Start as server")
	serverIP = flag.String("ip","127.0.0.1","Server ip")
	waitGroup sync.WaitGroup
)

func recovery() { if recv:= recover(); recv != nil {fmt.Println(recv)} }

func main() {
	defer recovery()
	flag.Parse()
	if *asserver {
		waitGroup.Add(1)
		go server.StartServer(*host, *port)
	} else {
		waitGroup.Add(1)
		go client.StartClient(*serverIP, *port)
	}
	waitGroup.Wait()
}

package main

import (
	servers "TCPServer/server"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go servers.StartTCPServer(&wg)
	go servers.StartWSServer(&wg)
	// keeps the main thread waiting and doesn't lets it exit
	wg.Wait()
}

package main

import (
	databases "TCPServer/internal/database"
	servers "TCPServer/internal/server"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(4)
	dbConn := databases.ConnectToDB(&wg)
	go servers.StartWebServer(&wg, &dbConn)
	go servers.StartTCPServer(&wg)
	go servers.StartWSServer(&wg)
	// keeps the main thread waiting and doesn't lets it exit
	wg.Wait()
}

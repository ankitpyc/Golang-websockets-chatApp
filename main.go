package main

import (
	databases "TCPServer/internal/database"
	servers "TCPServer/internal/server"
	"os"
	"os/signal"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(4)
	dbConn := databases.ConnectToDB(&wg)
	go servers.StartWebServer(&wg, &dbConn)
	go servers.StartTCPServer(&wg)
	go servers.StartWSServer(&wg, &dbConn)
	// keeps the main thread waiting and doesn't let it exit
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt)
	<-terminate
}

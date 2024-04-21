package servers

import (
	"fmt"
	"log"
	"net"
	"sync"
)

func StartTCPServer(wg *sync.WaitGroup) {
	defer wg.Done()
	tcpListener, err := net.Listen("tcp", ":4445")
	if err != nil {
		log.Fatal("error while listening for tcp connections at port 4444")
	}
	fmt.Println("Listening for TCP connections port 4444")
	conn, _ := tcpListener.Accept()
	go handleConnection(conn)
}

func handleConnection(conn net.Conn) {
	panic("unimplemented")
}

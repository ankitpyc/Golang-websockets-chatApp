package servers

import (
	cache "TCPServer/redis-cache"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func StartWSServer(wg *sync.WaitGroup) {
	defer wg.Done()
	go cache.InitRedisClient()
	hub := newSocketHub()
	go hub.startSocketHub()
	go hub.notifyOnlineUsers()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		log.Print("recieving connection")
		serveWS(hub, w, r)
	})

	fmt.Println("Listening for websockets on port ", 2019)
	http.ListenAndServe(":2019", corsHandler(http.DefaultServeMux))
}

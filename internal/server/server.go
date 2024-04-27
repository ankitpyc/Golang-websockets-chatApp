package servers

import (
	cache "TCPServer/internal/redis-cache"
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
	hub := newSocketHub(&SocketHub{})
	go hub.startSocketHub()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		log.Print("recieving connection")
		serveWS(hub, w, r)
	})
	fmt.Println("Listening for websockets on port ", 2019)
	http.ListenAndServe(":2023", corsHandler(http.DefaultServeMux))
}

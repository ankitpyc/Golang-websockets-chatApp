package servers

import (
	models "TCPServer/internal/database/models"
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

func StartWSServer(wg *sync.WaitGroup, db *models.DBServer) {
	defer wg.Done()
	go cache.InitRedisClient()
	hub := newSocketHub(&SocketHub{})
	hub.DB = db
	go hub.startSocketHub()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Print("receiving connection")
		serveWS(hub, w, r)
	})
	fmt.Println("Listening for websockets on port ", 2019)
	http.ListenAndServe(":2023", corsHandler(http.DefaultServeMux))
}

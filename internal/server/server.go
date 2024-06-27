package servers

import (
	models "TCPServer/internal/database"
	corsmiddleware "TCPServer/internal/middleware/cors"
	cache "TCPServer/internal/redis-cache"
	client "TCPServer/internal/server/client"
	"fmt"
	"log"
	"net/http"
	"sync"
)

func StartWSServer(wg *sync.WaitGroup, db *models.DBServer) {
	defer wg.Done()
	go cache.InitRedisClient()
	hub := client.NewSocketHub()
	hub.DB = db
	go hub.StartSocketHub()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Print("receiving connection")
		client.ServeWS(&hub, w, r)
	})
	fmt.Println("Listening for websockets on port ", 2019)
	http.ListenAndServe(":2023", corsmiddleware.CorsHandler(http.DefaultServeMux))
}

package servers

import (
	models "TCPServer/internal/database"
	corsmiddleware "TCPServer/internal/middleware/cors"
	client "TCPServer/internal/server/client"
	"fmt"
	"net/http"
	"sync"
)

func StartWSServer(wg *sync.WaitGroup, db *models.DBServer) {
	defer wg.Done()
	hub := client.NewSocketHub(db)
	go hub.StartSocketHub()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		client.ServeWS(&hub, w, r)
	})

	fmt.Println("Listening for websockets on port ", 2019)
	http.ListenAndServe(":2023", corsmiddleware.CorsHandler(http.DefaultServeMux))
}

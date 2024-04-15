package servers

import (
	cache "TCPServer/redis-cache"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Server struct {
	conns map[string]*websocket.Conn
}

type Message struct {
	MessageType string `json:"messageType"`
	Text        string `json:"text"`
	ID          string `json:"id"`
	Date        int64  `json:"date"`
}

type Client struct {
}

func (server *Server) addConn(userId string, conn *websocket.Conn) {
	server.conns[userId] = conn
}

func (server *Server) removeConn(userId string) {
	delete(server.conns, userId)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func echo(w http.ResponseWriter, r *http.Request, server *Server) {
	connection, _ := upgrader.Upgrade(w, r, nil)
	fmt.Println(r.RemoteAddr + " joined the chat ")
	defer connection.Close()
	for {
		var chatMessage Message
		messageType, message, err := connection.ReadMessage()
		fmt.Println("recieved message", message)

		fmt.Println("recieved message", json.Unmarshal(message, &chatMessage))
		fmt.Println("MESSAGE TYPE IS ", chatMessage.MessageType)

		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				server.removeConn(r.RemoteAddr)
				fmt.Printf("Connection closed normally")
				return
			}

			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				fmt.Println("Connection closed Abruptly by", connection.RemoteAddr())
				server.removeConn(r.RemoteAddr)
				return
			}

			if websocket.IsCloseError(err, websocket.CloseAbnormalClosure) {
				fmt.Println("Connection closed Abnormally by", connection.RemoteAddr())
				server.removeConn(r.RemoteAddr)
				return
			}

			fmt.Println("Error while reading message", err)
			return
		}

		if chatMessage.MessageType == "CONNECT_PING" {
			server.addConn(chatMessage.ID, connection)
		} else {
			switch messageType {
			case websocket.TextMessage:
				broadCastMessage(server, chatMessage, connection)
				messageHandler(message, connection)
			case websocket.CloseMessage:
				fmt.Println("close message recieved from ", r.RemoteAddr)
			}
		}
	}
}

func broadCastMessage(server *Server, message Message, connection *websocket.Conn) {
	jsonres, _ := json.Marshal(message)
	for k := range server.conns {
		server.conns[k].WriteMessage(websocket.TextMessage, jsonres)
	}
}

func messageHandler(message []byte, conn *websocket.Conn) {
	fmt.Println(conn.RemoteAddr().String() + string(message))
}

func StartWSServer(wg *sync.WaitGroup) {
	defer wg.Done()
	cache.InitRedisClient()
	wsConnSer := Server{
		conns: make(map[string]*websocket.Conn),
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		echo(w, r, &wsConnSer)
	})

	fmt.Println("Listening for websockets on port ", 2019)
	go http.ListenAndServe(":2019", corsHandler(http.DefaultServeMux))

}

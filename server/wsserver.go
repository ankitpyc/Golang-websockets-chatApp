package servers

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func newSocketHub(socketHub *SocketHub) *SocketHub {
	return &SocketHub{
		unsubcribe:       make(chan *Client),
		subcribe:         make(chan *Client),
		broadCastMessage: make(chan Message),
		connectionsMap:   make(map[string]*Client),
	}
}

func (hub *SocketHub) notifyOnlineUsers() {
	for {
		select {
		case <-time.After(15 * time.Second):
			// Send a message to all WebSocket clients
			log.Print("Broadcasting online statuses")
			for _, client := range hub.connectionsMap {
				fmt.Println("writing to client ", client.id)
				message := Message{
					MessageType: "CONNECT_PING",
					UserName:    client.username,
					ID:          client.id,
					Text:        "",
					RecieverID:  "",
					Date:        0,
				}
				hub.broadCastMessage <- message
			}
		}
	}
}

func (hub *SocketHub) startSocketHub() {
	go hub.notifyOnlineUsers()
	for {
		select {
		case client := <-hub.subcribe:
			log.Printf("client subscribed")
			hub.connectionsMap[client.id] = client
		case client := <-hub.unsubcribe:
			log.Printf("client subscribed")
			delete(hub.connectionsMap, client.id)
			client.conn.Close()
		case messages := <-hub.broadCastMessage:
			sendBroadCastMessage(hub, messages)
		}
	}
}

func sendBroadCastMessage(hub *SocketHub, chatMessage Message) {
	log.Println("broadcasting")
	jsonres, _ := json.Marshal(chatMessage)
	for _, client := range hub.connectionsMap {
		fmt.Println("writing to client ", client.id)
		client.conn.WriteMessage(websocket.TextMessage, jsonres)
	}
}

func sendMessage(hub *SocketHub, chatMessage Message) {
	byteMessage, _ := json.Marshal(chatMessage)
	recieverConn := hub.connectionsMap[chatMessage.RecieverID]
	recieverConn.conn.WriteMessage(websocket.TextMessage, byteMessage)
}

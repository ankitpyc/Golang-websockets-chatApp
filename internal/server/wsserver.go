package servers

import (
	"TCPServer/internal/domain"
	"encoding/json"
	"fmt"
	"gorm.io/gorm/logger"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func newSocketHub(socketHub *SocketHub) *SocketHub {
	return &SocketHub{
		unsubcribe:       make(chan *Client),
		subcribe:         make(chan *Client),
		broadCastMessage: make(chan domain.Message),
		connectionsMap:   make(map[string]*Client),
	}
}

func (hub *SocketHub) notifyOnlineUsers() {
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			// Send a message to all WebSocket clients
			for _, client := range hub.connectionsMap {
				message := domain.Message{
					MessageType: "CONNECT_PING",
					UserName:    client.username,
					ID:          client.id,
					Text:        "",
					ReceiverID:  "",
					Date:        "0",
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
			log.Println(logger.Blue, "client subscribed : "+client.id)
			hub.connectionsMap[client.id] = client
			break
		case client := <-hub.unsubcribe:
			log.Printf("client unsubscribed")
			delete(hub.connectionsMap, client.id)
			client.conn.Close()
		case messages := <-hub.broadCastMessage:
			sendBroadCastMessage(hub, messages)
		}
	}
}

func sendBroadCastMessage(hub *SocketHub, chatMessage domain.Message) {
	jsonres, _ := json.Marshal(chatMessage)
	for _, client := range hub.connectionsMap {
		fmt.Println("writing to client ", client.id)
		client.Lock()
		err := client.conn.WriteMessage(websocket.TextMessage, jsonres)
		client.Unlock()
		if err != nil {
			return
		}
	}
}

func sendMessage(hub *SocketHub, chatMessage domain.Message) {
	byteMessage, _ := json.Marshal(chatMessage)
	recieverConn := hub.connectionsMap[chatMessage.ReceiverID]
	recieverConn.conn.WriteMessage(websocket.TextMessage, byteMessage)
}

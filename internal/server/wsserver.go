package servers

import (
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
		broadCastMessage: make(chan Message),
		connectionsMap:   make(map[string]*Client),
	}
}

func (hub *SocketHub) notifyOnlineUsers() {
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			// Send a message to all WebSocket clients
			log.Println("Broadcasting online statuses")
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

			log.Println(logger.Blue, "client subscribed : "+client.id)
			hub.Lock()
			hub.connectionsMap[client.id] = client
			hub.Unlock()
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

func sendBroadCastMessage(hub *SocketHub, chatMessage Message) {
	log.Println("broadcasting")
	jsonres, _ := json.Marshal(chatMessage)
	for _, client := range hub.connectionsMap {
		fmt.Println("writing to client ", client.id)
		err := client.conn.WriteMessage(websocket.TextMessage, jsonres)
		if err != nil {
			return
		}
	}
}

func sendMessage(hub *SocketHub, chatMessage Message) {
	byteMessage, _ := json.Marshal(chatMessage)
	recieverConn := hub.connectionsMap[chatMessage.RecieverID]
	recieverConn.conn.WriteMessage(websocket.TextMessage, byteMessage)
}

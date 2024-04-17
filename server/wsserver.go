package servers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

func newSocketHub() *SocketHub {
	return &SocketHub{
		unsubcribe:       make(chan *Client),
		subcribe:         make(chan *Client),
		broadCastMessage: make(chan Message),
		connectionsMap:   make(map[string]*Client),
	}
}

func (hub *SocketHub) startSocketHub() {
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

// func parseMessageAndSend(hub *SocketHub, message Message) {
// 	var chatMessage Message
// 	switch message.MessageType {
// 	case "BROADCAST":
// 		sendBroadCastMessage(hub, chatMessage)
// 	case "CHAT_MESSAGE":
// 		sendBroadCastMessage(hub, chatMessage)
// 	}
// }

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
	recieverConn := hub.connectionsMap[chatMessage.recieverID]
	recieverConn.conn.WriteMessage(websocket.TextMessage, byteMessage)
}
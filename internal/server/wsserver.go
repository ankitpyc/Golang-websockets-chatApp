package servers

import (
	"TCPServer/internal/domain"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm/logger"

	"github.com/gorilla/websocket"
)

// newSocketHub initializes a new SocketHub and returns a pointer to it.
func newSocketHub(socketHub *SocketHub) *SocketHub {
	return &SocketHub{
		unsubcribe:       make(chan *Client),        // Channel for unsubscribing clients
		subcribe:         make(chan *Client),        // Channel for subscribing clients
		broadCastMessage: make(chan domain.Message), // Channel for broadcasting messages
		connectionsMap:   make(map[string]*Client),  // Map to hold client connections
	}
}

// notifyOnlineUsers periodically sends a "CONNECT_PING" message to all connected clients.
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
				hub.broadCastMessage <- message // Broadcast the message
			}
		}
	}
}

// startSocketHub starts the main event loop for the SocketHub.
func (hub *SocketHub) startSocketHub() {
	go hub.notifyOnlineUsers() // Start the online user notification goroutine
	for {
		select {
		case client := <-hub.subcribe:
			log.Println(logger.Blue, "client subscribed : "+client.id)
			client.hub.Lock()
			hub.connectionsMap[client.id] = client // Add client to connections map
			client.hub.Unlock()
		case client := <-hub.unsubcribe:
			log.Printf("client unsubscribed")
			client.hub.Lock()
			delete(hub.connectionsMap, client.id) // Remove client from connections map
			_ = client.conn.Close()
			client.hub.Lock()
			// Close the client's WebSocket connection
		case messages := <-hub.broadCastMessage:
			sendBroadCastMessage(hub, messages) // Broadcast the message to all clients
		}
	}
}

// sendBroadCastMessage sends a broadcast message to all connected clients.
func sendBroadCastMessage(hub *SocketHub, chatMessage domain.Message) {
	jsonres, _ := json.Marshal(chatMessage) // Marshal the message to JSON
	for _, client := range hub.connectionsMap {
		fmt.Println("writing to client ", client.id)
		client.Lock()
		err := client.conn.WriteMessage(websocket.TextMessage, jsonres) // Write message to WebSocket
		client.Unlock()
		if err != nil {
			return
		}
	}
}

// sendMessage sends a message to a specific client identified by ReceiverID.
func sendMessage(hub *SocketHub, chatMessage domain.Message) {
	byteMessage, _ := json.Marshal(chatMessage) // Marshal the message to JSON
	recieverConn := hub.connectionsMap[chatMessage.ReceiverID]
	recieverConn.conn.WriteMessage(websocket.TextMessage, byteMessage) // Write message to WebSocket
}

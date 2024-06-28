package servers

import (
	models "TCPServer/internal/database"
	dto "TCPServer/internal/domain/dto"

	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// newSocketHub initializes a new SocketHub and returns a pointer to it.
func NewSocketHub(db *models.DBServer) SocketHub {

	shub := SocketHub{
		Unsubcribe:       make(chan *Client),       // Channel for unsubscribing clients
		Subcribe:         make(chan *Client),       // Channel for subscribing clients
		BroadCastMessage: make(chan dto.Message),   // Channel for broadcasting messages
		ConnectionsMap:   make(map[string]*Client), // Map to hold client connections
	}
	shub.DB = db
	return shub
}

// notifyOnlineUsers periodically sends a "CONNECT_PING" message to all connected clients.
func (hub *SocketHub) NotifyOnlineUsers() {
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			// Send a message to all WebSocket clients
			for _, client := range hub.ConnectionsMap {
				message := dto.Message{
					MessageType: "CONNECT_PING",
					UserName:    client.Username,
					ID:          client.Id,
					Text:        "",
					ReceiverID:  "",
					Date:        "0",
				}
				hub.BroadCastMessage <- message // Broadcast the message
			}
		}
	}
}

// startSocketHub starts the main event loop for the SocketHub.
func (hub *SocketHub) StartSocketHub() {
	go hub.NotifyOnlineUsers() // Start the online user notification goroutine
	for {
		select {
		case client := <-hub.Subcribe:
			fmt.Println("Client is being subscribed ", client.Id)
			client.Hub.Lock()
			hub.ConnectionsMap[client.Id] = client // Add client to connections map
			client.Hub.Unlock()
		case client := <-hub.Unsubcribe:
			log.Printf("client unsubscribed")
			client.Hub.Lock()
			delete(hub.ConnectionsMap, client.Id) // Remove client from connections map
			_ = client.Conn.Close()
			client.Hub.Unlock()
			// Close the client's WebSocket connection
		case messages := <-hub.BroadCastMessage:
			sendBroadCastMessage(hub, messages) // Broadcast the message to all clients
		}
	}
}

// sendBroadCastMessage sends a broadcast message to all connected clients.
func sendBroadCastMessage(hub *SocketHub, chatMessage dto.Message) {
	jsonres, _ := json.Marshal(chatMessage) // Marshal the message to JSON
	for _, client := range hub.ConnectionsMap {
		fmt.Println("writing to client ", client.Id)
		client.Lock()
		err := client.Conn.WriteMessage(websocket.TextMessage, jsonres) // Write message to WebSocket
		client.Unlock()
		if err != nil {
			return
		}
	}
}

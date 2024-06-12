package servers

import (
	databases "TCPServer/internal/database/handlers"
	"TCPServer/internal/domain"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
)

// readWS reads messages from the WebSocket connection of a client.
func readWS(client *Client) {
	connection := client.conn
	defer func() {
		client.hub.unsubcribe <- client // Unsubscribe the client from the hub
		client.conn.Close()             // Close the WebSocket connection
	}()
	for {
		var chatMessage domain.Message
		_, message, err := client.conn.ReadMessage() // Read message from WebSocket
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				fmt.Printf("Connection closed normally")
			}
			// TODO: Handle concurrent deletes
			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				delete(client.hub.connectionsMap, client.id) // Remove client from connections map
				fmt.Println("Connection closed abruptly by", connection.RemoteAddr())
				closeMessage := &domain.Message{
					MessageType: "CLOSE",
					UserName:    client.username,
					ID:          client.id,
					Text:        "",
					ReceiverID:  "",
					Date:        "0",
				}
				client.hub.broadCastMessage <- *closeMessage // Broadcast close message
				return
			}

			if websocket.IsCloseError(err, websocket.CloseAbnormalClosure) {
				fmt.Println("Connection closed abnormally by", connection.RemoteAddr())
				return
			}

			fmt.Println("Error while reading message", err)
			return
		}
		err = json.Unmarshal(message, &chatMessage) // Unmarshal JSON message
		if err != nil {
			fmt.Print(err)
			os.Exit(0)
		}
		log.Println("message is ", chatMessage.ID)
		switch chatMessage.MessageType {
		case "CONNECT_PING":
			log.Printf("Broadcasting user the message")
			client.id = chatMessage.ID
			client.username = chatMessage.UserName
			client.hub.subcribe <- client              // Subscribe client to hub
			client.hub.broadCastMessage <- chatMessage // Broadcast the message
		case "BROADCAST":
			log.Printf("Broadcasting")
			client.hub.broadCastMessage <- chatMessage // Broadcast the message
		case "ACK":
			log.Printf("Received ACK")
			client.hub.connectionsMap[chatMessage.ReceiverID].message <- chatMessage // Send ACK to receiver
		case "CHAT_MESSAGE":
			client.hub.connectionsMap[chatMessage.ReceiverID].message <- chatMessage // Send chat message to receiver
		}
	}
}

// WriteMessage handles writing messages to the WebSocket connection of a client.
func WriteMessage(client *Client) {
	log.Println("Writing message go routine triggered")
	chatHandler := databases.ChatHandler{DBServer: client.hub.DB}
	defer func() {
		client.hub.unsubcribe <- client // Unsubscribe the client from the hub
		err := client.conn.Close()      // Close the WebSocket connection
		if err != nil {
			return
		}
	}()
	for {
		select {
		case mess := <-client.message: // Receive message from channel
			client.Lock()
			if mess.MessageType == "ACK" {
				log.Printf("Received ACK %s -> %s , Status : %s ", mess.ID, mess.ReceiverID, mess.MessageDeliveryStatus)
				byteMessage, _ := json.Marshal(mess)                                // Marshal message to JSON
				err := client.conn.WriteMessage(websocket.TextMessage, byteMessage) // Write message to WebSocket
				if err != nil {
					return
				}
			} else {
				err := chatHandler.PersistMessages(&mess)                 // Persist message to database
				ack, _ := chatHandler.SendAcknowledgement(&mess)          // Send ACK for the message
				client.hub.connectionsMap[ack.ReceiverID].message <- *ack // Send ACK to receiver
				byteMessage, err := json.Marshal(mess)                    // Marshal message to JSON
				if err != nil {
					log.Println("error while marshalling message", err)
				}
				err = client.conn.WriteMessage(websocket.TextMessage, byteMessage) // Write message to WebSocket
				if err != nil {
					return
				}
			}
			client.Unlock()
		}
	}
}

// serveWS upgrades the HTTP server connection to the WebSocket protocol and starts read and write goroutines for the client.
func serveWS(hub *SocketHub, w http.ResponseWriter, r *http.Request) {
	connection, _ := upgrader.Upgrade(w, r, nil) // Upgrade the HTTP connection to WebSocket
	client := newClient(connection, hub)
	// Individual client threads to read and write from socket
	go readWS(client)       // Start read goroutine
	go WriteMessage(client) // Start write goroutine
}

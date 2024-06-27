package servers

import (
	databases "TCPServer/internal/database/handlers"
	dto "TCPServer/internal/domain/dto"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

// readWS reads messages from the WebSocket connection of a client.
func readWS(client *Client) {
	connection := client.Conn
	defer func() {
		client.Hub.Unsubcribe <- client // Unsubscribe the client from the hub
		client.Conn.Close()             // Close the WebSocket connection
	}()
	for {
		var chatMessage dto.Message
		_, message, err := client.Conn.ReadMessage() // Read message from WebSocket
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				fmt.Printf("Connection closed normally")
			}
			// TODO: Handle concurrent deletes
			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				client.Hub.Lock()
				delete(client.Hub.ConnectionsMap, client.Id) // Remove client from connections map
				client.Hub.Unlock()
				fmt.Println("Connection closed abruptly by", connection.RemoteAddr())
				closeMessage := &dto.Message{
					MessageType: "CLOSE",
					UserName:    client.Username,
					ID:          client.Id,
					Text:        "",
					ReceiverID:  "",
					Date:        "0",
				}
				client.Hub.BroadCastMessage <- *closeMessage // Broadcast close message
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
		switch chatMessage.MessageType {
		case "CONNECT_PING":
			log.Printf("Broadcasting user the message")
			client.Id = chatMessage.ID
			client.Username = chatMessage.UserName
			client.Hub.Subcribe <- client              // Subscribe client to hub
			client.Hub.BroadCastMessage <- chatMessage // Broadcast the message
		case "BROADCAST":
			log.Printf("Broadcasting")
			client.Hub.BroadCastMessage <- chatMessage // Broadcast the message
		case "ACK":
			log.Printf("Received ACK")
			client.Hub.ConnectionsMap[chatMessage.ReceiverID].Message <- chatMessage // Send ACK to receiver
		case "CHAT_MESSAGE":
			client.Hub.ConnectionsMap[chatMessage.ReceiverID].Message <- chatMessage // Send chat message to receiver
		}
	}
}

// WriteMessage handles writing messages to the WebSocket connection of a client.
func WriteMessage(client *Client) {
	chatHandler := databases.ChatHandler{DBServer: client.Hub.DB}
	defer func() {
		client.Hub.Unsubcribe <- client // Unsubscribe the client from the hub
		err := client.Conn.Close()      // Close the WebSocket connection
		if err != nil {
			return
		}
	}()
	for {
		select {
		case mess := <-client.Message: // Receive message from channel
			client.Lock()
			if mess.MessageType == "ACK" {
				log.Printf("Received ACK %s -> %s , Status : %s ", mess.ID, mess.ReceiverID, mess.MessageDeliveryStatus)
				byteMessage, _ := json.Marshal(mess)                                // Marshal message to JSON
				err := client.Conn.WriteMessage(websocket.TextMessage, byteMessage) // Write message to WebSocket
				if err != nil {
					return
				}
			} else {
				_ = chatHandler.PersistMessages(&mess) // Persist message to database
				//TODO : this probably needs to be e moved to a service (eg :- Nessaging Service)
				fmt.Print("chat id is ", mess.ChatId)
				ack, _ := chatHandler.SendAcknowledgement(&mess)          // Send ACK for the message
				client.Hub.ConnectionsMap[ack.ReceiverID].Message <- *ack // Send ACK to receiver
				byteMessage, err := json.Marshal(mess)                    // Marshal message to JSON
				if err != nil {
					log.Println("error while marshalling message", err)
				}
				err = client.Conn.WriteMessage(websocket.TextMessage, byteMessage) // Write message to WebSocket
				if err != nil {
					return
				}
			}
			client.Unlock()
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// serveWS upgrades the HTTP server connection to the WebSocket protocol and starts read and write goroutines for the client.
func ServeWS(hub *SocketHub, w http.ResponseWriter, r *http.Request) {
	connection, _ := upgrader.Upgrade(w, r, nil) // Upgrade the HTTP connection to WebSocket
	client := newClient(connection, hub)
	// Individual client threads to read and write from socket
	go readWS(client)       // Start read goroutine
	go WriteMessage(client) // Start write goroutine
}

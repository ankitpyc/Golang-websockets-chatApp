package servers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

func readWS(client *Client) {
	connection := client.conn
	defer func() {
		client.hub.unsubcribe <- client
		client.conn.Close()
	}()
	for {
		var chatMessage Message
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				fmt.Printf("Connection closed normally")
			}

			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				delete(client.hub.connectionsMap, client.id)
				fmt.Println("Connection closed Abruptly by", connection.RemoteAddr())
				closeMessage := &Message{
					MessageType: "CLOSE",
					UserName:    client.username,
					ID:          client.id,
					Text:        "",
					ReceiverID:  "",
					Date:        "0",
				}
				client.hub.broadCastMessage <- *closeMessage
				return
			}

			if websocket.IsCloseError(err, websocket.CloseAbnormalClosure) {
				fmt.Println("Connection closed Abnormally by", connection.RemoteAddr())
				return
			}

			fmt.Println("Error while reading message", err)
			return
		}
		err = json.Unmarshal(message, &chatMessage)
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
			client.hub.subcribe <- client
			client.hub.broadCastMessage <- chatMessage
		case "BROADCAST":
			log.Printf("Broadcasting")
			client.hub.broadCastMessage <- chatMessage
		case "ACK":
			log.Printf("Recieved ACK")
			client.hub.connectionsMap[chatMessage.ReceiverID].message <- chatMessage
		case "CHAT_MESSAGE":
			mess := Message{
				MessageType:           "ACK",
				Text:                  "",
				MessageId:             chatMessage.MessageId,
				MessageDeliveryStatus: "SENT",
				UserName:              chatMessage.UserName,
				ID:                    chatMessage.ReceiverID,
				ReceiverID:            chatMessage.ID,
				Date:                  "0",
			}
			client.message <- mess
			time.Sleep(5000 * time.Millisecond)
			client.hub.connectionsMap[chatMessage.ReceiverID].message <- chatMessage
		}
	}
}

func WriteMessage(client *Client) {
	defer func() {
		client.hub.unsubcribe <- client
		err := client.conn.Close()
		if err != nil {
			return
		}
	}()
	for {
		select {
		case mess := <-client.message:
			byteMessage, err := json.Marshal(mess)
			if err != nil {
				log.Println("error while marshalling message", err)
			}
			client.Lock()
			err = client.conn.WriteMessage(websocket.TextMessage, byteMessage)
			client.Unlock()
			if err != nil {
				return
			}
		}
	}

}

func serveWS(hub *SocketHub, w http.ResponseWriter, r *http.Request) {
	connection, _ := upgrader.Upgrade(w, r, nil)
	client := newClient(connection, hub)
	//individual client threads to read and write from socket
	go readWS(client)
	go WriteMessage(client)
}

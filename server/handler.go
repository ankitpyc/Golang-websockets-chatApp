package servers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

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
				return
			}

			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				fmt.Println("Connection closed Abruptly by", connection.RemoteAddr())
				return
			}

			if websocket.IsCloseError(err, websocket.CloseAbnormalClosure) {
				fmt.Println("Connection closed Abnormally by", connection.RemoteAddr())
				return
			}

			fmt.Println("Error while reading message", err)
			return
		}
		error := json.Unmarshal(message, &chatMessage)
		if error != nil {
			fmt.Print(error)
			os.Exit(0)
		}
		log.Println("message is ", chatMessage.ID)

		switch chatMessage.MessageType {
		case "CONNECT_PING":
			client.id = chatMessage.ID
			client.hub.subcribe <- client
			client.hub.broadCastMessage <- chatMessage
		case "BROADCAST":
			log.Printf("Broadcasting")
			client.hub.broadCastMessage <- chatMessage
		case "CHAT_MESSAGE":
			client.hub.connectionsMap[chatMessage.recieverID].message <- chatMessage
		}
	}
}

func WriteMessage(client *Client) {
	defer func() {
		client.hub.unsubcribe <- client
		client.conn.Close()
	}()
	for {
		select {
		case mess := <-client.message:
			log.Printf("writing to client")
			byteMessage, err := json.Marshal(mess)
			if err != nil {
				fmt.Println("error while convering")
			}
			client.conn.WriteMessage(websocket.TextMessage, byteMessage)
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

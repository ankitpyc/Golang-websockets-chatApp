package servers

import "github.com/gorilla/websocket"

type Client struct {
	hub     *SocketHub
	id      string
	conn    *websocket.Conn
	message chan Message
}

func newClient(conn *websocket.Conn, hub *SocketHub) *Client {
	return &Client{
		hub:     hub,
		conn:    conn,
		message: make(chan Message),
	}
}

type Message struct {
	MessageType string `json:"messageType"`
	Text        string `json:"text"`
	ID          string `json:"userId"`
	recieverID  string `json:"recieverID"`
	Date        int64  `json:"date"`
}

type SocketHub struct {
	unsubcribe       chan *Client
	subcribe         chan *Client
	broadCastMessage chan Message
	connectionsMap   map[string]*Client
}

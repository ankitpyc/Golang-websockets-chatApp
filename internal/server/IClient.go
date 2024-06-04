package servers

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Client struct {
	hub      *SocketHub
	id       string
	username string
	conn     *websocket.Conn
	message  chan Message
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
	UserName    string `json:"userName"`
	ID          string `json:"userId"`
	RecieverID  string `json:"recieverID"`
	Date        int64  `json:"date"`
}

type SocketHub struct {
	sync.Mutex
	unsubcribe       chan *Client
	subcribe         chan *Client
	broadCastMessage chan Message
	connectionsMap   map[string]*Client
}
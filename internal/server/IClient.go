package servers

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Acknowledgement struct {
	messageId string
	status    int
}

type Client struct {
	sync.RWMutex
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
	MessageType           string `json:"messageType"`
	Text                  string `json:"text"`
	MessageId             string `json:"messageId"`
	UserName              string `json:"userName"`
	ID                    string `json:"userId"`
	ReceiverID            string `json:"receiverID"`
	Date                  string `json:"date"`
	MessageDeliveryStatus string `json:"MessageStatus"`
}

type SocketHub struct {
	sync.RWMutex
	unsubcribe       chan *Client
	subcribe         chan *Client
	broadCastMessage chan Message
	connectionsMap   map[string]*Client
}

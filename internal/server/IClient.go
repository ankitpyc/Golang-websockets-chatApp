package servers

import (
	models "TCPServer/internal/database/models"
	"TCPServer/internal/domain"
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
	message  chan domain.Message
}

func newClient(conn *websocket.Conn, hub *SocketHub) *Client {
	return &Client{
		hub:     hub,
		conn:    conn,
		message: make(chan domain.Message),
	}
}

type SocketHub struct {
	sync.RWMutex
	DB               *models.DBServer
	unsubcribe       chan *Client
	subcribe         chan *Client
	broadCastMessage chan domain.Message
	connectionsMap   map[string]*Client
}

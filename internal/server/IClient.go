package servers

import (
	models "TCPServer/internal/database"
	"TCPServer/internal/domain"
	"sync"

	"github.com/gorilla/websocket"
)

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

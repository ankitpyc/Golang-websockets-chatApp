package servers

import (
	models "TCPServer/internal/database"
	dto "TCPServer/internal/domain/dto"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	sync.RWMutex
	hub      *SocketHub
	id       string
	username string
	conn     *websocket.Conn
	message  chan dto.Message
}

func newClient(conn *websocket.Conn, hub *SocketHub) *Client {
	return &Client{
		hub:     hub,
		conn:    conn,
		message: make(chan dto.Message),
	}
}

type SocketHub struct {
	sync.RWMutex
	DB               *models.DBServer
	unsubcribe       chan *Client
	subcribe         chan *Client
	broadCastMessage chan dto.Message
	connectionsMap   map[string]*Client
}

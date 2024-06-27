package servers

import (
	models "TCPServer/internal/database"
	dto "TCPServer/internal/domain/dto"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	sync.RWMutex
	Hub      *SocketHub
	Id       string
	Username string
	Conn     *websocket.Conn
	Message  chan dto.Message
}

func newClient(conn *websocket.Conn, hub *SocketHub) *Client {
	return &Client{
		Hub:     hub,
		Conn:    conn,
		Message: make(chan dto.Message),
	}
}

type SocketHub struct {
	sync.RWMutex
	DB               *models.DBServer
	Unsubcribe       chan *Client
	Subcribe         chan *Client
	BroadCastMessage chan dto.Message
	ConnectionsMap   map[string]*Client
}

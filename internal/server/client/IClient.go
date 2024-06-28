package servers

import (
	models "TCPServer/internal/database"
	dto "TCPServer/internal/domain/dto"
	cache "TCPServer/internal/redis-cache"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	sync.RWMutex
	Hub         *SocketHub
	CacheClient *cache.CacheClient
	Id          string
	Username    string
	Conn        *websocket.Conn
	Message     chan dto.Message
}

func newClient(conn *websocket.Conn, hub *SocketHub, cacheClient *cache.CacheClient) *Client {
	return &Client{
		Hub:         hub,
		Conn:        conn,
		CacheClient: cacheClient,
		Message:     make(chan dto.Message),
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

package databases

import (
	databases "TCPServer/internal/database/models"
	repo "TCPServer/internal/domain/Repository"

	"gorm.io/gorm"
)

type DBServer struct {
	DB          *gorm.DB
	Config      *databases.DBConfig
	UserRepo    *repo.UserRepository
	ChatRepo    *repo.ChatRepository
	MessageRepo *repo.MessageRepository
}

func (server *DBServer) InitRepository() {
	server.UserRepo = repo.NewUserRepository(server.DB)
	server.ChatRepo = repo.NewChatRepository(server.DB)
	server.MessageRepo = repo.NewMessageRepository(server.DB)
}

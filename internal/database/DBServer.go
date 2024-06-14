package databases

import (
	"TCPServer/internal/database/models"
	repo "TCPServer/internal/domain/Repository"
	"gorm.io/gorm"
)

type DBServer struct {
	DB       *gorm.DB
	Config   *databases.DBConfig
	UserRepo *repo.UserRepository
	ChatRepo *repo.ChatRepository
}

func (server *DBServer) InitRepository() {
	server.UserRepo = repo.NewUserRepository(server.DB)
	server.ChatRepo = repo.NewChatRepository(server.DB)
}

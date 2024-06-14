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
	chatRepo *repo.ChatRepositoryInf
}

func (server *DBServer) InitRepository() {
	server.UserRepo = repo.NewUserRepository(server.DB)
}

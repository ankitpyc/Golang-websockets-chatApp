// This package will implement the repository interface and provide abstraction over direct db calls

package Repository

import databases "TCPServer/internal/database/models"

type UserRepositoryInf interface {
	CreateUser(user databases.User) (uint8, error)
	RemoveUser(user databases.User) (bool, error)
	Login(user databases.User) (bool, error)
	UpdateUserImage() (bool, error)
	FindById(id uint8) (databases.User, error)
}

type ChatRepositoryInf interface {
	FetchAllChats(user databases.User) ([]databases.Chats, error)
}

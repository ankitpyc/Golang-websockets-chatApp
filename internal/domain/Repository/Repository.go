// This package will implement the repository interface and provide abstraction over direct db calls

package Repository

import (
	databases "TCPServer/internal/database/models"
	"context"
)

type UserRepositoryInf interface {
	CreateUser(ctx context.Context, user databases.User) (uint8, error)
	RemoveUser(user databases.User) (bool, error)
	Login(user databases.User) (bool, error)
	UpdateUserImage() (bool, error)
	FindById(id uint8) (databases.User, error)
}

type ChatRepositoryInf interface {
	FetchChats(chatId string) ([]*databases.Message, error)
	LoadAllUserChats(user *databases.User) ([]*databases.Chat, error)
	FetchChatByUser(user1 string, user2 string)
}

type MessageRepositoryInf interface {
	SaveUserChat(chat *databases.Message) *databases.Message
}

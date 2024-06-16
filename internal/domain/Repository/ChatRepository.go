package Repository

import (
	databases "TCPServer/internal/database/models"
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type ChatRepository struct {
	db             *gorm.DB
	ctx            context.Context
	defaultTimeout time.Duration
}

func NewChatRepository(db *gorm.DB) *ChatRepository {
	return &ChatRepository{ctx: context.Background(), db: db}
}

func (cr *ChatRepository) FetchChats(chatId uint) ([]*databases.Message, error) {
	ctx, cancel := context.WithTimeout(cr.ctx, time.Second*5)
	mess := make([]*databases.Message, 0)
	defer cancel()
	result := cr.db.WithContext(ctx).Limit(20).Where("chat_id = ?", chatId).First(&mess)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed")
	}
	return mess, nil
}

func (cr *ChatRepository) LoadAllUserChats(userid string) ([]*databases.Chat, error) {
	ctx, cancel := context.WithTimeout(cr.ctx, time.Second*5)
	chats := make([]*databases.Chat, 0)
	defer cancel()
	result := cr.db.WithContext(ctx).Where("user_id1 = ? OR user_id2 = ?", userid, userid).
		Preload("Messages").
		Find(&chats)
	if result.Error != nil {
		return nil, errors.New("error fetching user chats")
	}
	return chats, nil
}

func (cr *ChatRepository) FetchChatByUser(user1 string, user2 string) (uint, error) {
	ctx, cancel := context.WithTimeout(cr.ctx, time.Second*5)
	var chat databases.Chat
	defer cancel()
	result := cr.db.WithContext(ctx).
		Where("(user_id_1 = ? AND user_id_2 = ?) OR (user_id_2 = ? AND user_id_1 = ?)", user1, user2, user1, user2).
		First(&chat)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return 0, nil

	}
	if result.Error != nil {
		return 0, nil
	}

	return chat.ID, nil
}

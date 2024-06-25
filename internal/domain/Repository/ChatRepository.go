package Repository

import (
	databases "TCPServer/internal/database/models"
	"context"
	"errors"
	"fmt"
	"strconv"
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
	result := cr.db.WithContext(ctx).Limit(20).Where("chat_id = ?", chatId).Order("created_at").Find(&mess)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed")
	}
	return mess, nil
}

func (cr *ChatRepository) LoadAllUserChats(userid uint) ([]databases.Chat, error) {
	ctx, cancel := context.WithTimeout(cr.ctx, time.Second*12)
	chats := make([]databases.Chat, 0)
	uid := strconv.FormatUint(uint64(userid), 10)
	defer cancel()
	result := cr.db.WithContext(ctx).Where("user_id_1 = ? OR user_id_2 = ?", uid, uid).
		Preload("User1").
		Preload("User2").
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at asc").Limit(30)
		}).
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

package Repository

import (
	databases "TCPServer/internal/database/models"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type ChatRepository struct {
	db             *gorm.DB
	ctx            context.Context
	defaultTimeout time.Duration
}

func NewChatRepository(db *gorm.DB) *ChatRepository {
	return &ChatRepository{ctx: context.Background(), db: db}
}

func (cr *ChatRepository) FetchChats(chatId string) ([]*databases.Message, error) {
	ctx, cancel := context.WithTimeout(cr.ctx, time.Second*5)
	mess := make([]*databases.Message, 0)
	defer cancel()
	result := cr.db.WithContext(ctx).Limit(20).Where("chat_id = ?", chatId).First(&mess)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed")
	}
	return mess, nil

}

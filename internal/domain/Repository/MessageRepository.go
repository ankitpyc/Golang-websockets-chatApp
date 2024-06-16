package Repository

import (
	databases "TCPServer/internal/database/models"
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type MessageRepository struct {
	db             *gorm.DB
	ctx            context.Context
	defaultTimeout time.Duration
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{ctx: context.Background(), db: db, defaultTimeout: 5 * time.Second}
}

func (mr *MessageRepository) SaveMessage(message *databases.Message) (*databases.Message, error) {
	ctx, cancel := context.WithTimeout(mr.ctx, time.Second*5)
	defer cancel()
	result := mr.db.WithContext(ctx).Create(&message)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed")
	}
	return message, nil
}

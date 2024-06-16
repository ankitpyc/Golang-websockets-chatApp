package databases

import (
	"time"

	"gorm.io/gorm"
)

type ChatStatus int

// Declare related constants for each weekday starting with index 1

type Chat struct {
	gorm.Model
	UserID1      string `gorm:"column:user_id_1"`
	UserID2      string `gorm:"column:user_id_2"`
	LastActivity time.Time
	IsDeleted    bool
	Messages     []Message `gorm:"foreignKey:chat_id;references:ID"`
}

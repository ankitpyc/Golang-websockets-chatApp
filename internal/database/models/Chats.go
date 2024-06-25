package databases

import (
	"time"

	"gorm.io/gorm"
)

type ChatStatus int

// Declare related constants for each weekday starting with index 1

type Chat struct {
	gorm.Model
	UserID1      string     `gorm:"column:user_id_1" json:"userID1"`
	UserID2      string     `gorm:"column:user_id_2" json:"userID2"`
	LastActivity time.Time  `json:"time"`
	IsDeleted    bool       `json:"isdeleted"`
	Messages     []*Message `gorm:"foreignKey:chat_id;references:ID" json:"messages"`
	User1        User       `gorm:"foreignKey:user_id_1" json:"user1Details"`
	User2        User       `gorm:"foreignKey:user_id_2" json:"user2Details"`
}

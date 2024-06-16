package databases

import "gorm.io/gorm"

type MessageStatus int

const (
	PUSHED    MessageStatus = iota + 1 // EnumIndex = 1
	DELIVERED                          // EnumIndex = 2
	READ                               // EnumIndex = 3
)

type Message struct {
	*gorm.Model
	MessageID  string        `gorm:"primaryKey"`
	ChatID     uint          `gorm:"column:chat_id"` // Matches the column name in the database
	SenderID   string        `gorm:"column:sender_id"`
	ReceiverID string        `gorm:"column:receiver_id"`
	Status     MessageStatus `gorm:"column:status"`
	Text       string        `gorm:"column:message"`
}

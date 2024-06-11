package databases

import "gorm.io/gorm"

type MessageStatus int

const (
	PUSHED    MessageStatus = iota + 1 // EnumIndex = 1
	DELIVERED                          // EnumIndex = 2
	READ                               // EnumIndex = 3
)

type Message struct {
	gorm.Model
	ChatId     string        `gorm:"column:chat_id"`
	SenderID   string        `gorm:"column:sender_id"` // Custom foreign key name
	ReceiverID string        `gorm:"column:receiver_id"`
	Status     MessageStatus `gorm:"column:status"`
	Text       string        `gorm:"column:message"`
}

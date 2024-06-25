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
	MessageID  string        `gorm:"primaryKey" json:"messageId"`
	ChatID     uint          `gorm:"column:chat_id" json:"chatId"` // Matches the column name in the database
	SenderID   string        `gorm:"column:sender_id" json:"senderId"`
	ReceiverID string        `gorm:"column:receiver_id" json:"receiverId"`
	Status     MessageStatus `gorm:"column:status" json:"status"`
	Text       string        `gorm:"column:message" json:"message"`
}

package databases

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	ChatId     uint `gorm:"column:chat_id"`
	SenderID   uint `gorm:"column:sender_id"` // Custom foreign key name
	ReceiverID uint `gorm:"column:receiver_id"`
	Text       string
}

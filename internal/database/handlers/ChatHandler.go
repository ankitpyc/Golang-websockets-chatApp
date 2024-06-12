package databases

import (
	databases "TCPServer/internal/database/models"
	"TCPServer/internal/domain"
	"TCPServer/internal/server/messageUtil"
	"fmt"
	"github.com/google/uuid"
)

type ChatHandlerInf interface {
	PersistMessages(message *databases.Message) error
	CreateChat(message *databases.Message) error
	SendAcknowledgment(message *databases.Message) error
}

type ChatHandler struct {
	*databases.DBServer
}

func NewChatHandler() *ChatHandler {
	return &ChatHandler{}
}

func (ch *ChatHandler) PersistMessages(message *domain.Message) error {
	if message.MessageType == "ACK" {
		return nil
	}

	if message.ChatId == "" {
		id, err := ch.CreateChat(message)
		message.ChatId = id
		if err != nil {
			return err
		}
	}
	chatMessage := messageUtil.ConvertToChatMessage(message)
	ch.DB.Create(&chatMessage)
	return nil
}

func (ch *ChatHandler) CreateChat(message *domain.Message) (string, error) {
	chat := databases.Chats{
		ChatId:    uuid.New().String(),
		UserID1:   message.ID,
		UserID2:   message.ReceiverID,
		IsDeleted: false,
	}
	id := ch.DB.Create(&chat)
	fmt.Println(id)
	return chat.ChatId, nil
}

func (ch *ChatHandler) SendAcknowledgement(message *domain.Message) (*domain.Message, error) {
	chat := &domain.Message{
		ChatId:                message.ChatId,
		ID:                    message.ReceiverID,
		ReceiverID:            message.ID,
		MessageType:           "ACK",
		Text:                  "OYE",
		MessageId:             message.MessageId,
		MessageDeliveryStatus: "SENT",
		Date:                  message.Date,
	}

	return chat, nil

}

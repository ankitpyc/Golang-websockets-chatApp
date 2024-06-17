package databases

import (
	databases2 "TCPServer/internal/database"
	databases "TCPServer/internal/database/models"
	"TCPServer/internal/domain"
	"TCPServer/internal/server/messageUtil"
	"fmt"
	"strconv"
	"time"
)

type ChatHandlerInf interface {
	PersistMessages(message *databases.Message) error
	CreateChat(message *databases.Message) error
	SendAcknowledgment(message *databases.Message) error
}

type ChatHandler struct {
	*databases2.DBServer
}

func NewChatHandler() *ChatHandler {
	return &ChatHandler{}
}

func (ch *ChatHandler) PersistMessages(message *domain.Message) error {
	if message.MessageType == "ACK" {
		return nil
	}
	ids, _ := ch.ChatRepo.FetchChatByUser(message.ID, message.ReceiverID)
	message.ChatId = strconv.FormatUint(uint64(ids), 10)
	if ids == 0 {
		id, err := ch.CreateChat(message)
		message.ChatId = strconv.FormatUint(uint64(id), 10)
		if err != nil {
			return err
		}
	}
	chatMessage := messageUtil.ConvertToChatMessage(message)
	ch.MessageRepo.SaveMessage(&chatMessage)
	return nil
}

func (ch *ChatHandler) CreateChat(message *domain.Message) (uint, error) {
	chat := databases.Chat{
		UserID1:      message.ID,
		UserID2:      message.ReceiverID,
		IsDeleted:    false,
		LastActivity: time.Now(),
	}
	result := ch.DB.Create(&chat)

	if result.Error != nil {
		fmt.Print(result.Error)
		return 0, nil
	}
	return chat.ID, nil
}

func (ch *ChatHandler) SendAcknowledgement(message *domain.Message) (*domain.Message, error) {
	chat := &domain.Message{
		ChatId:                message.ChatId,
		ID:                    message.ReceiverID,
		ReceiverID:            message.ID,
		MessageType:           "ACK",
		Text:                  "",
		MessageId:             message.MessageId,
		MessageDeliveryStatus: "SENT",
		Date:                  message.Date,
	}

	return chat, nil

}

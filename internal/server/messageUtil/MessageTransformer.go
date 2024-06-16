package messageUtil

import (
	databases "TCPServer/internal/database/models"
	"TCPServer/internal/domain"
	"strconv"

	"github.com/google/uuid"
)

func ConvertToChatMessage(message *domain.Message) databases.Message {
	chatId, _ := StringToUint(message.ChatId)
	mess := databases.Message{
		MessageID:  uuid.New().String(),
		ChatID:     chatId,
		Text:       message.Text,
		SenderID:   message.ID,
		ReceiverID: message.ReceiverID,
	}
	return mess
}

func StringToUint(str string) (uint, error) {
	value, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(value), nil
}

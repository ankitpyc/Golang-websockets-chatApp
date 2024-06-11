package messageUtil

import (
	databases "TCPServer/internal/database/models"
	"TCPServer/internal/domain"
)

func ConvertToChatMessage(message *domain.Message) databases.Message {
	mess := databases.Message{
		ChatId:     message.ChatId,
		Text:       message.Text,
		SenderID:   message.ID,
		ReceiverID: message.ReceiverID,
	}

	return mess
}

package dto

type Message struct {
	MessageType           string `json:"messageType"`
	Text                  string `json:"text"`
	MessageId             string `json:"messageId"`
	UserName              string `json:"userName"`
	ID                    string `json:"userId"`
	ReceiverID            string `json:"receiverID"`
	Date                  string `json:"date"`
	ChatId                string `json:"chatId,omitempty"`
	MessageDeliveryStatus string `json:"MessageStatus"`
}

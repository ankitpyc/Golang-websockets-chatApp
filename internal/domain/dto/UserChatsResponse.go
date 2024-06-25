package dto

import (
	dbmodels "TCPServer/internal/database/models"
)

type UserDataResponse struct {
	UserId       uint            `json:"userid"`
	UserName     string          `json:"username"`
	ChatsHistory []dbmodels.Chat `json:"history"`
}

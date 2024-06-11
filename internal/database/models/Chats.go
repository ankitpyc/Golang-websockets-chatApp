package databases

import "gorm.io/gorm"

type ChatStatus int

// Declare related constants for each weekday starting with index 1

type Chats struct {
	*gorm.Model
	ChatId    string
	UserID1   string // Foreign key referencing User.ID
	UserID2   string
	IsDeleted bool
}

package databases

type Chats struct {
	ChatId    string
	UserID1   uint // Foreign key referencing User.ID
	UserID2   uint
	IsDeleted bool
}
